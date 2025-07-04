export const ChatService = {
  /**
   * Send a message to the backend
   * @param {string} message - The message to send
   * @param {string} userId - The user ID
   * @returns {Promise<Object>} Response from the server
   */
  async sendMessage(message, userId) {
    try {
      const response = await fetch('http://localhost:8081/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          message: message,
          userId: userId
        })
      })

      if (response.ok) {
        return await response.json()
      } else {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorData.message || 'Unknown error'}`)
      }
    } catch (error) {
      console.error('Error sending message:', error)
      throw error
    }
  },

  /**
   * Create an SSE connection for streaming responses
   * @param {string} userId - The user ID
   * @param {string} messageId - The message ID
   * @returns {EventSource} SSE connection
   */
  createSSEConnection(userId, messageId) {
    const url = `http://localhost:8081/ask-chatgpt?userId=${encodeURIComponent(userId)}&messageId=${encodeURIComponent(messageId)}`
    return new EventSource(url)
  },

  /**
   * Check if the backend is available
   * @returns {Promise<boolean>} True if backend is available
   */
  async checkBackendHealth() {
    try {
      const response = await fetch('http://localhost:8081/health', {
        method: 'GET',
        timeout: 5000
      })
      return response.ok
    } catch (error) {
      console.error('Backend health check failed:', error)
      return false
    }
  },

  /**
   * Get chat history from backend (if supported)
   * @param {string} userId - The user ID
   * @returns {Promise<Array>} Array of chat messages
   */
  async getChatHistory(userId) {
    try {
      const response = await fetch(`http://localhost:8081/chat-history/${encodeURIComponent(userId)}`, {
        method: 'GET'
      })

      if (response.ok) {
        return await response.json()
      } else {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
    } catch (error) {
      console.error('Error fetching chat history:', error)
      throw error
    }
  }
}