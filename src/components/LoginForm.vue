<template>
  <div class="auth-container">
    <h2>Login</h2>
    <form @submit.prevent="submitLogin">
      <input type="email" v-model="email" placeholder="Email" required />
      <input type="password" v-model="password" placeholder="Password" required />
      <button type="submit">Login</button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script>
export default {
  name: 'LoginForm',
  data() {
    return {
      email: '',
      password: '',
      error: null,
    };
  },
  methods: {
    async submitLogin() {
      try {
        await this.$store.dispatch('login', {
          email: this.email,
          password: this.password,
        });
        this.$router.push('/');
      } catch (err) {
        this.error = err.response?.data?.message || 'Login failed';
      }
    },
  },
};
</script>

<style scoped>
.auth-container {
  max-width: 400px;
  margin: 2rem auto;
}
.error {
  color: red;
  margin-top: 1rem;
}
</style>

