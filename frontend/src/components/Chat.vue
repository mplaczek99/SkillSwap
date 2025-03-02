<template>
  <div class="chat-page">
    <div class="container">
      <div class="chat-container">
        <!-- Sidebar with conversations -->
        <div class="chat-sidebar">
          <div class="sidebar-header">
            <h3>Conversations</h3>
            <button
              class="btn btn-outline btn-sm"
              @click="refreshConversations"
            >
              <font-awesome-icon icon="sync" />
            </button>
          </div>

          <!-- Search for users to start a conversation -->
          <div class="chat-search">
            <input
              type="text"
              v-model="searchQuery"
              placeholder="Search users..."
              @input="searchUsers"
            />
            <div
              v-if="showSearchResults && searchResults.length > 0"
              class="search-results"
            >
              <div
                v-for="user in searchResults"
                :key="user.id"
                class="search-result-item"
                @click="startNewConversation(user)"
              >
                <div class="user-avatar">
                  <font-awesome-icon icon="user" v-if="!user.avatar" />
                  <img v-else :src="user.avatar" :alt="user.name" />
                </div>
                <div class="user-info">
                  <div class="user-name">{{ user.name }}</div>
                  <div class="user-email">{{ user.email }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Loading state -->
          <div v-if="loadingConversations" class="chat-loading">
            <font-awesome-icon icon="spinner" class="spin" />
            <p>Loading conversations...</p>
          </div>

          <!-- Conversation list -->
          <div v-else-if="conversations.length > 0" class="conversation-list">
            <div
              v-for="convo in conversations"
              :key="convo.id"
              class="conversation-item"
              :class="{
                active:
                  activeConversation && activeConversation.id === convo.id,
              }"
              @click="loadConversation(convo.id)"
            >
              <div class="conversation-avatar">
                <font-awesome-icon icon="user" v-if="!convo.recipient.avatar" />
                <img
                  v-else
                  :src="convo.recipient.avatar"
                  :alt="convo.recipient.name"
                />
              </div>
              <div class="conversation-details">
                <div class="conversation-name">{{ convo.recipient.name }}</div>
                <div class="conversation-preview">
                  {{ convo.lastMessage.text }}
                </div>
              </div>
              <div class="conversation-meta">
                <div class="conversation-time">
                  {{ formatTime(convo.lastMessage.timestamp) }}
                </div>
                <div v-if="convo.unreadCount > 0" class="unread-badge">
                  {{ convo.unreadCount }}
                </div>
              </div>
            </div>
          </div>

          <!-- Empty state -->
          <div v-else class="empty-state">
            <font-awesome-icon icon="comment-alt" class="empty-icon" />
            <p>No conversations yet</p>
            <p class="empty-hint">
              Search for users to start chatting or wait for someone to message
              you.
            </p>
          </div>
        </div>

        <!-- Main chat area -->
        <div class="chat-main">
          <template v-if="activeConversation">
            <!-- Chat header -->
            <div class="chat-header">
              <div class="chat-user">
                <div class="user-avatar">
                  <font-awesome-icon
                    icon="user"
                    v-if="!activeConversation.recipient.avatar"
                  />
                  <img
                    v-else
                    :src="activeConversation.recipient.avatar"
                    :alt="activeConversation.recipient.name"
                  />
                </div>
                <div class="user-info">
                  <div class="user-name">
                    {{ activeConversation.recipient.name }}
                  </div>
                </div>
              </div>
              <div class="chat-actions">
                <button class="btn btn-outline btn-sm" @click="openUserProfile">
                  <font-awesome-icon icon="user" />
                  <span>View Profile</span>
                </button>
              </div>
            </div>

            <!-- Messages -->
            <div class="messages-container" ref="messagesContainer">
              <div
                class="messages-date-divider"
                v-if="activeConversation.messages.length > 0"
              >
                {{
                  formatMessageDate(activeConversation.messages[0].timestamp)
                }}
              </div>

              <div
                v-for="(message, index) in activeConversation.messages"
                :key="message.id"
                class="message-wrapper"
              >
                <!-- Date divider -->
                <div
                  v-if="shouldShowDateDivider(message, index)"
                  class="messages-date-divider"
                >
                  {{ formatMessageDate(message.timestamp) }}
                </div>

                <!-- Message bubble -->
                <div
                  class="message"
                  :class="{
                    outgoing: message.isOutgoing,
                    incoming: !message.isOutgoing,
                  }"
                >
                  <div class="message-content">
                    <p>{{ message.text }}</p>
                  </div>
                  <div class="message-time">
                    {{ formatMessageTime(message.timestamp) }}
                  </div>
                </div>
              </div>

              <div v-if="sending" class="message-sending">
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>

            <!-- Message input -->
            <form @submit.prevent="sendMessage" class="message-form">
              <textarea
                v-model="newMessage"
                placeholder="Type your message..."
                rows="1"
                @input="autoGrow"
                ref="messageInput"
                @keydown.enter.prevent="handleKeyDown"
              ></textarea>
              <button
                type="submit"
                class="btn btn-primary message-send"
                :disabled="sending || !newMessage.trim()"
              >
                <font-awesome-icon icon="paper-plane" />
              </button>
            </form>
          </template>

          <!-- No conversation selected state -->
          <div v-else class="no-conversation">
            <div class="no-conversation-content">
              <font-awesome-icon icon="comments" class="no-conversation-icon" />
              <h3>Select a conversation</h3>
              <p>
                Choose an existing conversation or search for a user to start
                chatting.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ChatService from "@/services/ChatService";
import { debounce } from "lodash";

export default {
  name: "Chat",
  data() {
    return {
      conversations: [],
      activeConversation: null,
      newMessage: "",
      loadingConversations: false,
      sending: false,
      searchQuery: "",
      searchResults: [],
      showSearchResults: false,
      searchingUsers: false,
      // Mock user search results
      mockUsers: [
        {
          id: 2,
          name: "Alice Smith",
          email: "alice@example.com",
          avatar: null,
        },
        { id: 3, name: "Bob Johnson", email: "bob@example.com", avatar: null },
        {
          id: 4,
          name: "Carol Williams",
          email: "carol@example.com",
          avatar: null,
        },
        {
          id: 5,
          name: "David Brown",
          email: "david@example.com",
          avatar: null,
        },
      ],
    };
  },
  created() {
    this.loadConversations();
    this.searchUsers = debounce(this.performSearch, 300);

    // Set up demo mock message receiver
    this.setupMockMessageReceiver();

    // Check for conversation ID in route query
    if (this.$route.query.conversation) {
      this.loadConversation(this.$route.query.conversation);
    }

    // Check for user ID in route query (to start a new conversation)
    if (this.$route.query.user && this.$route.query.userName) {
      this.startNewConversation({
        id: parseInt(this.$route.query.user),
        name: this.$route.query.userName,
      });
    }
  },
  mounted() {
    // Close search results when clicking outside
    document.addEventListener("click", this.handleOutsideClick);
  },
  beforeUnmount() {
    document.removeEventListener("click", this.handleOutsideClick);
    clearInterval(this.mockMessageInterval);
  },
  methods: {
    async loadConversations() {
      this.loadingConversations = true;
      try {
        this.conversations = await ChatService.getConversations();
      } catch (error) {
        console.error("Failed to load conversations:", error);
      } finally {
        this.loadingConversations = false;
      }
    },

    async loadConversation(conversationId) {
      try {
        this.activeConversation =
          await ChatService.getConversation(conversationId);

        // Update route query parameter
        this.$router.replace({
          query: { ...this.$route.query, conversation: conversationId },
        });

        // Update the unread count in the conversations list
        const convoIndex = this.conversations.findIndex(
          (c) => c.id === conversationId,
        );
        if (convoIndex !== -1) {
          this.conversations[convoIndex].unreadCount = 0;
        }

        // Scroll to bottom of messages
        this.$nextTick(() => {
          this.scrollToBottom();
        });
      } catch (error) {
        console.error("Failed to load conversation:", error);
      }
    },

    async sendMessage() {
      if (!this.newMessage.trim() || !this.activeConversation) return;

      this.sending = true;
      try {
        const message = await ChatService.sendMessage(
          this.activeConversation.id,
          this.newMessage.trim(),
        );

        // Add message to conversation
        this.activeConversation.messages.push(message);

        // Update last message in conversations list
        const convoIndex = this.conversations.findIndex(
          (c) => c.id === this.activeConversation.id,
        );
        if (convoIndex !== -1) {
          this.conversations[convoIndex].lastMessage = {
            text: message.text,
            timestamp: message.timestamp,
          };

          // Move conversation to top of list
          const convo = this.conversations.splice(convoIndex, 1)[0];
          this.conversations.unshift(convo);
        }

        // Clear input
        this.newMessage = "";
        this.$refs.messageInput.style.height = "auto";

        // Scroll to bottom
        this.$nextTick(() => {
          this.scrollToBottom();
        });
      } catch (error) {
        console.error("Failed to send message:", error);
      } finally {
        this.sending = false;
      }
    },

    scrollToBottom() {
      if (this.$refs.messagesContainer) {
        this.$refs.messagesContainer.scrollTop =
          this.$refs.messagesContainer.scrollHeight;
      }
    },

    autoGrow(e) {
      // Auto-resize textarea
      e.target.style.height = "auto";
      e.target.style.height = e.target.scrollHeight + "px";
    },

    handleKeyDown(e) {
      // Send on Enter, allow Shift+Enter for new line
      if (!e.shiftKey) {
        this.sendMessage();
      }
    },

    formatTime(timestamp) {
      const date = new Date(timestamp);
      const now = new Date();
      const yesterday = new Date(now);
      yesterday.setDate(yesterday.getDate() - 1);

      // Today: show time
      if (date.toDateString() === now.toDateString()) {
        return date.toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
        });
      }

      // Yesterday: show "Yesterday"
      if (date.toDateString() === yesterday.toDateString()) {
        return "Yesterday";
      }

      // This week: show day name
      if (now.getTime() - date.getTime() < 7 * 24 * 60 * 60 * 1000) {
        return date.toLocaleDateString([], { weekday: "short" });
      }

      // Older: show date
      return date.toLocaleDateString([], { month: "short", day: "numeric" });
    },

    formatMessageTime(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    },

    formatMessageDate(timestamp) {
      const date = new Date(timestamp);
      const now = new Date();
      const yesterday = new Date(now);
      yesterday.setDate(yesterday.getDate() - 1);

      // Today
      if (date.toDateString() === now.toDateString()) {
        return "Today";
      }

      // Yesterday
      if (date.toDateString() === yesterday.toDateString()) {
        return "Yesterday";
      }

      // Show full date for older messages
      return date.toLocaleDateString([], {
        weekday: "long",
        month: "long",
        day: "numeric",
        year: date.getFullYear() !== now.getFullYear() ? "numeric" : undefined,
      });
    },

    shouldShowDateDivider(message, index) {
      if (index === 0) return false;

      const currentDate = new Date(message.timestamp).toDateString();
      const prevDate = new Date(
        this.activeConversation.messages[index - 1].timestamp,
      ).toDateString();

      return currentDate !== prevDate;
    },

    refreshConversations() {
      this.loadConversations();
    },

    performSearch() {
      if (!this.searchQuery.trim()) {
        this.searchResults = [];
        this.showSearchResults = false;
        return;
      }

      this.searchingUsers = true;

      // In a real app, we would call an API
      // For now, use our mock users
      const query = this.searchQuery.toLowerCase();
      this.searchResults = this.mockUsers.filter(
        (user) =>
          user.name.toLowerCase().includes(query) ||
          user.email.toLowerCase().includes(query),
      );

      this.showSearchResults = true;
      this.searchingUsers = false;
    },

    async startNewConversation(user) {
      // Hide search results
      this.showSearchResults = false;
      this.searchQuery = "";

      try {
        // Start new conversation (or get existing one)
        const conversationId = await ChatService.startConversation(
          user.id,
          user.name,
        );

        // Refresh conversations list
        await this.loadConversations();

        // Load the conversation
        await this.loadConversation(conversationId);

        // Focus message input
        this.$nextTick(() => {
          this.$refs.messageInput.focus();
        });
      } catch (error) {
        console.error("Failed to start conversation:", error);
      }
    },

    openUserProfile() {
      if (!this.activeConversation) return;

      // In a real app, navigate to user profile
      alert(`Viewing profile for ${this.activeConversation.recipient.name}`);
    },

    handleOutsideClick(event) {
      // Close search results when clicking outside
      if (this.showSearchResults && !event.target.closest(".chat-search")) {
        this.showSearchResults = false;
      }
    },

    // For demo purposes - simulate receiving messages
    setupMockMessageReceiver() {
      // Every 45-90 seconds, simulate receiving a message in one of the conversations
      this.mockMessageInterval = setInterval(
        async () => {
          if (this.conversations.length === 0) return;

          // Pick a random conversation
          const randomIndex = Math.floor(
            Math.random() * this.conversations.length,
          );
          const conversation = this.conversations[randomIndex];

          // Get a random message
          const messages = [
            "Hello there! How's your day going?",
            "Just checking in. Have you had a chance to practice?",
            "I have some free time next week if you'd like to schedule a session.",
            "I found a great resource that might help you learn faster!",
            "Quick question about our last session...",
            "Are you still interested in learning this skill?",
            "I enjoyed our last skill exchange session! Hope we can do it again soon.",
            "Do you think you could teach me more about advanced techniques?",
            "I've been practicing what you taught me last time. Making progress!",
          ];
          const randomMessage =
            messages[Math.floor(Math.random() * messages.length)];

          // Simulate receiving a message
          const message = await ChatService.simulateIncomingMessage(
            conversation.id,
            randomMessage,
          );

          // Update UI
          this.handleIncomingMessage(message);
        },
        Math.floor(Math.random() * 45000) + 45000,
      ); // Random between 45-90 seconds
    },

    handleIncomingMessage(message) {
      // Update active conversation if it's the one receiving the message
      if (
        this.activeConversation &&
        this.activeConversation.id === message.conversationId
      ) {
        this.activeConversation.messages.push({
          ...message,
          isOutgoing: false,
        });

        // Scroll to bottom
        this.$nextTick(() => {
          this.scrollToBottom();
        });
      } else {
        // Update conversation in list
        const convoIndex = this.conversations.findIndex(
          (c) => c.id === message.conversationId,
        );
        if (convoIndex !== -1) {
          const convo = this.conversations[convoIndex];
          convo.lastMessage = {
            text: message.text,
            timestamp: message.timestamp,
          };
          convo.unreadCount += 1;

          // Move to top of list
          this.conversations.splice(convoIndex, 1);
          this.conversations.unshift(convo);
        }

        // Show notification
        this.$root.$emit("show-notification", {
          type: "message",
          title: `New message from ${this.conversations[0].recipient.name}`,
          message: message.text,
          duration: 5000,
        });
      }
    },
  },
};
</script>

<style scoped>
.chat-page {
  padding: var(--space-4) 0;
}

.chat-container {
  display: flex;
  height: calc(100vh - 200px);
  min-height: 500px;
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}

/* Chat Sidebar */
.chat-sidebar {
  width: 320px;
  border-right: 1px solid var(--light);
  display: flex;
  flex-direction: column;
  background-color: var(--white);
}

.sidebar-header {
  padding: var(--space-4);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--light);
}

.sidebar-header h3 {
  margin: 0;
  font-size: var(--font-size-lg);
}

.chat-search {
  padding: var(--space-3);
  position: relative;
}

.chat-search input {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
}

.search-results {
  position: absolute;
  top: 100%;
  left: var(--space-3);
  right: var(--space-3);
  background-color: var(--white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  z-index: 10;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid var(--light);
}

.search-result-item {
  display: flex;
  align-items: center;
  padding: var(--space-3);
  cursor: pointer;
  transition: background-color var(--transition-fast) ease;
}

.search-result-item:hover {
  background-color: var(--light);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: var(--primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary-color);
  font-size: var(--font-size-lg);
  margin-right: var(--space-3);
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-weight: var(--font-weight-medium);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: var(--font-size-xs);
  color: var(--medium);
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
}

.conversation-item {
  display: flex;
  padding: var(--space-3);
  cursor: pointer;
  border-bottom: 1px solid var(--light);
  transition: background-color var(--transition-fast) ease;
}

.conversation-item:hover,
.conversation-item.active {
  background-color: var(--primary-light);
}

.conversation-details {
  flex: 1;
  min-width: 0;
  margin-right: var(--space-2);
}

.conversation-name {
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--space-1);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conversation-preview {
  font-size: var(--font-size-xs);
  color: var(--medium);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 170px;
}

.conversation-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  font-size: var(--font-size-xs);
}

.conversation-time {
  color: var(--medium);
  margin-bottom: var(--space-1);
}

.unread-badge {
  background-color: var(--primary-color);
  color: white;
  font-size: 10px;
  font-weight: var(--font-weight-bold);
  border-radius: 50%;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-loading,
.empty-state,
.no-conversation {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-6);
  color: var(--medium);
  text-align: center;
  height: 100%;
}

.empty-icon,
.no-conversation-icon {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
  opacity: 0.5;
}

.empty-hint {
  font-size: var(--font-size-sm);
  margin-top: var(--space-2);
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* Chat Main Area */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: var(--light);
}

.chat-header {
  padding: var(--space-4);
  border-bottom: 1px solid var(--light);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: var(--white);
}

.chat-user {
  display: flex;
  align-items: center;
}

.chat-actions {
  display: flex;
  gap: var(--space-2);
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
}

.messages-date-divider {
  text-align: center;
  margin: var(--space-4) 0;
  position: relative;
  color: var(--medium);
  font-size: var(--font-size-xs);
}

.messages-date-divider::before,
.messages-date-divider::after {
  content: "";
  position: absolute;
  top: 50%;
  width: 30%;
  height: 1px;
  background-color: var(--light);
}

.messages-date-divider::before {
  left: 0;
}

.messages-date-divider::after {
  right: 0;
}

.message-wrapper {
  margin-bottom: var(--space-3);
}

.message {
  max-width: 80%;
  padding: var(--space-3);
  border-radius: var(--radius-lg);
  position: relative;
  margin-bottom: var(--space-1);
}

.message.outgoing {
  align-self: flex-end;
  background-color: var(--primary-color);
  color: white;
  border-bottom-right-radius: var(--radius-sm);
  margin-left: auto;
}

.message.incoming {
  align-self: flex-start;
  background-color: var(--white);
  border-bottom-left-radius: var(--radius-sm);
}

.message-content p {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.message-time {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.7);
  margin-top: var(--space-1);
  text-align: right;
}

.message.incoming .message-time {
  color: var(--medium);
}

.message-sending {
  display: flex;
  justify-content: flex-start;
  margin-bottom: var(--space-3);
}

.typing-indicator {
  background-color: var(--white);
  padding: 10px 15px;
  border-radius: 1.5rem;
  border-bottom-left-radius: var(--radius-sm);
  display: flex;
  align-items: center;
}

.typing-indicator span {
  height: 8px;
  width: 8px;
  background-color: var(--medium);
  border-radius: 50%;
  display: inline-block;
  margin: 0 1px;
  opacity: 0.4;
  animation: typing 1.2s infinite;
}

.typing-indicator span:nth-child(1) {
  animation-delay: 0s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0% {
    opacity: 0.4;
    transform: translateY(0);
  }
  50% {
    opacity: 1;
    transform: translateY(-5px);
  }
  100% {
    opacity: 0.4;
    transform: translateY(0);
  }
}

.message-form {
  display: flex;
  padding: var(--space-3);
  background-color: var(--white);
  border-top: 1px solid var(--light);
  align-items: flex-end;
}

.message-form textarea {
  flex: 1;
  background-color: var(--light);
  border: none;
  border-radius: var(--radius-lg);
  padding: var(--space-3);
  resize: none;
  max-height: 150px;
  font-family: inherit;
  font-size: var(--font-size-md);
  margin-right: var(--space-2);
}

.message-form textarea:focus {
  outline: none;
}

.message-send {
  align-self: stretch;
  padding: 0 var(--space-4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.no-conversation-content {
  max-width: 300px;
}

/* Responsive */
@media (max-width: 768px) {
  .chat-container {
    flex-direction: column;
    height: calc(100vh - 180px);
  }

  .chat-sidebar {
    width: 100%;
    height: 300px;
    border-right: none;
    border-bottom: 1px solid var(--light);
  }

  .chat-main {
    height: calc(100vh - 500px);
  }
}
</style>
