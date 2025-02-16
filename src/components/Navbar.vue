<template>
  <nav class="navbar">
    <router-link to="/">SkillSwap</router-link>
    <router-link to="/search">Search</router-link>
    <router-link v-if="isAuthenticated" to="/profile">Profile</router-link>
    <router-link v-if="isAuthenticated" to="/chat">Chat</router-link>
    <div class="auth-links">
      <router-link v-if="!isAuthenticated" to="/login">Login</router-link>
      <router-link v-if="!isAuthenticated" to="/register">Register</router-link>
      <button v-if="isAuthenticated" @click="handleLogout">Logout</button>
    </div>
  </nav>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex';

export default {
  name: 'Navbar',
  computed: {
    ...mapGetters(['isAuthenticated']),
  },
  methods: {
    ...mapMutations(['logout']),
    handleLogout() {
      this.logout();
      this.$router.push('/login');
    },
  },
};
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  padding: 1rem;
  background: #3498db;
  color: #fff;
}
.auth-links > * {
  margin-left: 1rem;
}
</style>

