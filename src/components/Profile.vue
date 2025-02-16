<template>
  <div class="profile">
    <h2>Your Profile</h2>
    <p>Name: {{ user.name }}</p>
    <p>Email: {{ user.email }}</p>
    <!-- More profile fields and an edit form can be added here -->
  </div>
</template>

<script>
export default {
  name: 'Profile',
  computed: {
    user() {
      return this.$store.state.user || {};
    },
  },
};
</script>

<style scoped>
.profile {
  padding: 2rem;
}
</style>

