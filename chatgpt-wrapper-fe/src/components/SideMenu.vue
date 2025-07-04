<template>
  <div class="side-menu">
    <button @click="$emit('new-chat')" class="new-chat-btn">
      New Chat
    </button>
    
    <div class="chat-list">
      <div 
        v-for="chat in chats" 
        :key="chat.id" 
        @click="$emit('select-chat', chat.id)"
        :class="['chat-item', { active: currentChat === chat.id }]"
      >
        <div class="chat-name">{{ chat.name }}</div>
        <div class="chat-info">{{ chat.messages.length }} messages</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SideMenu',
  props: {
    chats: {
      type: Array,
      required: true
    },
    currentChat: {
      type: String,
      default: null
    }
  },
  emits: ['select-chat', 'new-chat']
}
</script>

<style scoped>
.side-menu {
  width: 250px;
  background: #f5f5f5;
  padding: 10px;
  border-right: 1px solid #ccc;
  display: flex;
  flex-direction: column;
}

.new-chat-btn {
  width: 100%;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
}

.new-chat-btn:hover {
  background: #e9e9e9;
}

.chat-list {
  flex: 1;
  overflow-y: auto;
}

.chat-item {
  padding: 10px;
  cursor: pointer;
  margin-bottom: 5px;
  border: 1px solid #ddd;
  background: white;
}

.chat-item:hover {
  background: #f0f0f0;
}

.chat-item.active {
  background: #e0e0e0;
}

.chat-name {
  font-weight: bold;
  margin-bottom: 5px;
}

.chat-info {
  font-size: 12px;
  color: #666;
}
</style>