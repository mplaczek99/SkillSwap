<template>
  <div class="profile">
    <h2>Your Profile</h2>
    <!-- Profile Overview Card -->
    <div class="profile-card">
      <div class="profile-avatar-container">
        <template v-if="user && user.avatar">
          <img
            :src="user.avatar"
            alt="Profile Picture"
            class="profile-avatar"
          />
        </template>
        <template v-else>
          <font-awesome-icon
            icon="user"
            class="profile-avatar default-avatar"
          />
        </template>
      </div>
      <div class="profile-info">
        <h3>{{ user ? user.name : "Guest" }}</h3>
        <p>{{ user ? user.email : "" }}</p>
        <p v-if="user && user.bio">{{ user.bio }}</p>
        <p v-else>No bio provided.</p>
        <p>
          <strong>SkillPoints:</strong> {{ user ? user.skillPoints || 0 : 0 }}
        </p>
      </div>
    </div>

    <!-- Toggle Edit Mode -->
    <button @click="toggleEdit" class="edit-button">
      {{ editing ? "Cancel Edit" : "Edit Profile" }}
    </button>

    <!-- Edit Form -->
    <div v-if="editing" class="edit-profile-form">
      <form @submit.prevent="submitProfile">
        <div>
          <label for="name">Name:</label>
          <input id="name" v-model="editedProfile.name" type="text" required />
        </div>
        <div>
          <label for="email">Email:</label>
          <input
            id="email"
            v-model="editedProfile.email"
            type="email"
            required
          />
        </div>
        <div>
          <label for="bio">Bio:</label>
          <textarea id="bio" v-model="editedProfile.bio"></textarea>
        </div>
        <button type="submit">Save Changes</button>
      </form>
    </div>

    <!-- My Skills Section -->
    <div class="my-skills">
      <h3>My Skills</h3>
      <div v-if="userSkills.length">
        <ProfileCard
          v-for="(skill, index) in userSkills"
          :key="index"
          :title="skill.name"
          :description="skill.description"
          :imageSrc="skill.image || defaultSkillImage"
          @open-profile="viewSkill(skill)"
        />
      </div>
      <p v-else>You haven't added any skills yet.</p>
      <button @click="addSkill">Add New Skill</button>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import ProfileCard from "./ProfileCard.vue";

export default {
  name: "Profile",
  components: {
    ProfileCard,
  },
  data() {
    return {
      editing: false,
      editedProfile: {
        name: "",
        email: "",
        bio: "",
      },
      userSkills: [],
      defaultSkillImage: "https://via.placeholder.com/80",
    };
  },
  computed: {
    ...mapGetters(["user"]),
  },
  created() {
    if (this.user) {
      this.editedProfile = {
        name: this.user.name,
        email: this.user.email,
        bio: this.user.bio || "",
      };
    }
    this.fetchUserSkills();
  },
  methods: {
    toggleEdit() {
      this.editing = !this.editing;
      if (!this.editing && this.user) {
        // Reset changes if canceled
        this.editedProfile = {
          name: this.user.name,
          email: this.user.email,
          bio: this.user.bio || "",
        };
      }
    },
    submitProfile() {
      this.$store.dispatch("updateProfile", this.editedProfile);
      this.editing = false;
    },
    fetchUserSkills() {
      // Simulate user skills with dummy data
      const dummySkills = [
        {
          name: "Go Programming",
          description: "Learn the basics of Go",
          image: "",
        },
        {
          name: "Vue.js",
          description: "Frontend development with Vue",
          image: "",
        },
      ];
      this.userSkills = dummySkills;
    },
    addSkill() {
      this.$router.push("/add-skill");
    },
    viewSkill(skill) {
      this.$router.push({
        name: "SkillDetails",
        params: { skillName: skill.name },
      });
    },
  },
  watch: {
    user(newVal) {
      if (newVal) {
        this.editedProfile = {
          name: newVal.name,
          email: newVal.email,
          bio: newVal.bio || "",
        };
      }
    },
  },
};
</script>

<style scoped>
.profile {
  padding: 2rem;
}
.profile-card {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
  background-color: #f7f7f7;
  padding: 1rem;
  border-radius: 8px;
}
.profile-avatar-container {
  width: 80px;
  height: 80px;
  margin-right: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}
.profile-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}
.profile-info h3 {
  margin: 0;
  font-size: 1.5rem;
}
.edit-button {
  margin-bottom: 1rem;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
}
.edit-profile-form {
  margin-bottom: 2rem;
  border: 1px solid #ccc;
  padding: 1rem;
  border-radius: 4px;
}
.edit-profile-form form > div {
  margin-bottom: 1rem;
}
.my-skills {
  margin-top: 2rem;
}
.my-skills h3 {
  margin-bottom: 1rem;
}
.my-skills button {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
}
/* Default Font Awesome icon styling for avatar */
.default-avatar {
  font-size: 3rem;
  color: #ccc;
}
</style>
