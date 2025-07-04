<template>
  <div class="messages-container">
    <div 
      v-for="message in messages" 
      :key="message.id" 
      :class="['message-wrapper', message.type]"
    >
      <div class="message-bubble">
        {{ message.content }}
      </div>
    </div>
    
    <div v-if="streamingMessage" class="message-wrapper bot">
      <div class="message-bubble streaming">
        {{ streamingMessage }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MessagesComponent',
  props: {
    messages: {
      type: Array,
      required: true
    },
    streamingMessage: {
      type: String,
      default: ''
    }
  }
}
</script>

<style scoped>
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.message-wrapper {
  display: flex;
  margin-bottom: 10px;
}

.message-wrapper.user {
  justify-content: flex-end;
}

.message-wrapper.bot {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 70%;
  padding: 10px;
  border: 1px solid #ccc;
  word-wrap: break-word;
}

.message-wrapper.user .message-bubble {
  background-color: #e3f2fd;
}

.message-wrapper.bot .message-bubble {
  background-color: #f5f5f5;
}

.message-bubble.streaming {
  border-color: #007bff;
}
</style>