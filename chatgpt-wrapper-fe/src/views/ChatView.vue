<template>
  <div class="chat-container">
    <SideMenu :chats="chats" :current-chat="currentChat" @select-chat="selectChat" @new-chat="createNewChat" />

    <div class="chat-main">
      <div class="content-container">
        <MessagesComponent :messages="currentChatMessages" :streaming-message="streamingMessage" />

        <MessageInput @send-message="sendMessage" />
      </div>
    </div>

  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SideMenu from '../components/SideMenu.vue'
import MessagesComponent from '../components/MessagesComponent.vue'
import MessageInput from '../components/MessageInput.vue'
import { StorageService } from '../services/StorageService.js'
import { ChatService } from '../services/ChatService.js'
import { generateUserId, generateChatId } from '../utils/helpers.js'

export default {
  name: 'ChatView',
  components: {
    SideMenu,
    MessagesComponent,
    MessageInput
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    const userId = ref('')
    const chats = ref([])
    const currentChat = ref(null)
    const streamingMessage = ref('')

    // Initialize app
    onMounted(() => {
      initializeUserId()
      loadChats()

      // Set current chat from route params
      if (route.params.id) {
        currentChat.value = route.params.id
      }
    })

    // Watch for route changes
    watch(() => route.params.id, (newId) => {
      if (newId && chats.value.find(c => c.id === newId)) {
        currentChat.value = newId
        streamingMessage.value = ''
      }
    })

    // Initialize user ID
    function initializeUserId() {
      let storedUserId = StorageService.getItem('userId')
      if (!storedUserId) {
        storedUserId = generateUserId()
        StorageService.setItem('userId', storedUserId)
      }
      userId.value = storedUserId
    }

    // Load chats from localStorage
    function loadChats() {
      const storedChats = StorageService.getItem('chats')
      if (storedChats) {
        try {
          chats.value = JSON.parse(storedChats)
          if (chats.value.length > 0 && !currentChat.value) {
            currentChat.value = chats.value[0].id
            router.push(`/chat/${currentChat.value}`)
          } else if (chats.value.length === 0) {
            createNewChat()
          }
        } catch (e) {
          console.error('Error parsing stored chats:', e)
          createNewChat()
        }
      } else {
        createNewChat()
      }
    }

    // Save chats to localStorage
    function saveChats() {
      try {
        StorageService.setItem('chats', JSON.stringify(chats.value))
      } catch (e) {
        console.error('Error saving chats:', e)
      }
    }

    // Get current chat messages
    const currentChatMessages = computed(() => {
      const chat = chats.value.find(c => c.id === currentChat.value)
      return chat ? chat.messages : []
    })

    // Create new chat
    function createNewChat() {
      const newChat = {
        id: generateChatId(),
        name: `Chat ${chats.value.length + 1}`,
        messages: [],
        createdAt: new Date().toISOString()
      }
      chats.value.push(newChat)
      currentChat.value = newChat.id
      saveChats()
      router.push(`/chat/${currentChat.value}`)
    }

    // Select chat
    function selectChat(chatId) {
      currentChat.value = chatId
      streamingMessage.value = ''
      router.push(`/chat/${chatId}`)
    }

    // Add message to current chat
    function addMessage(content, type) {
      const chat = chats.value.find(c => c.id === currentChat.value)
      if (chat) {
        const message = {
          id: Date.now() + Math.random(),
          content,
          type,
          timestamp: new Date().toISOString()
        }
        chat.messages.push(message)

        // Update chat name based on first user message
        if (type === 'user' && chat.messages.filter(m => m.type === 'user').length === 1) {
          chat.name = content.substring(0, 30) + (content.length > 30 ? '...' : '')
        }

        saveChats()
      }
    }

    // Send message
    async function sendMessage(message) {
      addMessage(message, 'user')

      try {
        const response = await ChatService.sendMessage(message, userId.value)

        if (response.status === 'approved') {
          // Start SSE connection
          const eventSource = ChatService.createSSEConnection(userId.value, response.messageId)

          // Reset streaming message
          streamingMessage.value = ''

          // Listen for connection event
          eventSource.addEventListener('connection', function (event) {
            console.log('Connection established:', event.data)
          })

          // Listen for data events (where the actual content comes through)
          eventSource.addEventListener('data', function (event) {
            console.log('Received data chunk:', event.data)

            // Accumulate the streaming message in real-time
            streamingMessage.value += event.data
          })

          // Listen for done event to finalize the message
          eventSource.addEventListener('done', function (event) {
            console.log('Stream completed:', event.data)

            // Add the final message to chat history
            if (streamingMessage.value.trim()) {
              addMessage(streamingMessage.value, 'bot')
            }

            // Reset streaming message and close connection
            streamingMessage.value = ''
            eventSource.close()
          })

          // Handle errors
          eventSource.onerror = function (event) {
            console.error('SSE error:', event)
            console.log('EventSource readyState:', eventSource.readyState)

            // If we have partial content, save it
            if (streamingMessage.value.trim()) {
              addMessage(streamingMessage.value, 'bot')
              streamingMessage.value = ''
            }

            eventSource.close()
          }

          // Connection opened
          eventSource.onopen = function (event) {
            console.log('SSE connection opened')
          }

        } else {
          addMessage('Message was not approved', 'bot')
        }
      } catch (error) {
        console.error('Error sending message:', error)
        addMessage('Error: Could not connect to server', 'bot')
      }
    }

    return {
      userId,
      chats,
      currentChat,
      currentChatMessages,
      streamingMessage,
      sendMessage,
      selectChat,
      createNewChat
    }
  }
}
</script>

<style scoped>
.chat-container {
  height: 100%;
  display: flex;
}

.chat-main {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  flex: 1;
  margin-top: 15px;
  margin-left: 5px;
  background-color: rgb(31, 26, 36);
  border: 1px solid rgb(50, 32, 40);
  border-top-left-radius: 15px;
}
.content-container{
    width: 720px;
    height: 100%;
    display: flex;
    flex-direction: column;
    position: relative;
}
</style>