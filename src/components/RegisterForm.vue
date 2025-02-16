<template>
  <div class="auth-container">
    <h2>Register</h2>
    <form @submit.prevent="submitRegister">
      <input type="text" v-model="name" placeholder="Name" required />
      <input type="email" v-model="email" placeholder="Email" required />
      <input type="password" v-model="password" placeholder="Password" required />
      <button type="submit">Register</button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script>
export default {
  name: 'RegisterForm',
  data() {
    return {
      name: '',
      email: '',
      password: '',
      error: null,
    };
  },
  methods: {
    async submitRegister() {
      try {
        await this.$store.dispatch('register', {
          name: this.name,
          email: this.email,
          password: this.password,
        });
        this.$router.push('/');
      } catch (err) {
        this.error = err.response?.data?.message || 'Registration failed';
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

