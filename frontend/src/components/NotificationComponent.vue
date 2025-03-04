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
      notificationTimeouts: new Map(),
    };
  },
  created() {
    // Listen for notification events using eventBus instead of $root
    eventBus.on("show-notification", this.showNotification);
  },
  beforeUnmount() {
    // Clean up event listener
    eventBus.off("show-notification", this.showNotification);

    // Clear any pending timeouts
    this.notificationTimeouts.forEach((timeoutId) => {
      clearTimeout(timeoutId);
    });
    this.notificationTimeouts.clear();
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
      const timeoutId = setTimeout(() => {
        this.removeNotification(id);
      }, completeNotification.duration);

      // Store the timeout ID for cleanup
      this.notificationTimeouts.set(id, timeoutId);
    },
    removeNotification(id) {
      const index = this.notifications.findIndex((n) => n.id === id);
      if (index !== -1) {
        this.notifications.splice(index, 1);
      }

      // Clear and remove the timeout
      if (this.notificationTimeouts.has(id)) {
        clearTimeout(this.notificationTimeouts.get(id));
        this.notificationTimeouts.delete(id);
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
/* Style remains unchanged */
</style>
