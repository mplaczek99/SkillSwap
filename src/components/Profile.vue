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
        <div
          v-for="(skill, index) in userSkills"
          :key="index"
          class="skill-item"
        >
          <SkillImage :skill="skill" />
          <div class="skill-details">
            <h4>{{ skill.name }}</h4>
            <p>{{ skill.description }}</p>
            <button @click="viewSkill(skill)">Learn More</button>
          </div>
        </div>
      </div>
      <p v-else>You haven't added any skills yet.</p>
      <button @click="addSkill">Add New Skill</button>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import SkillImage from "@/components/SkillImage.vue";

export default {
  name: "Profile",
  components: { SkillImage },
  data() {
    return {
      editing: false,
      editedProfile: { name: "", email: "", bio: "" },
      userSkills: [],
    };
  },
  computed: {
    ...mapGetters(["user"]),
  },
  async created() {
    if (this.user) {
      this.editedProfile = {
        name: this.user.name,
        email: this.user.email,
        bio: this.user.bio || "",
      };
    }
    // Simulate fetching user skills via API.
    this.userSkills = [
      {
        name: "Go Programming",
        description: "Learn the basics of Go",
        image: "",
      },
      {
        name: "Vue.js Development",
        description: "Frontend development with Vue",
        image: "",
      },
      {
        name: "Guitar Lessons",
        description: "Play your favorite tunes",
        image: "",
      },
      {
        name: "Creative Cooking",
        description: "Exploring cuisines",
        image: "",
      },
    ];
  },
  methods: {
    toggleEdit() {
      this.editing = !this.editing;
      if (!this.editing && this.user) {
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
.my-skills {
  margin-top: 2rem;
}
.skill-item {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}
.skill-details {
  margin-left: 1rem;
}
</style>
