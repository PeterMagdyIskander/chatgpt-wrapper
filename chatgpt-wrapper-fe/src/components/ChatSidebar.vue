<template>
  <div>
    <div class="sidebar" :class="{ open: isOpen }">
      <div class="sidebar-header">
        <h3>Chat with AI</h3>
        <button class="close-btn" @click="onClose">×</button>
      </div>

      <div class="sidebar-content">
        <ChatMessages
          ref="chatMessagesComponent"
          :messages="messages"
          :is-loading="isLoading"
        />
      </div>

      <ChatInput
        ref="chatInputComponent"
        :is-loading="isLoading"
        @send-message="onSendMessage"
      />
    </div>

    <!-- Overlay -->
    <div class="overlay" :class="{ visible: isOpen }" @click="onClose"></div>
  </div>
</template>

<script>
import { ref } from 'vue'
import ChatMessages from './ChatMessages.vue'
import ChatInput from './ChatInput.vue'

export default {
  name: 'ChatSidebar',
  components: {
    ChatMessages,
    ChatInput
  },
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    messages: {
      type: Array,
      default: () => []
    },
    isLoading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['close', 'send-message'],
  setup(props, { emit }) {
    const chatMessagesComponent = ref(null)
    const chatInputComponent = ref(null)

    const onClose = () => {
      emit('close')
    }

    const onSendMessage = (message) => {
      emit('send-message', message)
    }

    const scrollToBottom = () => {
      if (chatMessagesComponent.value) {
        chatMessagesComponent.value.setShouldScrollToBottom()
      }
    }

    const focusInput = () => {
      if (chatInputComponent.value) {
        chatInputComponent.value.focus()
      }
    }

    return {
      chatMessagesComponent,
      chatInputComponent,
      onClose,
      onSendMessage,
      scrollToBottom,
      focusInput
    }
  }
}
</script>

<style scoped>
.sidebar {
  position: fixed;
  right: -400px;
  top: 0;
  width: 400px;
  height: 100vh;
  background: white;
  box-shadow: -2px 0 20px rgba(0, 0, 0, 0.1);
  transition: right 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  pointer-events: auto;
  z-index: 1002;
}

.sidebar.open {
  right: 0;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}

.sidebar-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6b7280;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #e5e7eb;
  color: #374151;
}

.sidebar-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
  pointer-events: none;
  z-index: 1001;
}

.overlay.visible {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
}

@media (max-width: 768px) {
  .sidebar {
    width: 100vw;
    right: -100vw;
  }

  .sidebar.open {
    right: 0;
  }
}
</style>