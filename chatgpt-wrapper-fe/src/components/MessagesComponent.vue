<template>
  <div class="messages-container">
    <div v-for="message in messages" :key="message.id" :class="['message-wrapper', message.type]">
      <div class="message-bubble">
        {{ message.content }}
      </div>
    </div>

    <div v-if="streamingMessage" class="message-wrapper bot">
      <div class="message-bubble">
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
  padding: 40px 16px;
  flex: 1;
  overflow-y: auto;
   scrollbar-width: thin;
    scrollbar-color: rgba(242, 235, 250,0.3) transparent;
}
.message-wrapper {
  display: flex;
  margin: 20px 0;
  color: rgb(242, 235, 250);
}

.message-wrapper.user {
  justify-content: flex-end;
}

.message-wrapper.bot {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 100%;
  word-wrap: break-word;
}

.message-wrapper.user .message-bubble {
  background-color: rgba(54, 45, 61, 0.5);
  padding: 12px 16px;
  border: 1px solid rgba(54, 45, 61, 0.5);
  border-radius: 12px;
}
</style>