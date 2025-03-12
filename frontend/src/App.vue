<template>
  <div id="app">
    <Navbar />
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <footer class="footer">
      <div class="container">
        <p>&copy; {{ currentYear }} SkillSwap. All rights reserved.</p>
      </div>
    </footer>

    <!-- Single Notification Component -->
    <NotificationComponent />
  </div>
</template>

<script>
import Navbar from "./components/Navbar.vue";
import NotificationComponent from "./components/NotificationComponent.vue";

export default {
  name: "App",
  components: {
    Navbar,
    NotificationComponent,
  },
  computed: {
    currentYear() {
      return new Date().getFullYear();
    },
  },
};
</script>

<style>
/* Import design system */
@import "./assets/css/design-system.css";

/* App specific styles */
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  padding: var(--space-6) 0;
}

.footer {
  padding: var(--space-6) 0;
  background-color: var(--dark);
  color: var(--white);
  text-align: center;
}

/* Page transition animations */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-normal) ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
