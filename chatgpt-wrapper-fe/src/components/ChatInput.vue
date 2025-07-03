<template>
  <div class="input-area">
    <div class="input-container">
      <textarea
        v-model="currentMessage"
        @keydown.enter="handleKeyDown"
        placeholder="Ask me anything..."
        class="message-input"
        rows="1"
        ref="messageInput"
      ></textarea>
      <button
        @click="onSend"
        :disabled="!currentMessage.trim() || isLoading"
        class="send-button"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M2 21L23 12L2 3V10L17 12L2 14V21Z" fill="currentColor"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script>
import { ref, nextTick } from 'vue'

export default {
  name: 'ChatInput',
  props: {
    isLoading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['send-message'],
  setup(props, { emit }) {
    const currentMessage = ref('')
    const messageInput = ref(null)

    const handleKeyDown = (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        onSend()
      }
    }

    const onSend = () => {
      if (!currentMessage.value.trim() || props.isLoading) return
      
      const message = currentMessage.value.trim()
      currentMessage.value = ''
      emit('send-message', message)
    }

    const focus = () => {
      nextTick(() => {
        messageInput.value?.focus()
      })
    }

    const clear = () => {
      currentMessage.value = ''
    }

    return {
      currentMessage,
      messageInput,
      handleKeyDown,
      onSend,
      focus,
      clear
    }
  }
}
</script>

<style scoped>
.input-area {
  padding: 20px;
  border-top: 1px solid #e5e7eb;
  background: white;
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.message-input {
  flex: 1;
  border: 1px solid #d1d5db;
  border-radius: 20px;
  padding: 12px 16px;
  resize: none;
  outline: none;
  font-family: 'Segoe UI', 'Tahoma', 'Arial', 'Helvetica Neue', sans-serif;
  font-size: 14px;
  line-height: 1.4;
  max-height: 120px;
  transition: border-color 0.2s;
}

.message-input:focus {
  border-color: #667eea;
}

.send-button {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-button:hover:not(:disabled) {
  transform: scale(1.05);
}

.send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}
</style>