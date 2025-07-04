/**
 * Generate a unique user ID
 * @returns {string} Unique user ID
 */
export function generateUserId() {
  return 'user_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
}

/**
 * Generate a unique chat ID
 * @returns {string} Unique chat ID
 */
export function generateChatId() {
  return 'chat_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
}

/**
 * Generate a unique message ID
 * @returns {string} Unique message ID
 */
export function generateMessageId() {
  return 'msg_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
}

/**
 * Format timestamp for display
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Formatted time string
 */
export function formatTimestamp(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

/**
 * Format date for display
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Formatted date string
 */
export function formatDate(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleDateString()
}

/**
 * Truncate text to specified length
 * @param {string} text - The text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string} Truncated text
 */
export function truncateText(text, maxLength = 50) {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

/**
 * Debounce function to limit function calls
 * @param {Function} func - Function to debounce
 * @param {number} delay - Delay in milliseconds
 * @returns {Function} Debounced function
 */
export function debounce(func, delay) {
  let timeoutId
  return function (...args) {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => func.apply(this, args), delay)
  }
}

/**
 * Validate if a string is a valid URL
 * @param {string} string - String to validate
 * @returns {boolean} True if valid URL
 */
export function isValidUrl(string) {
  try {
    new URL(string)
    return true
  } catch (_) {
    return false
  }
}

/**
 * Escape HTML characters in a string
 * @param {string} text - Text to escape
 * @returns {string} Escaped text
 */
export function escapeHtml(text) {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  }
  return text.replace(/[&<>"']/g, (m) => map[m])
}

/**
 * Deep clone an object
 * @param {Object} obj - Object to clone
 * @returns {Object} Cloned object
 */
export function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') return obj
  if (obj instanceof Date) return new Date(obj)
  if (obj instanceof Array) return obj.map(item => deepClone(item))
  if (typeof obj === 'object') {
    const clonedObj = {}
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key])
      }
    }
    return clonedObj
  }
}