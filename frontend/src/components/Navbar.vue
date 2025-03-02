<template>
  <nav class="navbar">
    <div class="container navbar-container">
      <div class="navbar-brand">
        <router-link to="/" class="navbar-logo">
          <font-awesome-icon icon="code" class="logo-icon" />
          SkillSwap
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
          <!-- Primary Links - Always visible -->
          <router-link to="/" class="navbar-link" active-class="active">
            <font-awesome-icon icon="home" />
            <span>Home</span>
          </router-link>
          <router-link to="/search" class="navbar-link" active-class="active">
            <font-awesome-icon icon="search" />
            <span>Search</span>
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
                <router-link to="/chat" class="dropdown-item">
                  <font-awesome-icon icon="comments" />
                  <span>Chat</span>
                </router-link>
              </div>
            </div>
          </template>
        </div>

        <div class="navbar-auth">
          <template v-if="!isAuthenticated">
            <router-link to="/login" class="btn btn-outline btn-sm">Login</router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">Register</router-link>
          </template>
          <button v-else @click="handleLogout" class="btn btn-outline btn-sm">
            <font-awesome-icon icon="sign-out-alt" />
            Logout
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { mapGetters, mapMutations } from "vuex";

export default {
  name: "Navbar",
  data() {
    return {
      menuActive: false
    };
  },
  computed: {
    ...mapGetters(["isAuthenticated"]),
  },
  methods: {
    ...mapMutations(["logout"]),
    toggleMenu() {
      this.menuActive = !this.menuActive;
    },
    handleLogout() {
      this.logout();
      this.$router.push("/login");
      this.menuActive = false;
    }
  },
  watch: {
    $route() {
      // Close mobile menu when route changes
      this.menuActive = false;
    },
  },
};
</script>

<style scoped>
.navbar {
  background-color: var(--white);
  box-shadow: var(--shadow-md);
  position: sticky;
  top: 0;
  z-index: var(--z-header);
}

.navbar-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 4rem;
}

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
}

.navbar-logo:hover {
  text-decoration: none;
}

.logo-icon {
  color: var(--primary-color);
}

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
}

.navbar-link:hover {
  color: var(--primary-color);
  text-decoration: none;
  background-color: var(--primary-light);
}

.navbar-link.active {
  color: var(--primary-color);
  font-weight: var(--font-weight-semibold);
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
  min-width: 200px;
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
}

.dropdown-item:hover {
  background-color: var(--primary-light);
  color: var(--primary-color);
  text-decoration: none;
}

.navbar-auth {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

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
  height: 3px;
  background-color: var(--dark);
  transition: transform var(--transition-fast) ease;
}

/* Mobile navigation */
@media (max-width: 768px) {
  .navbar-container {
    flex-wrap: wrap;
  }

  .navbar-toggler {
    display: flex;
  }

  .navbar-menu {
    display: none;
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
    padding: var(--space-4) 0;
    margin-left: 0;
  }

  .navbar-menu.is-active {
    display: flex;
  }

  .navbar-links {
    flex-direction: column;
    align-items: flex-start;
    width: 100%;
    gap: var(--space-1);
  }

  .navbar-link, .dropdown-trigger {
    width: 100%;
    padding: var(--space-3) var(--space-2);
  }

  .navbar-dropdown-wrapper {
    width: 100%;
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
  }

  .navbar-auth .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
