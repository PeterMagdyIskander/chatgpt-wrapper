import {
  createOpenAIMessage,
  createOpenAIRequest,
  createFormattedChatMessage,
} from "../models/chat.models.js";
import { marked } from "marked";
export class ChatService {
  constructor(httpClient) {
    this.httpClient = httpClient;
    this.lastApiCall = 0;
    this.API_CALL_DELAY = 1000; // 1 second between calls
  }

  async callOpenAI(message, chatHistory, config) {
    if (!config.openaiApiKey) {
      throw new Error("OpenAI API key not provided");
    }

    const systemPrompt = this.createSystemPrompt(config);

    const messages = [
      createOpenAIMessage("system", systemPrompt),
      ...chatHistory
        .filter(
          (msg) =>
            msg.role !== "assistant" ||
            (!msg.content.includes("Hello! I'm here to help") &&
              !msg.content.includes("How can I assist you today?"))
        )
        .slice(-10)
        .map((msg) =>
          createOpenAIMessage(
            msg.role === "assistant" ? "assistant" : "user",
            msg.content
          )
        ),
      createOpenAIMessage("user", message),
    ];

    const requestBody = createOpenAIRequest(
      "gpt-4o-mini",
      messages,
      1000,
      0.7,
      false
    );

    try {
      const response = await this.httpClient.post(
        config.openaiApiUrl,
        requestBody,
        {
          headers: {
            Authorization: `Bearer ${config.openaiApiKey}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (
        !response.data ||
        !response.data.choices ||
        response.data.choices.length === 0
      ) {
        throw new Error("Invalid response from OpenAI");
      }

      return response.data.choices[0].message.content;
    } catch (error) {
      console.error("OpenAI API Error:", error);
      return this.getErrorMessage(error);
    }
  }

  getErrorMessage(error) {
    const errorMessages = {
      429: "I'm currently experiencing high demand. Please try again in a few moments, or check your OpenAI API quota if you're using your own key.",
      401: "There seems to be an authentication issue with the API. Please check your API key.",
      403: "Access forbidden. Please check your API key permissions.",
      500: "The AI service is temporarily unavailable. Please try again later.",
      default:
        "I encountered an unexpected error. Please try again or contact support if the issue persists.",
    };

    const status = error.response?.status || error.status;
    return errorMessages[status] || errorMessages.default;
  }

  createSystemPrompt(config) {
    // Generic system prompt - can be customized based on config
    const basePrompt = `Task Instructions:
You will be provided with a question from the user, and you need to provide an answer based on the data provided.
If the question can be answered with the data provided, you should provide a direct answer.
If the question requires reasoning or analysis, you should provide a detailed explanation of your reasoning process and the steps you took to arrive at your answer.
Your answer should always be structured and use fun emojis.
You should also provide a concise title for your response, it should be the very first thing you provide in your message, 
and it should be formatted like this title="your title should be here" you should 2 $$ after it i.e it should look something like this
title="Todo list creation"$$your response should be here, you should only return the title for the first message in the conversation`;

    // If config contains additional data, include it in the prompt
    if (config.additionalData) {
      const jsonData = JSON.stringify(config.additionalData, null, 2);
      return `${basePrompt}\n\nAdditional context data: ${jsonData}`;
    }

    return basePrompt;
  }

  createWelcomeMessage() {
    const welcomeContent =
      "Hello! I'm your AI assistant. How can I help you today?";

    return createFormattedChatMessage(
      welcomeContent,
      `<p>${welcomeContent}</p>`,
      "assistant",
      new Date(),
      { isError: false }
    );
  }

  createUserMessage(content) {
    return createFormattedChatMessage(
      content,
      this.responseFormatter?.formatToHTML?.(content) || `<p>${content}</p>`,
      "user",
      new Date(),
      { isError: false }
    );
  }

  createAssistantMessage(content, isError = false) {
    if (content.includes("title=")) {
      [title, content] = content.split("$$");
    }
    return createFormattedChatMessage(
      content,
      marked.parse(content),
      "assistant",
      new Date(),
      { isError }
    );
  }

  // Rate limiting check
  async checkRateLimit() {
    const now = Date.now();
    const timeSinceLastCall = now - this.lastApiCall;
    if (timeSinceLastCall < this.API_CALL_DELAY) {
      await new Promise((resolve) =>
        setTimeout(resolve, this.API_CALL_DELAY - timeSinceLastCall)
      );
    }
    this.lastApiCall = Date.now();
  }
}
