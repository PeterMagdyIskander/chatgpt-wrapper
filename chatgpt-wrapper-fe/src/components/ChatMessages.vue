<template>
  <div class="chat-messages" ref="chatMessages">
    <div
      v-for="message in messages"
      :key="message.timestamp"
      class="message"
      :class="{
        user: isUserMessage(message),
        assistant: isAssistantMessage(message),
        error: message.isError
      }"
    >
      <div class="message-content">
        <div class="message-text" v-html="message.htmlContent"></div>
        <div class="message-meta">
          <div class="message-time">
            {{ formatTime(message.timestamp) }}
          </div>
          <div
            v-if="message.tokenUsage && isAssistantMessage(message)"
            class="token-usage"
          >
            {{ getTokenUsageSummary(message) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Loading indicator -->
    <div v-if="isLoading" class="message assistant">
      <div class="message-content">
        <div class="loading-dots">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick, watch } from 'vue'

export default {
  name: 'ChatMessages',
  props: {
    messages: {
      type: Array,
      default: () => []
    },
    isLoading: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const chatMessages = ref(null)
    let shouldScrollToBottom = false

    const isUserMessage = (message) => {
      return message.role === 'user'
    }

    const isAssistantMessage = (message) => {
      return message.role === 'assistant'
    }

    const formatTime = (timestamp) => {
      return new Date(timestamp).toLocaleString()
    }

    const getTokenUsageSummary = (message) => {
      if (!message.tokenUsage) return ''
      return `Tokens: ${message.tokenUsage.total} (${message.tokenUsage.prompt}+${message.tokenUsage.completion})`
    }

    const scrollToBottom = () => {
      if (chatMessages.value) {
        const element = chatMessages.value
        element.scrollTop = element.scrollHeight
      }
    }

    const setShouldScrollToBottom = () => {
      shouldScrollToBottom = true
      nextTick(() => {
        if (shouldScrollToBottom) {
          scrollToBottom()
          shouldScrollToBottom = false
        }
      })
    }

    // Watch for changes in messages or loading state
    watch([() => props.messages, () => props.isLoading], () => {
      setShouldScrollToBottom()
    })

    return {
      chatMessages,
      isUserMessage,
      isAssistantMessage,
      formatTime,
      getTokenUsageSummary,
      scrollToBottom,
      setShouldScrollToBottom
    }
  }
}
</script>

<style scoped>
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background-color: #fcfcfc;
  scroll-behavior: smooth;
}

.message {
  display: flex;
  max-width: 80%;
  word-break: break-word;
}

.message.user {
  align-self: flex-end;
  justify-content: flex-end;
}

.message.assistant {
  align-self: flex-start;
  justify-content: flex-start;
}

.message-content {
  padding: 12px 18px;
  border-radius: 20px;
  max-width: 100%;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  background-color: #ffffff;
  color: #333;
  border: 1px solid #e0e0e0;
}

.message.user .message-content {
  background: linear-gradient(135deg, #4c68d1 0%, #6f42c1 100%);
  color: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  border: none;
}

.message.error .message-content {
  background: #ffebee;
  border: 1px solid #ef9a9a;
  color: #d32f2f;
}

.message-text {
  line-height: 1.6;
  font-size: 15px;
  color: #333;
}

.message-text :deep(p) {
  margin: 8px 0;
}

.message-text :deep(p:first-child) {
  margin-top: 0;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(ul),
.message-text :deep(ol) {
  margin: 8px 0;
  list-style-position: outside;
  padding-left: 16px;
}

.message-text :deep(li) {
  margin: 4px 0;
}

.message-text :deep(ul ul),
.message-text :deep(ol ol),
.message-text :deep(ul ol),
.message-text :deep(ol ul) {
  margin-top: 4px;
  margin-bottom: 4px;
}

.message-text :deep(strong) {
  font-weight: 700;
}

.message-text :deep(em) {
  font-style: italic;
}

.message-text :deep(code) {
  background-color: #f0f0f0;
  padding: 3px 6px;
  border-radius: 5px;
  font-family: 'SFMono-Regular', 'Menlo', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  font-size: 0.88em;
  color: #c7254e;
}

.message.user .message-text :deep(code) {
  background: rgba(255, 255, 255, 0.25);
  color: #fff;
}

.message-text :deep(blockquote) {
  border-left: 4px solid #ccc;
  margin: 10px 0;
  padding: 5px 15px;
  background-color: #f9f9f9;
  color: #555;
  font-style: italic;
}

.message-text :deep(pre) {
  background-color: #2d2d2d;
  color: #f8f8f2;
  padding: 10px 15px;
  border-radius: 8px;
  overflow-x: auto;
  font-family: 'SFMono-Regular', 'Menlo', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  font-size: 0.9em;
  line-height: 1.4;
  margin: 10px 0;
}

.message-text :deep(pre code) {
  background: none;
  color: inherit;
  padding: 0;
  font-size: inherit;
  border-radius: 0;
}

.message.user .message-text {
  color: white;
}

.message.user .message-text :deep(p),
.message.user .message-text :deep(ul),
.message.user .message-text :deep(ol),
.message.user .message-text :deep(li),
.message.user .message-text :deep(strong),
.message.user .message-text :deep(em) {
  color: white;
}

.message.user .message-text :deep(blockquote) {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  border-color: rgba(255, 255, 255, 0.5);
}

.message.user .message-text :deep(pre) {
  background-color: rgba(0, 0, 0, 0.3);
  color: #f8f8f2;
}

.message-meta {
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.message-time {
  font-size: 11px;
  opacity: 0.7;
  text-align: left;
}

.token-usage {
  font-size: 10px;
  opacity: 0.6;
  text-align: left;
  font-style: italic;
}

.loading-dots {
  display: flex;
  gap: 4px;
  padding: 8px 0;
}

.loading-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #9ca3af;
  animation: loading 1.4s infinite ease-in-out;
}

.loading-dots span:nth-child(1) { 
  animation-delay: -0.32s; 
}

.loading-dots span:nth-child(2) { 
  animation-delay: -0.16s; 
}

@keyframes loading {
  0%, 80%, 100% { 
    transform: scale(0.8); 
    opacity: 0.5; 
  }
  40% { 
    transform: scale(1); 
    opacity: 1; 
  }
}
</style>