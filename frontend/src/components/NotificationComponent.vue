<template>
  <transition name="message-pop">
    <div
      v-if="activeMessage"
      class="message-preview-container"
      @click="viewConversation"
    >
      <div class="message-preview" :class="getNotificationClass">
        <div class="message-avatar">
          <font-awesome-icon
            v-if="!activeMessage.senderAvatar"
            :icon="getNotificationIcon"
          />
          <img
            v-else
            :src="activeMessage.senderAvatar"
            :alt="activeMessage.senderName"
          />
        </div>
        <div class="message-content">
          <div class="message-header">
            <span class="message-sender">{{ activeMessage.senderName }}</span>
            <span class="message-time">{{ formattedTime }}</span>
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
import {
  defineComponent,
  ref,
  computed,
  onMounted,
  onBeforeUnmount,
} from "vue";
import eventBus from "@/utils/eventBus";
import { useRouter } from "vue-router";

// Precompute formatter function outside component
const createTimeFormatter = () => {
  const formatter = new Intl.DateTimeFormat(undefined, {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (timestamp) => {
    if (!timestamp) return "";
    return formatter.format(new Date(timestamp));
  };
};

export default defineComponent({
  name: "NotificationComponent",
  setup() {
    const router = useRouter();

    // State with refs for better performance
    const activeMessage = ref(null);
    const messageQueue = ref([]);

    // Use compute property for formatted time to avoid recalculating on every render
    const formattedTime = computed(() => {
      if (!activeMessage.value || !activeMessage.value.timestamp) return "";
      return formatTime(activeMessage.value.timestamp);
    });

    // Compute properties for notification styling and icons
    const getNotificationClass = computed(() => {
      if (!activeMessage.value) return "";

      // If it's a chat notification, use default styling
      if (activeMessage.value.conversationId) return "";

      // For other notification types, add a specific class
      return `notification-${activeMessage.value.type || "info"}`;
    });

    const getNotificationIcon = computed(() => {
      if (!activeMessage.value) return "user";

      // For chat notifications, use user icon
      if (activeMessage.value.conversationId) return "user";

      // For other notification types, use appropriate icons
      const iconMap = {
        success: "check-circle",
        error: "exclamation-circle",
        warning: "exclamation-triangle",
        info: "info-circle",
      };

      return iconMap[activeMessage.value.type] || "bell";
    });

    // Precomputed time formatter
    const formatTime = createTimeFormatter();

    // Scheduler for better batching and performance
    const scheduler = {
      callbacks: new Map(),
      timeouts: new Map(),

      // Schedule a callback with specified delay
      schedule(callbackId, callback, delay) {
        // Clear existing timeout
        this.clear(callbackId);

        // Store the callback
        this.callbacks.set(callbackId, callback);

        // Create new timeout
        const timeoutId = setTimeout(() => {
          this.execute(callbackId);
        }, delay);

        this.timeouts.set(callbackId, timeoutId);

        return timeoutId;
      },

      // Execute callback by id
      execute(callbackId) {
        const callback = this.callbacks.get(callbackId);
        if (callback) {
          callback();
          this.callbacks.delete(callbackId);
          this.timeouts.delete(callbackId);
        }
      },

      // Clear timeout by id
      clear(callbackId) {
        if (this.timeouts.has(callbackId)) {
          clearTimeout(this.timeouts.get(callbackId));
          this.timeouts.delete(callbackId);
        }
      },

      // Clear all timeouts
      clearAll() {
        for (const [, timeoutId] of this.timeouts) {
          clearTimeout(timeoutId);
        }
        this.callbacks.clear();
        this.timeouts.clear();
      },
    };

    // Handle general notifications
    const handleNotification = (notification) => {
      // Transform general notification to the message format
      const message = {
        senderName: notification.title || "Notification",
        text: notification.message,
        timestamp: new Date(),
        type: notification.type || "info",
        duration: notification.duration || 5000,
      };

      // Add to queue and show if no active message
      messageQueue.value.push(message);
      if (!activeMessage.value) {
        showNextMessage();
      }
    };

    // Event handlers for chat notifications
    const handleIncomingMessage = (message) => {
      // Skip if user is already viewing this conversation
      if (
        window.location.href.includes(`conversation=${message.conversationId}`)
      ) {
        return;
      }

      // Using push with direct array access for better performance than reactive arrays
      messageQueue.value.push(message);

      // Show only if no active message
      if (!activeMessage.value) {
        showNextMessage();
      }
    };

    const handleReadMessages = () => {
      // Clear active message and queue
      activeMessage.value = null;
      messageQueue.value.length = 0; // Faster than reassigning
      scheduler.clearAll();
    };

    const showNextMessage = () => {
      if (messageQueue.value.length === 0) {
        activeMessage.value = null;
        return;
      }

      // Shift is faster than splice(0, 1) and maintains reference
      activeMessage.value = messageQueue.value.shift();

      // Auto-dismiss after specified duration (default 5 seconds)
      scheduler.clear("dismissMessage");
      const duration = activeMessage.value.duration || 5000;
      scheduler.schedule("dismissMessage", closeMessage, duration);
    };

    const closeMessage = () => {
      activeMessage.value = null;
      scheduler.clear("dismissMessage");

      // Show next message after a short delay if queue isn't empty
      if (messageQueue.value.length > 0) {
        scheduler.schedule("showNext", showNextMessage, 300);
      }
    };

    const viewConversation = () => {
      if (!activeMessage.value) return;

      // Only navigate for chat notifications
      if (activeMessage.value.conversationId) {
        router.push({
          name: "Chat",
          query: { conversation: activeMessage.value.conversationId },
        });
      }

      closeMessage();
    };

    // Lifecycle hooks
    onMounted(() => {
      // Set up all event listeners
      eventBus.on("show-notification", handleNotification);
      eventBus.on("chat:incoming-message", handleIncomingMessage);
      eventBus.on("chat:read-messages", handleReadMessages);
    });

    onBeforeUnmount(() => {
      // Proper cleanup of all event listeners
      eventBus.off("show-notification", handleNotification);
      eventBus.off("chat:incoming-message", handleIncomingMessage);
      eventBus.off("chat:read-messages", handleReadMessages);
      scheduler.clearAll();
    });

    // Return all necessary functions and state
    return {
      activeMessage,
      messageQueue,
      formattedTime,
      getNotificationClass,
      getNotificationIcon,
      closeMessage,
      viewConversation,
      formatTime,
    };
  },
});
</script>

<style scoped>
/* Scoped styles for better CSS performance */
.message-preview-container {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: var(--z-toast);
  max-width: 350px;
  width: 100%;
  will-change: transform, opacity;
  /* Hints for browser optimization */
  contain: content;
  /* Containment for better rendering performance */
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
  transition:
    transform 150ms ease,
    box-shadow 150ms ease;
  /* Fixed values over variables for better performance */
  backface-visibility: hidden;
  /* Prevents rendering artifacts */
  transform: translateZ(0);
  /* Forces GPU acceleration */
}

/* Notification type styling */
.message-preview.notification-success {
  border-left-color: var(--success-color);
}

.message-preview.notification-error {
  border-left-color: var(--error-color);
}

.message-preview.notification-warning {
  border-left-color: var(--warning-color);
}

.message-preview.notification-info {
  border-left-color: var(--info-color);
}

.message-preview:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
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
  font-size: 1.25rem;
  /* Fixed value for better rendering */
  flex-shrink: 0;
}

/* Avatar colors for notification types */
.notification-success .message-avatar {
  background-color: var(--success-color);
  color: white;
}

.notification-error .message-avatar {
  background-color: var(--error-color);
  color: white;
}

.notification-warning .message-avatar {
  background-color: var(--warning-color);
  color: white;
}

.notification-info .message-avatar {
  background-color: var(--info-color);
  color: white;
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
  margin-bottom: 0.25rem;
  /* Fixed value for better performance */
}

.message-sender {
  font-weight: 600;
  /* Fixed value over CSS var */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.message-time {
  font-size: 0.75rem;
  /* Fixed value over CSS var */
  color: var(--medium);
  white-space: nowrap;
  margin-left: 0.5rem;
  /* Fixed value over CSS var */
}

.message-text {
  font-size: 0.875rem;
  /* Fixed value over CSS var */
  color: var(--medium);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  /* Standard property for compatibility */
  -webkit-box-orient: vertical;
}

.message-close {
  background: transparent;
  border: none;
  color: var(--medium);
  font-size: 0.875rem;
  /* Fixed value over CSS var */
  cursor: pointer;
  opacity: 0.7;
  padding: 0.25rem;
  /* Fixed value over CSS var */
  margin-left: 0.25rem;
  /* Fixed value over CSS var */
  width: 24px;
  /* Fixed width for better click target */
  height: 24px;
  /* Fixed height for better click target */
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 150ms ease;
}

.message-close:hover {
  opacity: 1;
  color: var(--dark);
  background-color: rgba(0, 0, 0, 0.05);
}

/* Optimized animation with hardware acceleration */
.message-pop-enter-active,
.message-pop-leave-active {
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  transform: translateZ(0);
  /* Force GPU acceleration */
}

.message-pop-enter-from,
.message-pop-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.9) translateZ(0);
}
</style>
