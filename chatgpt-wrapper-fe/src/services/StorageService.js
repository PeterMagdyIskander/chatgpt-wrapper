export const StorageService = {
  setItem(key, value) {
    try {
      localStorage.setItem(key, value);
    } catch (e) {
      console.error("Failed to save to localStorage:", e);
    }
  },

  getItem(key) {
    try {
      return localStorage.getItem(key);
    } catch (e) {
      console.error("Failed to read from localStorage:", e);
      return null;
    }
  },
};
