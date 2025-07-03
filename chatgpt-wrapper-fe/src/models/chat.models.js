// Chat message interface
export const createChatMessage = (role, content, timestamp = new Date()) => ({
  role, // 'user' | 'assistant'
  content,
  timestamp
})

// OpenAI API interfaces
export const createOpenAIMessage = (role, content) => ({
  role, // 'system' | 'user' | 'assistant'
  content
})

export const createOpenAIRequest = (model, messages, maxTokens, temperature,stream) => ({
  model,
  messages,
  max_tokens: maxTokens,
  temperature,
  stream
})

// Default OpenAI response structure
export const createOpenAIResponse = () => ({
  id: '',
  object: '',
  created: 0,
  model: '',
  choices: [],
  usage: {
    prompt_tokens: 0,
    completion_tokens: 0,
    total_tokens: 0
  }
})

// Formatted chat message interface
export const createFormattedChatMessage = (content, htmlContent, role, timestamp = new Date(), options = {}) => ({
  content,
  htmlContent,
  role,
  timestamp,
  tokenUsage: options.tokenUsage || null,
  finishReason: options.finishReason || null,
  isError: options.isError || false
})

// Chat configuration interface
export const createChatConfig = (openaiApiKey, openaiApiUrl, additionalData = {}) => ({
  openaiApiKey,
  openaiApiUrl,
  ...additionalData
})

// Token usage helper
export const createTokenUsage = (prompt, completion, total) => ({
  prompt,
  completion,
  total
})