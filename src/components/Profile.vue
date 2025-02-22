<template>
  <div class="profile">
    <h2>Your Profile</h2>
    <div v-if="!editing">
      <p><strong>Name:</strong> {{ user.name }}</p>
      <p><strong>Email:</strong> {{ user.email }}</p>
      <p>
        <strong>Bio:</strong>
        <span v-if="user.bio && user.bio.length">{{ user.bio }}</span>
        <span v-else>No bio provided.</span>
      </p>
      <button @click="startEdit">Edit Profile</button>
    </div>
    <div v-else>
      <form @submit.prevent="submitProfile">
        <div>
          <label for="name">Name:</label>
          <input id="name" v-model="editedProfile.name" type="text" required />
        </div>
        <div>
          <label for="email">Email:</label>
          <input id="email" v-model="editedProfile.email" type="email" required />
        </div>
        <div>
          <label for="bio">Bio:</label>
          <textarea id="bio" v-model="editedProfile.bio"></textarea>
        </div>
        <button type="submit">Save</button>
        <button type="button" @click="cancelEdit">Cancel</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Profile',
  data() {
    return {
      editing: false,
      editedProfile: {
        name: '',
        email: '',
        bio: '',
      },
    };
  },
  computed: {
    user() {
      return this.$store.getters.user || {};
    },
  },
  methods: {
    startEdit() {
      this.editing = true;
      this.editedProfile = {
        name: this.user.name,
        email: this.user.email,
        bio: this.user.bio || '',
      };
    },
    submitProfile() {
      this.$store.dispatch('updateProfile', this.editedProfile);
      this.editing = false;
    },
    cancelEdit() {
      this.editing = false;
    },
  },
  watch: {
    user: {
      immediate: true,
      handler(newVal) {
        if (newVal) {
          this.editedProfile = {
            name: newVal.name,
            email: newVal.email,
            bio: newVal.bio || '',
          };
        }
      },
    },
  },
};
</script>

<style scoped>
.profile {
  padding: 2rem;
}
.profile form {
  display: flex;
  flex-direction: column;
  max-width: 400px;
}
.profile form div {
  margin-bottom: 1rem;
}
.profile form label {
  display: block;
  margin-bottom: 0.5rem;
}
.profile form input,
.profile form textarea {
  padding: 0.5rem;
  font-size: 1rem;
  width: 100%;
}
.profile form button {
  margin-right: 0.5rem;
}
</style>

