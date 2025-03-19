<template>
  <nav class="navbar" :class="{ 'dark-mode': isDark }">
    <div class="container navbar-container">
      <div class="navbar-brand">
        <router-link to="/" class="navbar-logo">
          <font-awesome-icon icon="code" class="logo-icon" />
          <span class="logo-text">SkillSwap</span>
        </router-link>
        <button
          class="navbar-toggler"
          @click="toggleMenu"
          aria-label="Toggle navigation"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>

      <div class="navbar-menu" :class="{ 'is-active': menuActive }">
        <div class="navbar-links">
          <!-- Primary Links -->
          <router-link to="/" class="navbar-link" active-class="active">
            <font-awesome-icon icon="home" />
            <span>Home</span>
          </router-link>
          <router-link to="/search" class="navbar-link" active-class="active">
            <font-awesome-icon icon="search" />
            <span>Search</span>
          </router-link>

          <!-- Jobs Link -->
          <router-link to="/jobs" class="navbar-link" active-class="active">
            <font-awesome-icon icon="briefcase" />
            <span>Jobs</span>
          </router-link>

          <!-- Authenticated-only content -->
          <template v-if="isAuthenticated">
            <!-- Profile Dropdown -->
            <div class="navbar-dropdown-wrapper">
              <div class="navbar-link dropdown-trigger">
                <font-awesome-icon icon="user" />
                <span>Profile</span>
                <font-awesome-icon icon="chevron-down" class="dropdown-icon" />
              </div>
              <div class="dropdown-menu">
                <router-link to="/profile" class="dropdown-item">
                  <font-awesome-icon icon="user" />
                  <span>My Profile</span>
                </router-link>
                <router-link to="/feedback" class="dropdown-item">
                  <font-awesome-icon icon="star" />
                  <span>Feedback</span>
                </router-link>
              </div>
            </div>

            <!-- Learning Dropdown -->
            <div class="navbar-dropdown-wrapper">
              <div class="navbar-link dropdown-trigger">
                <font-awesome-icon icon="graduation-cap" />
                <span>Learning</span>
                <font-awesome-icon icon="chevron-down" class="dropdown-icon" />
              </div>
              <div class="dropdown-menu">
                <router-link to="/videos" class="dropdown-item">
                  <font-awesome-icon icon="film" />
                  <span>My Videos</span>
                </router-link>
                <router-link to="/upload-video" class="dropdown-item">
                  <font-awesome-icon icon="video" />
                  <span>Upload Video</span>
                </router-link>
              </div>
            </div>

            <!-- Exchange Dropdown -->
            <div class="navbar-dropdown-wrapper">
              <div class="navbar-link dropdown-trigger">
                <font-awesome-icon icon="exchange-alt" />
                <span>Exchange</span>
                <font-awesome-icon icon="chevron-down" class="dropdown-icon" />
              </div>
              <div class="dropdown-menu">
                <router-link to="/transactions" class="dropdown-item">
                  <font-awesome-icon icon="coins" />
                  <span>SkillPoints</span>
                </router-link>
                <router-link to="/schedule" class="dropdown-item">
                  <font-awesome-icon icon="calendar-alt" />
                  <span>Schedule</span>
                </router-link>
                <router-link to="/chat" class="dropdown-item position-relative">
                  <font-awesome-icon icon="comments" />
                  <span>Chat</span>
                  <span v-if="unreadMessagesCount > 0" class="unread-badge">{{
                    unreadMessagesCount
                  }}</span>
                </router-link>

                <!-- Post Job Link -->
                <router-link to="/post-job" class="dropdown-item">
                  <font-awesome-icon icon="plus-circle" />
                  <span>Post a Job</span>
                </router-link>
              </div>
            </div>
          </template>
        </div>

        <div class="navbar-auth">
          <template v-if="!isAuthenticated">
            <router-link to="/login" class="btn btn-outline btn-sm">
              Login
            </router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">
              Register
            </router-link>
          </template>
          <template v-else>
            <!-- Chat Button -->
            <router-link
              to="/chat"
              class="btn btn-icon btn-sm navbar-chat-btn"
              v-tooltip="'Messages'"
            >
              <font-awesome-icon icon="comments" />
              <span v-if="unreadMessagesCount > 0" class="unread-badge">{{
                unreadMessagesCount
              }}</span>
            </router-link>
            <button @click="handleLogout" class="btn btn-outline btn-sm logout-btn">
              <font-awesome-icon icon="sign-out-alt" />
              <span>Logout</span>
            </button>
          </template>
          
          <!-- Theme Toggle Button -->
          <button 
            @click="toggleDark()" 
            class="btn btn-icon btn-sm theme-toggle-btn" 
            v-tooltip="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
          >
            <font-awesome-icon :icon="isDark ? 'sun' : 'moon'" />
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { mapGetters, mapMutations } from "vuex";
import ChatService from "@/services/ChatService";
import eventBus from "@/utils/eventBus";
import { useDark, useToggle } from '@vueuse/core';

export default {
  name: "Navbar",
  data() {
    return {
      menuActive: false,
      unreadMessagesCount: 0,
    };
  },
  computed: {
    ...mapGetters(["isAuthenticated"]),
  },
  setup() {
    const isDark = useDark();
    const toggleDark = useToggle(isDark);
    
    return {
      isDark,
      toggleDark,
    };
  },
  mounted() {
    // Listen for message events using the event bus.
    eventBus.on("chat:new-message", this.handleNewMessage);
    eventBus.on("chat:read-messages", this.handleReadMessages);
    this.fetchUnreadCount();
    
    // Close menu when clicking outside
    document.addEventListener('click', this.closeMenuOnClickOutside);
  },
  beforeUnmount() {
    eventBus.off("chat:new-message", this.handleNewMessage);
    eventBus.off("chat:read-messages", this.handleReadMessages);
    document.removeEventListener('click', this.closeMenuOnClickOutside);
  },
  methods: {
    ...mapMutations(["logout"]),
    toggleMenu() {
      this.menuActive = !this.menuActive;
    },
    closeMenuOnClickOutside(event) {
      const navbar = document.querySelector('.navbar');
      if (this.menuActive && navbar && !navbar.contains(event.target)) {
        this.menuActive = false;
      }
    },
    handleLogout() {
      this.logout();
      this.$router.push("/login");
      this.menuActive = false;
    },
    async fetchUnreadCount() {
      if (!this.isAuthenticated) return;
      try {
        this.unreadMessagesCount = await ChatService.getUnreadCount();
      } catch (error) {
        console.error("Error fetching unread count:", error);
      }
    },
    handleNewMessage() {
      // Increment the unread count.
      this.unreadMessagesCount += 1;
    },
    handleReadMessages() {
      // Refresh the unread count when messages are read.
      this.fetchUnreadCount();
    },
  },
  directives: {
    tooltip: {
      mounted(el, binding) {
        el.setAttribute("title", binding.value);
      },
      updated(el, binding) {
        el.setAttribute("title", binding.value);
      },
    },
  },
};
</script>

<style scoped>
:root {
  --primary-color: #4f46e5;
  --primary-hover: #4338ca;
  --primary-light: #eef2ff;
  --dark: #1f2937;
  --gray-light: #f3f4f6;
  --gray-medium: #9ca3af;
  --white: #ffffff;
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --radius-sm: 0.125rem;
  --radius-md: 0.375rem;
  --radius-lg: 0.5rem;
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;
  --font-weight-normal: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-5: 1.25rem;
  --space-6: 1.5rem;
  --transition-fast: 150ms;
  --transition-normal: 300ms;
  --z-header: 50;
}

/* Navbar Base */
.navbar {
  background-color: var(--white);
  box-shadow: var(--shadow-md);
  position: sticky;
  top: 0;
  z-index: var(--z-header);
  border-bottom: 1px solid rgba(229, 231, 235, 0.8);
  height: 4rem;
}

.navbar-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1rem;
}

/* Brand & Logo */
.navbar-brand {
  display: flex;
  align-items: center;
}

.navbar-logo {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--primary-color);
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: var(--space-2);
  transition: transform var(--transition-fast) ease;
}

.navbar-logo:hover {
  text-decoration: none;
  transform: scale(1.05);
}

.logo-icon {
  color: var(--primary-color);
  font-size: 1.5rem;
}

.logo-text {
  background: linear-gradient(135deg, var(--primary-color), #6366f1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 800;
}

/* Menu & Links */
.navbar-menu {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-grow: 1;
  margin-left: var(--space-6);
}

.navbar-links {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.navbar-link {
  color: var(--dark);
  font-weight: var(--font-weight-medium);
  text-decoration: none;
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast) ease;
  display: flex;
  align-items: center;
  gap: var(--space-2);
  position: relative;
  font-size: var(--font-size-sm);
}

.navbar-link:hover {
  color: var(--primary-color);
  text-decoration: none;
  background-color: var(--primary-light);
}

.navbar-link.active {
  color: var(--primary-color);
  font-weight: var(--font-weight-semibold);
  background-color: var(--primary-light);
}

.navbar-link.active::after {
  content: '';
  position: absolute;
  bottom: -0.5rem;
  left: 50%;
  transform: translateX(-50%);
  width: 1.5rem;
  height: 2px;
  background-color: var(--primary-color);
  border-radius: 1px;
}

/* Dropdown styles */
.navbar-dropdown-wrapper {
  position: relative;
}

.dropdown-trigger {
  cursor: pointer;
  user-select: none;
}

.dropdown-icon {
  font-size: 0.7em;
  margin-left: var(--space-1);
  transition: transform var(--transition-fast) ease;
}

.navbar-dropdown-wrapper:hover .dropdown-icon {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 220px;
  background-color: var(--white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all var(--transition-normal) ease;
  pointer-events: none;
  z-index: 100;
  overflow: hidden;
  border: 1px solid rgba(229, 231, 235, 0.8);
  /* Add padding to create space between trigger and dropdown */
  margin-top: 0.5rem;
}

/* Fix for dropdown disappearing - add padding to create a hover bridge */
.navbar-dropdown-wrapper::after {
  content: '';
  position: absolute;
  height: 0.5rem;
  left: 0;
  right: 0;
  bottom: -0.5rem;
  z-index: 99;
}

.navbar-dropdown-wrapper:hover .dropdown-menu {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
  pointer-events: auto;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  color: var(--dark);
  text-decoration: none;
  transition: all var(--transition-fast) ease;
  position: relative;
  font-size: var(--font-size-sm);
}

.dropdown-item:hover {
  background-color: var(--primary-light);
  color: var(--primary-color);
  text-decoration: none;
}

/* Auth Section */
.navbar-auth {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast) ease;
  cursor: pointer;
  text-decoration: none;
  padding: 0.5rem 1rem;
  font-size: var(--font-size-sm);
}

.btn-sm {
  padding: 0.375rem 0.75rem;
  font-size: var(--font-size-xs);
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
  border: 1px solid var(--primary-color);
}

.btn-primary:hover {
  background-color: var(--primary-hover);
  border-color: var(--primary-hover);
}

.btn-outline {
  background-color: transparent;
  color: var(--dark);
  border: 1px solid var(--gray-medium);
}

.btn-outline:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.logout-btn {
  color: #ef4444;
  border-color: #ef4444;
}

.logout-btn:hover {
  background-color: #fef2f2;
  color: #dc2626;
  border-color: #dc2626;
}

/* Chat Button */
.navbar-chat-btn {
  position: relative;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: var(--primary-light);
  color: var(--primary-color);
  transition: all var(--transition-fast) ease;
  border: none;
}

.navbar-chat-btn:hover {
  background-color: var(--primary-color);
  color: white;
  transform: scale(1.05);
}

/* Theme Toggle Button */
.theme-toggle-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: var(--gray-light);
  color: var(--dark);
  transition: all var(--transition-fast) ease;
  border: none;
}

.theme-toggle-btn:hover {
  background-color: var(--gray-medium);
  color: var(--white);
  transform: scale(1.05);
}

.btn-icon {
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Badges */
.position-relative {
  position: relative;
}

.unread-badge {
  position: absolute;
  top: 0;
  right: 0;
  background-color: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: var(--font-weight-bold);
  border-radius: 50%;
  min-width: 18px;
  height: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
  transform: translate(50%, -50%);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.dropdown-item .unread-badge {
  right: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
}

/* Mobile Toggle */
.navbar-toggler {
  display: none;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: var(--space-2);
  flex-direction: column;
  gap: 5px;
}

.navbar-toggler span {
  display: block;
  width: 25px;
  height: 2px;
  background-color: var(--dark);
  transition: transform var(--transition-fast) ease, opacity var(--transition-fast) ease;
}

.navbar-toggler:hover span {
  background-color: var(--primary-color);
}

.is-active .navbar-toggler span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}

.is-active .navbar-toggler span:nth-child(2) {
  opacity: 0;
}

.is-active .navbar-toggler span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

/* Dark Mode Styles */
.dark-mode {
  background-color: #1e293b;
  border-bottom-color: #334155;
}

.dark-mode .navbar-link {
  color: #e2e8f0;
}

.dark-mode .navbar-link:hover {
  background-color: rgba(79, 70, 229, 0.2);
}

.dark-mode .navbar-link.active {
  background-color: rgba(79, 70, 229, 0.2);
}

.dark-mode .dropdown-menu {
  background-color: #1e293b;
  border-color: #334155;
}

.dark-mode .dropdown-item {
  color: #e2e8f0;
}

.dark-mode .dropdown-item:hover {
  background-color: rgba(79, 70, 229, 0.2);
}

.dark-mode .btn-outline {
  color: #e2e8f0;
  border-color: #64748b;
}

.dark-mode .theme-toggle-btn {
  background-color: #334155;
  color: #f8fafc;
}

.dark-mode .navbar-toggler span {
  background-color: #e2e8f0;
}

/* Mobile navigation */
@media (max-width: 768px) {
  .navbar {
    height: auto;
  }
  
  .navbar-container {
    flex-wrap: wrap;
    padding: 0.75rem 1rem;
  }
  
  .navbar-toggler {
    display: flex;
    margin-left: auto;
  }
  
  .navbar-menu {
    display: none;
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
    padding: var(--space-4) 0;
    margin-left: 0;
    background-color: var(--white);
    border-top: 1px solid rgba(229, 231, 235, 0.8);
    margin-top: 0.75rem;
  }
  
  .dark-mode .navbar-menu {
    background-color: #1e293b;
    border-top-color: #334155;
  }
  
  .navbar-menu.is-active {
    display: flex;
    animation: slideDown 0.3s ease forwards;
  }
  
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .navbar-links {
    flex-direction: column;
    align-items: flex-start;
    width: 100%;
    gap: var(--space-1);
  }
  
  .navbar-link,
  .dropdown-trigger {
    width: 100%;
    padding: var(--space-3) var(--space-2);
  }
  
  .navbar-link.active::after {
    display: none;
  }
  
  .navbar-dropdown-wrapper {
    width: 100%;
  }
  
  .navbar-dropdown-wrapper::after {
    display: none;
  }
  
  .dropdown-menu {
    position: static;
    box-shadow: none;
    opacity: 1;
    visibility: visible;
    transform: none;
    pointer-events: auto;
    max-height: 0;
    transition: max-height var(--transition-normal) ease;
    padding: 0;
    margin-left: var(--space-6);
    overflow: hidden;
    border: none;
    border-left: 2px solid var(--primary-light);
    margin-top: 0;
  }
  
  .dark-mode .dropdown-menu {
    border-left-color: rgba(79, 70, 229, 0.4);
  }
  
  .navbar-dropdown-wrapper:hover .dropdown-menu {
    max-height: 500px;
  }
  
  .dropdown-item {
    padding: var(--space-2) var(--space-4);
  }
  
  .navbar-auth {
    margin-top: var(--space-4);
    width: 100%;
    justify-content: center;
    gap: var(--space-4);
  }
  
  .navbar-auth .btn {
    width: 100%;
    justify-content: center;
  }
  
  .dropdown-item .unread-badge {
    right: auto;
    left: 85px;
  }
}

/* Medium screens */
@media (min-width: 769px) and (max-width: 1024px) {
  .navbar-link {
    padding: var(--space-2) var(--space-2);
  }
  
  .navbar-menu {
    margin-left: var(--space-4);
  }
}
</style>