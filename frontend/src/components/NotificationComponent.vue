<template>
  <div class="notification-container">
    <transition-group name="notification-slide">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="notification-toast"
        :class="notification.type"
      >
        <div class="notification-icon">
          <font-awesome-icon :icon="getIcon(notification.type)" />
        </div>
        <div class="notification-content">
          <div class="notification-title">{{ notification.title }}</div>
          <div class="notification-message">{{ notification.message }}</div>
        </div>
        <button
          class="notification-close"
          @click="removeNotification(notification.id)"
        >
          <font-awesome-icon icon="times" />
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script>
import eventBus from "@/utils/eventBus";

export default {
  name: "NotificationComponent",
  data() {
    return {
      notifications: [],
      counter: 0,
    };
  },
  created() {
    // Listen for notification events using eventBus instead of $root
    eventBus.on("show-notification", this.showNotification);
  },
  beforeUnmount() {
    // Clean up event listener
    eventBus.off("show-notification", this.showNotification);
  },
  methods: {
    showNotification(notification) {
      // Assign unique ID
      const id = this.counter++;
      const completeNotification = {
        id,
        type: notification.type || "info",
        title: notification.title || "",
        message: notification.message || "",
        duration: notification.duration || 5000,
      };

      // Add to notifications array
      this.notifications.push(completeNotification);

      // Auto-remove after duration
      setTimeout(() => {
        this.removeNotification(id);
      }, completeNotification.duration);
    },
    removeNotification(id) {
      const index = this.notifications.findIndex((n) => n.id === id);
      if (index !== -1) {
        this.notifications.splice(index, 1);
      }
    },
    getIcon(type) {
      switch (type) {
        case "success":
          return "check-circle";
        case "error":
          return "exclamation-circle";
        case "warning":
          return "exclamation-triangle";
        case "message":
          return "comment-alt";
        default:
          return "info-circle";
      }
    },
  },
};
</script>

<style scoped>
.notification-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: var(--z-toast);
  max-width: 350px;
  width: 100%;
}

.notification-toast {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  margin-bottom: var(--space-3);
  padding: var(--space-3) var(--space-4);
  display: flex;
  align-items: flex-start;
  overflow: hidden;
  border-left: 4px solid var(--primary-color);
}

.notification-toast.success {
  border-left-color: var(--success-color);
}

.notification-toast.error {
  border-left-color: var(--error-color);
}

.notification-toast.warning {
  border-left-color: var(--warning-color);
}

.notification-toast.message {
  border-left-color: var(--primary-color);
}

.notification-icon {
  font-size: var(--font-size-lg);
  margin-right: var(--space-3);
  padding-top: var(--space-1);
}

.notification-toast.success .notification-icon {
  color: var(--success-color);
}

.notification-toast.error .notification-icon {
  color: var(--error-color);
}

.notification-toast.warning .notification-icon {
  color: var(--warning-color);
}

.notification-toast.message .notification-icon {
  color: var(--primary-color);
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-1);
}

.notification-message {
  font-size: var(--font-size-sm);
  color: var(--medium);
}

.notification-close {
  background: transparent;
  border: none;
  color: var(--medium);
  padding: var(--space-1);
  cursor: pointer;
  opacity: 0.7;
  margin-left: var(--space-2);
}

.notification-close:hover {
  opacity: 1;
}

/* Animations */
.notification-slide-enter-active,
.notification-slide-leave-active {
  transition: all 0.3s ease;
}

.notification-slide-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.notification-slide-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
