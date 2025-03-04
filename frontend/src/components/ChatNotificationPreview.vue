<template>
  <transition name="message-pop">
    <div
      v-if="activeMessage"
      class="message-preview-container"
      @click="viewConversation"
    >
      <div class="message-preview">
        <div class="message-avatar">
          <font-awesome-icon v-if="!activeMessage.senderAvatar" icon="user" />
          <img
            v-else
            :src="activeMessage.senderAvatar"
            :alt="activeMessage.senderName"
          />
        </div>
        <div class="message-content">
          <div class="message-header">
            <span class="message-sender">{{ activeMessage.senderName }}</span>
            <span class="message-time">{{
              formatTime(activeMessage.timestamp)
            }}</span>
          </div>
          <div class="message-text">{{ activeMessage.text }}</div>
        </div>
        <button class="message-close" @click.stop="closeMessage">
          <font-awesome-icon icon="times" />
        </button>
      </div>
    </div>
  </transition>
</template>

<script>
import eventBus from "@/utils/eventBus";

export default {
  name: "ChatNotificationPreview",
  data() {
    return {
      activeMessage: null,
      messageQueue: [],
      messageTimeout: null,
    };
  },
  created() {
    // Listen for incoming message events using the event bus.
    eventBus.on("chat:incoming-message", this.handleIncomingMessage);
  },
  beforeUnmount() {
    // Clean up the event listener and any pending timeout.
    eventBus.off("chat:incoming-message", this.handleIncomingMessage);
    this.clearTimeouts();
  },
  methods: {
    handleIncomingMessage(message) {
      // If currently in the chat conversation for the message, ignore.
      if (
        this.$route.name === "Chat" &&
        this.$route.query.conversation === message.conversationId.toString()
      ) {
        return;
      }
      this.messageQueue.push(message);
      if (!this.activeMessage) {
        this.showNextMessage();
      }
    },
    showNextMessage() {
      if (this.messageQueue.length === 0) {
        this.activeMessage = null;
        return;
      }
      this.activeMessage = this.messageQueue.shift();
      // Auto-dismiss after 5 seconds.
      this.clearTimeouts();
      this.messageTimeout = setTimeout(() => {
        this.closeMessage();
      }, 5000);
    },
    closeMessage() {
      this.activeMessage = null;
      this.clearTimeouts();
      if (this.messageQueue.length > 0) {
        setTimeout(() => {
          this.showNextMessage();
        }, 300);
      }
    },
    viewConversation() {
      if (!this.activeMessage) return;
      this.$router.push({
        name: "Chat",
        query: { conversation: this.activeMessage.conversationId },
      });
      this.closeMessage();
    },
    formatTime(timestamp) {
      return new Date(timestamp).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    },
    clearTimeouts() {
      if (this.messageTimeout) {
        clearTimeout(this.messageTimeout);
        this.messageTimeout = null;
      }
    },
  },
};
</script>

<style scoped>
/* (CSS remains unchanged) */
.message-preview-container {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: var(--z-toast);
  max-width: 350px;
  width: 100%;
}
.message-preview {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  padding: var(--space-3);
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  overflow: hidden;
  cursor: pointer;
  border-left: 4px solid var(--primary-color);
  transition: transform var(--transition-fast) ease;
}
.message-preview:hover {
  transform: translateY(-3px);
}
.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: var(--primary-light);
  color: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-lg);
  flex-shrink: 0;
}
.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}
.message-content {
  flex: 1;
  min-width: 0;
}
.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-1);
}
.message-sender {
  font-weight: var(--font-weight-semibold);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.message-time {
  font-size: var(--font-size-xs);
  color: var(--medium);
  white-space: nowrap;
  margin-left: var(--space-2);
}
.message-text {
  font-size: var(--font-size-sm);
  color: var(--medium);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.message-close {
  background: transparent;
  border: none;
  color: var(--medium);
  font-size: var(--font-size-sm);
  cursor: pointer;
  opacity: 0.7;
  padding: var(--space-1);
  margin-left: var(--space-1);
}
.message-close:hover {
  opacity: 1;
  color: var(--dark);
}
.message-pop-enter-active,
.message-pop-leave-active {
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
.message-pop-enter-from,
.message-pop-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.9);
}
</style>
