<template>
  <div class="message-input-container">
    <input 
      v-model="inputMessage" 
      @keyup.enter="sendMessage"
      placeholder="Type your message..."
      class="message-input"
    >
    <button @click="sendMessage" class="send-button">
      Send
    </button>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'MessageInput',
  emits: ['send-message'],
  setup(props, { emit }) {
    const inputMessage = ref('')

    function sendMessage() {
      if (!inputMessage.value.trim()) return
      
      const message = inputMessage.value.trim()
      emit('send-message', message)
      inputMessage.value = ''
    }

    return {
      inputMessage,
      sendMessage
    }
  }
}
</script>

<style scoped>
.message-input-container {
  padding: 10px;
  border-top: 1px solid #ccc;
}

.message-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  margin-bottom: 10px;
}

.message-input:focus {
  outline: none;
  border-color: #007bff;
}

.send-button {
  padding: 10px 20px;
  border: 1px solid #007bff;
  background: #007bff;
  color: white;
  cursor: pointer;
}

.send-button:hover {
  background: #0056b3;
}

.send-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style>