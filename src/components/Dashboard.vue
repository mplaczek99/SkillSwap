<template>
  <div class="dashboard">
    <h2>Welcome to SkillSwap!</h2>
    <p>Your SkillPoints: {{ user.skillPoints || 0 }}</p>
    <!-- Additional dashboard content can go here -->
  </div>
</template>

<script>
export default {
  name: 'Dashboard',
  computed: {
    user() {
      return this.$store.state.user || {};
    },
  },
};
</script>

<style scoped>
.dashboard {
  padding: 2rem;
}
</style>

