<template>
  <div class="container">
    <div class="message-input-container">
      <input v-model="inputMessage" @keyup.enter="sendMessage" placeholder="Type your message..." class="message-input" :disabled="inputMessage.length>charLimit">
      
      <button @click="sendMessage" class="send-button" >
        Send
      </button>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'MessageInput',
  emits: ['send-message'],
  props: {
    charLimit: {
      type: Number,
      required: true
    },
  },
  setup(props,{ emit }) {
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
.container {
  height: 122px;
  background-color: rgba(51, 46, 56, 0.4);
  border: 1px solid rgb(39, 36, 44);
  padding: 8px 8px 0 8px;
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;
  position: relative;
}

.message-input-container {
  padding: 12px;
  background-color: #2c2431;
  border: 1px solid #2c2431;
  color: rgb(212, 199, 225);
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;
}

.message-input {
  width: 100%;
  font-size: 16px;
  background-color: transparent;
  border: unset;
  color: rgb(249, 248, 251);
}

.message-input:focus-visible {
  border: unset;
  outline: unset;
}

.message-input:disabled{
  border: 1px solid red !important;
}

.send-button {
  position: absolute;
  bottom: 4px;
  right: 4px;
  padding: 10px;
  color: rgb(253, 242, 248);
  cursor: pointer;
  font-size: 16px;
  border-radius: 8px;
  background-color: rgba(163, 0, 76, 0.2);
  border: 1px solid rgba(163, 0, 76, 0.1);
}
</style>