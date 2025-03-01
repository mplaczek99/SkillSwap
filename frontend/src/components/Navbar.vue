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
          <router-link to="/" class="navbar-link" active-class="active">
            <font-awesome-icon icon="home" />
            <span>Home</span>
          </router-link>
          <router-link to="/search" class="navbar-link" active-class="active">
            <font-awesome-icon icon="search" />
            <span>Search</span>
          </router-link>
          <router-link
            v-if="isAuthenticated"
            to="/profile"
            class="navbar-link"
            active-class="active"
          >
            <font-awesome-icon icon="user" />
            <span>Profile</span>
          </router-link>
          <router-link
            v-if="isAuthenticated"
            to="/chat"
            class="navbar-link"
            active-class="active"
          >
            <font-awesome-icon icon="comments" />
            <span>Chat</span>
          </router-link>
          <!-- Add this new link for video uploads -->
          <router-link
            v-if="isAuthenticated"
            to="/upload-video"
            class="navbar-link"
            active-class="active"
          >
            <font-awesome-icon icon="video" />
            <span>Upload Video</span>
          </router-link>
          <router-link
            v-if="isAuthenticated"
            to="/videos"
            class="navbar-link"
            active-class="active"
          >
            <font-awesome-icon icon="film" />
            <span>My Videos</span>
          </router-link>
        </div>

        <div class="navbar-auth">
          <template v-if="!isAuthenticated">
            <router-link to="/login" class="btn btn-outline btn-sm"
              >Login</router-link
            >
            <router-link to="/register" class="btn btn-primary btn-sm"
              >Register</router-link
            >
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
      menuActive: false,
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
    },
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
  gap: var(--space-4);
}

.navbar-link {
  color: var(--dark);
  font-weight: var(--font-weight-medium);
  text-decoration: none;
  padding: var(--space-2) var(--space-1);
  border-radius: var(--radius-md);
  transition:
    color var(--transition-fast) ease,
    background-color var(--transition-fast) ease;
  display: flex;
  align-items: center;
  gap: var(--space-2);
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
    gap: var(--space-2);
  }

  .navbar-link {
    width: 100%;
    padding: var(--space-3) var(--space-2);
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
