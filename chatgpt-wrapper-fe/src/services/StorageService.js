export const StorageService = {
  /**
   * Store data in localStorage
   * @param {string} key - The key to store data under
   * @param {string} value - The value to store
   */
  setItem(key, value) {
    try {
      localStorage.setItem(key, value)
    } catch (e) {
      console.error('Failed to save to localStorage:', e)
      // Handle storage quota exceeded or other errors
      if (e.name === 'QuotaExceededError') {
        console.warn('localStorage quota exceeded. Consider clearing old data.')
      }
    }
  },

  /**
   * Retrieve data from localStorage
   * @param {string} key - The key to retrieve data for
   * @returns {string|null} The stored value or null if not found
   */
  getItem(key) {
    try {
      return localStorage.getItem(key)
    } catch (e) {
      console.error('Failed to read from localStorage:', e)
      return null
    }
  },

  /**
   * Remove data from localStorage
   * @param {string} key - The key to remove
   */
  removeItem(key) {
    try {
      localStorage.removeItem(key)
    } catch (e) {
      console.error('Failed to remove from localStorage:', e)
    }
  },

  /**
   * Clear all data from localStorage
   */
  clear() {
    try {
      localStorage.clear()
    } catch (e) {
      console.error('Failed to clear localStorage:', e)
    }
  },

  /**
   * Get all keys from localStorage
   * @returns {string[]} Array of all keys
   */
  getAllKeys() {
    try {
      return Object.keys(localStorage)
    } catch (e) {
      console.error('Failed to get keys from localStorage:', e)
      return []
    }
  },

  /**
   * Check if localStorage is available
   * @returns {boolean} True if localStorage is available
   */
  isAvailable() {
    try {
      const testKey = '__localStorage_test__'
      localStorage.setItem(testKey, 'test')
      localStorage.removeItem(testKey)
      return true
    } catch (e) {
      return false
    }
  }
}