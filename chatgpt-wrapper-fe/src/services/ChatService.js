export const ChatService = {
  async sendMessage(message, userId) {
    try {
      const response = await fetch("http://localhost:8081/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          message: message,
          userId: userId,
        }),
      });

      if (response.ok) {
        return await response.json();
      } else {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(
          `HTTP error! status: ${response.status}, message: ${
            errorData.message || "Unknown error"
          }`
        );
      }
    } catch (error) {
      console.error("Error sending message:", error);
      throw error;
    }
  },

  async getChatLimit() {
    try {
      const response = await fetch("http://localhost:8081/char-limit", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        return await response.json();
      } else {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(
          `HTTP error! status: ${response.status}, message: ${
            errorData.message || "Unknown error"
          }`
        );
      }
    } catch (error) {
      console.error("Error sending message:", error);
      throw error;
    }
  },

  createSSEConnection(userId, messageId) {
    const url = `http://localhost:8081/ask-chatgpt?userId=${encodeURIComponent(
      userId
    )}&messageId=${encodeURIComponent(messageId)}`;
    return new EventSource(url);
  },
};
