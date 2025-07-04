<template>
  <div class="side-menu">
    <h1 class="title">ChatGPT Wrapper</h1>
    <button @click="$emit('new-chat')" class="new-chat-btn">
      New Chat
    </button>

    <div class="chat-list">
      <div v-for="chat in chats" :key="chat.id" @click="$emit('select-chat', chat.id)"
        :class="['chat-item', { active: currentChat === chat.id }]">
        <div class="chat-name">{{ chat.name }}</div>
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
  background-color: #171017;
  padding: 10px;
  display: flex;
  flex-direction: column;
  row-gap: 8px;
}

.title {
  font-size: 18px;
  color: rgb(227, 186, 209);
  text-align: center;
}

.new-chat-btn {
  width: 100%;
  border: 1px solid rgb(39, 36, 44);
  background: rgba(163, 0, 76, 0.2);
  color: rgb(251, 208, 232);
  border-radius: 8px;
  padding: 8px 16px;
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
  cursor: pointer;
  color: rgb(244, 244, 245);
  font-size: 14px;
  
  padding: 8px 16px;
  border-radius: 6px;
}


.chat-item.active {
  background: rgb(38, 25, 34);
}
</style>