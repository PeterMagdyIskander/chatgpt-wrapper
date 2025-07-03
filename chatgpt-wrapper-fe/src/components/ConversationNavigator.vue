<template>
  <div class="navigator-container" :class="{ show: show }">
    <!-- Overlay -->
    <div class="navigator-overlay" @click="$emit('close')" v-if="show"></div>
    
    <!-- Navigator Panel -->
    <div class="navigator-panel" :class="{ show: show }">
      <div class="navigator-header">
        <h3>Conversations</h3>
        <div class="header-actions">
          <button class="new-conversation-btn" @click="$emit('new-conversation')" title="New Conversation">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
          </button>
          <button class="close-btn" @click="$emit('close')">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>

      <div class="conversations-list">
        <div v-if="conversations.length === 0" class="empty-state">
          <div class="empty-icon">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
          </div>
          <p>No conversations yet</p>
          <button class="start-chat-btn" @click="$emit('new-conversation')">
            Start Your First Chat
          </button>
        </div>

        <div v-else>
          <div class="conversations-header">
            <span class="conversation-count">{{ conversations.length }} conversation{{ conversations.length !== 1 ? 's' : '' }}</span>
          </div>
          
          <div 
            v-for="conversation in sortedConversations" 
            :key="conversation.id"
            class="conversation-item"
            :class="{ active: conversation.id === currentConversationId }"
            @click="$emit('select-conversation', conversation.id)"
          >
            <div class="conversation-content">
              <div class="conversation-title">
                {{ conversation.title || 'Untitled Conversation' }}
              </div>
              <div class="conversation-meta">
                <span class="message-count">{{ getMessageCount(conversation) }} messages</span>
                <span class="conversation-date">{{ formatDate(conversation.updatedAt) }}</span>
              </div>
              <div class="conversation-preview" v-if="getLastMessage(conversation)">
                {{ getLastMessage(conversation) }}
              </div>
            </div>
            
            <div class="conversation-actions">
              <button 
                class="delete-btn" 
                @click.stop="handleDeleteConversation(conversation.id)"
                :title="'Delete conversation: ' + conversation.title"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"></polyline>
                  <path d="m19,6v14a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6m3,0V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2v2"></path>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'ConversationNavigator',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    conversations: {
      type: Array,
      default: () => []
    },
    currentConversationId: {
      type: String,
      default: null
    }
  },
  emits: ['close', 'select-conversation', 'delete-conversation', 'new-conversation'],
  setup(props, { emit }) {
    const sortedConversations = computed(() => {
      return [...props.conversations].sort((a, b) => {
        return new Date(b.updatedAt) - new Date(a.updatedAt)
      })
    })

    const formatDate = (dateString) => {
      const date = new Date(dateString)
      const now = new Date()
      const diffInHours = (now - date) / (1000 * 60 * 60)
      
      if (diffInHours < 1) {
        return 'Just now'
      } else if (diffInHours < 24) {
        return `${Math.floor(diffInHours)}h ago`
      } else if (diffInHours < 24 * 7) {
        return `${Math.floor(diffInHours / 24)}d ago`
      } else {
        return date.toLocaleDateString()
      }
    }

    const getMessageCount = (conversation) => {
      return conversation.messages ? conversation.messages.length : 0
    }

    const getLastMessage = (conversation) => {
      if (!conversation.messages || conversation.messages.length === 0) {
        return ''
      }
      
      const lastMessage = conversation.messages[conversation.messages.length - 1]
      if (lastMessage.role === 'user') {
        return lastMessage.content.substring(0, 100) + (lastMessage.content.length > 100 ? '...' : '')
      }
      return ''
    }

    const handleDeleteConversation = (conversationId) => {
      if (confirm('Are you sure you want to delete this conversation?')) {
        emit('delete-conversation', conversationId)
      }
    }

    return {
      sortedConversations,
      formatDate,
      getMessageCount,
      getLastMessage,
      handleDeleteConversation
    }
  }
}
</script>

<style scoped>
.navigator-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1003;
  pointer-events: none;
}

.navigator-container.show {
  pointer-events: auto;
}

.navigator-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.navigator-container.show .navigator-overlay {
  opacity: 1;
}

.navigator-panel {
  position: absolute;
  top: 0;
  left: -400px;
  width: 400px;
  height: 100%;
  background: white;
  box-shadow: 2px 0 20px rgba(0, 0, 0, 0.1);
  transition: left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

.navigator-panel.show {
  left: 0;
}

.navigator-header {
  display: flex;
}
</style>