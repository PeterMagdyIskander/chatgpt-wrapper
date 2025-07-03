<template>
  <div class="chat-view">
    <!-- Header -->
    <div class="chat-header">
      <div class="header-left">
        <button class="nav-toggle-btn" @click="toggleNavigator" :class="{ active: showNavigator }">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
          </svg>
          <span class="conversation-count" v-if="conversationCount > 0">{{ conversationCount }}</span>
        </button>
        <h1>{{ currentConversationTitle }}</h1>
      </div>
      <div class="header-right">
        <button class="new-chat-btn" @click="onNewConversation" title="New Conversation">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </button>
        <div class="status-indicator" :class="{ connected: isConnected }">
          <span class="status-dot"></span>
          {{ isConnected ? 'Connected' : 'Disconnected' }}
        </div>
      </div>
    </div>

    <!-- Conversation Navigator -->
    <ConversationNavigator
      :show="showNavigator"
      :conversations="conversations"
      :current-conversation-id="currentConversationId"
      @select-conversation="onSelectConversation"
      @delete-conversation="onDeleteConversation"
      @new-conversation="onNewConversation"
      @close="showNavigator = false"
    />

    <!-- Messages Container -->
    <div class="messages-container">
      <ChatMessages
        ref="chatMessagesComponent"
        :messages="messages"
        :is-loading="isLoading"
      />
    </div>

    <!-- Chat Input at Bottom -->
    <div class="input-container">
      <ChatInput
        ref="chatInputComponent"
        :is-loading="isLoading"
        @send-message="onSendMessage"
      />
    </div>

    <!-- Loading Overlay -->
    <div v-if="isInitializing" class="loading-overlay">
      <div class="loading-spinner"></div>
      <p>Initializing chat...</p>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, nextTick, computed } from 'vue'
import ChatMessages from '../components/ChatMessages.vue'
import ChatInput from '../components/ChatInput.vue'
import ConversationNavigator from '../components/ConversationNavigator.vue'
import { ChatService } from '../services/chat.service.js'

export default {
  name: 'ChatView',
  components: {
    ChatMessages,
    ChatInput,
    ConversationNavigator
  },
  setup() {
    const messages = ref([])
    const isLoading = ref(false)
    const isConnected = ref(false)
    const isInitializing = ref(true)
    const chatMessagesComponent = ref(null)
    const chatInputComponent = ref(null)
    const chatService = ref(null)
    const currentConversationId = ref(null)
    const currentConversationTitle = ref('New Conversation')
    const showNavigator = ref(false)
    const conversations = ref([])

    // Storage keys
    const CONVERSATIONS_KEY = 'ai_chat_conversations'
    const CURRENT_CONVERSATION_KEY = 'ai_chat_current_conversation'

    // Configuration - in a real app, this would come from environment variables or config
    const config = {
      openaiApiKey: process.env.VUE_APP_OPENAI_API_KEY || '', // Set your API key
      openaiApiUrl: 'https://api.openai.com/v1/chat/completions',
      additionalData: null // Can be used to pass context data
    }

    // HTTP client mock - in a real app, use axios or fetch
    const httpClient = {
      async post(url, data, options) {
        const response = await fetch(url, {
          method: 'POST',
          headers: options.headers,
          body: JSON.stringify(data)
        })
        
        if (!response.ok) {
          const error = new Error(`HTTP ${response.status}`)
          error.response = { status: response.status }
          throw error
        }
        
        return { data: await response.json() }
      }
    }

    // Conversation management functions
    const generateConversationId = () => {
      return `conv_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    }

    const loadConversations = () => {
      try {
        const stored = localStorage.getItem(CONVERSATIONS_KEY)
        return stored ? JSON.parse(stored) : []
      } catch (error) {
        console.error('Error loading conversations:', error)
        return []
      }
    }

    const saveConversations = (convs) => {
      try {
        localStorage.setItem(CONVERSATIONS_KEY, JSON.stringify(convs))
        conversations.value = convs
      } catch (error) {
        console.error('Error saving conversations:', error)
      }
    }

    const saveCurrentConversation = () => {
      if (!currentConversationId.value) return

      const convs = loadConversations()
      const conversationIndex = convs.findIndex(c => c.id === currentConversationId.value)
      
      const conversationData = {
        id: currentConversationId.value,
        title: currentConversationTitle.value,
        messages: messages.value,
        updatedAt: new Date().toISOString(),
        createdAt: conversationIndex === -1 ? new Date().toISOString() : (convs[conversationIndex]?.createdAt || new Date().toISOString())
      }

      if (conversationIndex === -1) {
        convs.unshift(conversationData)
      } else {
        convs[conversationIndex] = conversationData
      }

      saveConversations(convs)
      localStorage.setItem(CURRENT_CONVERSATION_KEY, currentConversationId.value)
    }

    const loadConversation = (conversationId) => {
      const convs = loadConversations()
      const conversation = convs.find(c => c.id === conversationId)
      
      if (conversation) {
        currentConversationId.value = conversation.id
        currentConversationTitle.value = conversation.title
        messages.value = conversation.messages || []
        localStorage.setItem(CURRENT_CONVERSATION_KEY, conversationId)
        
        // Scroll to bottom after loading
        nextTick(() => {
          scrollToBottom()
        })
      }
    }

    const createNewConversation = () => {
      // Save current conversation if it has messages
      if (messages.value.length > 1) { // More than just welcome message
        saveCurrentConversation()
      }

      // Create new conversation
      currentConversationId.value = generateConversationId()
      currentConversationTitle.value = 'New Conversation'
      messages.value = []
      
      // Add welcome message
      if (chatService.value) {
        const welcomeMessage = chatService.value.createWelcomeMessage()
        messages.value.push(welcomeMessage)
      }

      showNavigator.value = false
    }

    const deleteConversation = (conversationId) => {
      const convs = loadConversations()
      const updatedConvs = convs.filter(c => c.id !== conversationId)
      saveConversations(updatedConvs)

      // If we're deleting the current conversation, create a new one
      if (currentConversationId.value === conversationId) {
        createNewConversation()
      }
    }

    // Response formatter - can be customized for markdown, etc.
    const responseFormatter = {
      formatToHTML(content) {
        // Simple HTML formatting - can be enhanced with markdown parser
        return content
          .replace(/\n/g, '<br>')
          .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
          .replace(/\*(.*?)\*/g, '<em>$1</em>')
      }
    }

    const initializeChat = async () => {
      try {
        // Initialize chat service
        chatService.value = new ChatService(httpClient)
        
        // Load conversations list
        conversations.value = loadConversations()
        
        // Try to load the last active conversation
        const lastConversationId = localStorage.getItem(CURRENT_CONVERSATION_KEY)
        if (lastConversationId && conversations.value.some(c => c.id === lastConversationId)) {
          loadConversation(lastConversationId)
        } else {
          // Create new conversation
          currentConversationId.value = generateConversationId()
          currentConversationTitle.value = 'New Conversation'
          
          // Add welcome message
          const welcomeMessage = chatService.value.createWelcomeMessage()
          messages.value.push(welcomeMessage)
        }
        
        isConnected.value = true
        
        // Focus input after initialization
        await nextTick()
        if (chatInputComponent.value) {
          chatInputComponent.value.focus()
        }
      } catch (error) {
        console.error('Failed to initialize chat:', error)
        const errorMessage = {
          id: Date.now(),
          content: 'Failed to initialize chat. Please refresh the page and try again.',
          htmlContent: '<p>Failed to initialize chat. Please refresh the page and try again.</p>',
          role: 'assistant',
          timestamp: new Date(),
          metadata: { isError: true }
        }
        messages.value.push(errorMessage)
        isConnected.value = false
      } finally {
        isInitializing.value = false
      }
    }

    const onSendMessage = async (messageContent) => {
      if (!messageContent.trim() || isLoading.value || !chatService.value) {
        return
      }

      try {
        isLoading.value = true

        // Add user message
        const userMessage = chatService.value.createUserMessage(messageContent)
        messages.value.push(userMessage)

        // Scroll to bottom
        await nextTick()
        scrollToBottom()

        // Check rate limit
        await chatService.value.checkRateLimit()

        // Get chat history (exclude the current user message)
        const chatHistory = messages.value.slice(0, -1).map(msg => ({
          role: msg.role,
          content: msg.content
        }))

        // Call AI service
        const aiResponse = await chatService.value.callOpenAI(
          messageContent,
          chatHistory,
          config
        )

        // Create AI response and extract title if present
        let title = null
        let content = aiResponse
        
        if (aiResponse.includes('title=') && aiResponse.includes('$$')) {
          const parts = aiResponse.split('$$')
          const titlePart = parts[0]
          content = parts.slice(1).join('$$')
          
          // Extract title from title="..." format
          const titleMatch = titlePart.match(/title="([^"]*)"/)
          if (titleMatch) {
            title = titleMatch[1]
          }
        }

        // Update conversation title if this is the first response and we have a title
        if (title && messages.value.length === 2) { // Welcome + user message
          currentConversationTitle.value = title
        }

        // Add AI response
        const assistantMessage = chatService.value.createAssistantMessage(content)
        messages.value.push(assistantMessage)

        // Save conversation after each exchange
        saveCurrentConversation()

        // Scroll to bottom
        await nextTick()
        scrollToBottom()

      } catch (error) {
        console.error('Error sending message:', error)
        const errorMessage = chatService.value.createAssistantMessage(
          'Sorry, I encountered an error processing your message. Please try again.',
          true
        )
        messages.value.push(errorMessage)
      } finally {
        isLoading.value = false
        
        // Focus input after response
        await nextTick()
        if (chatInputComponent.value) {
          chatInputComponent.value.focus()
        }
      }
    }

    const scrollToBottom = () => {
      if (chatMessagesComponent.value) {
        chatMessagesComponent.value.setShouldScrollToBottom()
      }
    }

    const clearChat = () => {
      createNewConversation()
    }

    const toggleNavigator = () => {
      showNavigator.value = !showNavigator.value
    }

    const onSelectConversation = (conversationId) => {
      if (conversationId !== currentConversationId.value) {
        // Save current conversation before switching
        if (messages.value.length > 1) {
          saveCurrentConversation()
        }
        loadConversation(conversationId)
      }
      showNavigator.value = false
    }

    const onDeleteConversation = (conversationId) => {
      deleteConversation(conversationId)
    }

    const onNewConversation = () => {
      createNewConversation()
    }

    // Computed property for conversation count
    const conversationCount = computed(() => conversations.value.length)

    // Initialize on mount
    onMounted(() => {
      initializeChat()
    })

    return {
      messages,
      isLoading,
      isConnected,
      isInitializing,
      chatMessagesComponent,
      chatInputComponent,
      currentConversationTitle,
      showNavigator,
      conversations,
      conversationCount,
      onSendMessage,
      clearChat,
      toggleNavigator,
      onSelectConversation,
      onDeleteConversation,
      onNewConversation
    }
  }
}
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f9fafb;
  position: relative;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
  min-width: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.nav-toggle-btn {
  position: relative;
  background: none;
  border: 1px solid #e5e7eb;
  padding: 0.5rem;
  border-radius: 8px;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.nav-toggle-btn:hover {
  background: #f3f4f6;
  color: #374151;
}

.nav-toggle-btn.active {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.conversation-count {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ef4444;
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.125rem 0.375rem;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
  line-height: 1.2;
}

.new-chat-btn {
  background: none;
  border: 1px solid #e5e7eb;
  padding: 0.5rem;
  border-radius: 8px;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.new-chat-btn:hover {
  background: #f3f4f6;
  color: #374151;
  border-color: #d1d5db;
}

.chat-header h1 {
  margin: 0;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
  truncate: true;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  background: #fee2e2;
  color: #dc2626;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.3s ease;
}

.status-indicator.connected {
  background: #dcfce7;
  color: #16a34a;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 2s infinite;
}

.messages-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0;
  margin: 0;
}

.input-container {
  background: white;
  border-top: 1px solid #e5e7eb;
  box-shadow: 0 -1px 3px 0 rgba(0, 0, 0, 0.1);
  padding: 0;
  z-index: 10;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(249, 250, 251, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

.loading-overlay p {
  color: #6b7280;
  font-size: 1rem;
  margin: 0;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Mobile responsive */
@media (max-width: 768px) {
  .chat-header {
    padding: 1rem;
  }

  .header-left {
    gap: 0.5rem;
  }

  .header-right {
    gap: 0.5rem;
  }

  .chat-header h1 {
    font-size: 1rem;
  }

  .status-indicator {
    font-size: 0.75rem;
    padding: 0.25rem 0.75rem;
  }

  .nav-toggle-btn,
  .new-chat-btn {
    padding: 0.375rem;
  }
}

/* Ensure proper layout */
.chat-view :deep(.chat-messages) {
  height: 100%;
}

.chat-view :deep(.chat-input) {
  border-radius: 0;
  border: none;
  border-top: 1px solid #e5e7eb;
}
</style>