<template>
  <div class="profile-page">
    <!-- Header Section -->
    <header class="profile-header">
      <div class="container">
        <h1>Welcome, {{ user ? user.name : "Guest" }}</h1>
        <p>Earn SkillPoints by sharing your skills and learning new ones!</p>
      </div>
    </header>

    <!-- Main Content -->
    <div class="container profile-content">
      <!-- Profile Card -->
      <div class="card profile-card">
        <div class="profile-avatar">
          <template v-if="user && user.avatar">
            <img :src="user.avatar" alt="Profile Picture" />
          </template>
          <template v-else>
            <font-awesome-icon icon="user" class="default-avatar" />
          </template>
        </div>
        <div class="profile-details">
          <h2>{{ user ? user.name : "Guest" }}</h2>
          <p class="email">{{ user ? user.email : "" }}</p>
          <p class="bio" v-if="user && user.bio">{{ user.bio }}</p>
          <p class="bio" v-else>No bio provided.</p>
          <p class="skillpoints">
            SkillPoints: <strong>{{ user ? user.skillPoints || 0 : 0 }}</strong>
          </p>
          <button class="edit-btn" data-test="edit-button" @click="toggleEdit">
            {{ editing ? "Cancel Edit" : "Edit Profile" }}
          </button>
        </div>
      </div>

      <!-- Edit Profile Form (only visible when editing) -->
      <div v-if="editing" class="card edit-profile-card">
        <h3>Edit Profile</h3>
        <form data-test="edit-profile-form" @submit.prevent="submitProfile">
          <div class="form-group">
            <label for="name">Name:</label>
            <input
              id="name"
              v-model="editedProfile.name"
              type="text"
              required
            />
          </div>
          <div class="form-group">
            <label for="email">Email:</label>
            <input
              id="email"
              v-model="editedProfile.email"
              type="email"
              required
            />
          </div>
          <div class="form-group">
            <label for="bio">Bio:</label>
            <textarea id="bio" v-model="editedProfile.bio"></textarea>
          </div>
          <button type="submit" class="save-btn">Save Changes</button>
        </form>
      </div>

      <!-- My Skills Section -->
      <div class="card skills-card">
        <div class="my-skills">
          <h3>My Skills</h3>
          <div class="skills-list">
            <div
              v-for="(skill, index) in userSkills"
              :key="index"
              class="skill-item"
            >
              <SkillImage :skill="skill" />
              <div class="skill-info">
                <h4>{{ skill.name }}</h4>
                <p>{{ skill.description }}</p>
                <button @click="viewSkill(skill)">Learn More</button>
              </div>
            </div>
          </div>
          <button class="add-skill-btn" @click="addSkill">Add New Skill</button>
        </div>
      </div>

      <!-- My Schedules Section -->
      <div class="card schedules-card">
        <h3>My Schedules</h3>
        <div class="schedule-form-wrapper">
          <form @submit.prevent="createSchedule" class="schedule-form">
            <div class="form-group">
              <label>Start Time:</label>
              <input
                type="datetime-local"
                v-model="newSchedule.startTime"
                required
              />
            </div>
            <div class="form-group">
              <label>End Time:</label>
              <input
                type="datetime-local"
                v-model="newSchedule.endTime"
                required
              />
            </div>
            <button type="submit" class="schedule-btn">Schedule Session</button>
          </form>
          <div v-if="scheduleError" class="error-message">
            {{ scheduleError }}
          </div>
          <div v-if="scheduleLoading" class="loading-message">
            Processing...
          </div>
          <ul class="schedule-list" v-if="schedules.length">
            <li v-for="(schedule, index) in schedules" :key="index">
              Session on Skill ID: {{ schedule.skill_id }} from
              {{ formatDate(schedule.startTime) }} to
              {{ formatDate(schedule.endTime) }}
            </li>
          </ul>
          <div
            v-else-if="!scheduleLoading && schedulesFetched"
            class="no-schedules"
          >
            <p>No scheduled sessions found.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import SkillImage from "@/components/SkillImage.vue";
import axios from "axios";

export default {
  name: "Profile",
  components: { SkillImage },
  data() {
    return {
      editing: false,
      editedProfile: { name: "", email: "", bio: "" },
      userSkills: [],
      // Scheduling-related state
      newSchedule: {
        skill_id: 1, // Default for demonstration; can be dynamic.
        startTime: "",
        endTime: "",
      },
      schedules: [],
      scheduleError: null,
      scheduleLoading: false,
      schedulesFetched: false,
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
    // Dummy skills for demonstration
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
    this.fetchSchedules();
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
    async createSchedule() {
      this.scheduleError = null;
      this.scheduleLoading = true;
      try {
        const response = await axios.post("/api/schedule", this.newSchedule);
        this.schedules.push(response.data);
      } catch (error) {
        if (process.env.NODE_ENV !== "test") {
          console.error("Error creating schedule:", error);
        }
        this.scheduleError =
          "Failed to create schedule. Ensure the session is in the future and try again.";
      } finally {
        this.scheduleLoading = false;
      }
    },
    async fetchSchedules() {
      this.scheduleLoading = true;
      this.scheduleError = null;
      try {
        const response = await axios.get("/api/schedule");
        this.schedules = response.data;
      } catch (error) {
        if (process.env.NODE_ENV !== "test") {
          console.error("Error fetching schedules:", error);
        }
        this.scheduleError = "Unable to load schedules.";
      } finally {
        this.scheduleLoading = false;
        this.schedulesFetched = true;
      }
    },
    formatDate(dateStr) {
      const options = {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      };
      return new Date(dateStr).toLocaleDateString(undefined, options);
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
/* General page styling */
.profile-page {
  font-family: "Helvetica Neue", Arial, sans-serif;
  color: #333;
  background: #f4f7f9;
  min-height: 100vh;
}

/* Header section */
.profile-header {
  background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
  color: #fff;
  padding: 2rem 0;
  text-align: center;
}
.profile-header h1 {
  margin: 0;
  font-size: 2.5rem;
}
.profile-header p {
  margin: 0.5rem 0 0;
  font-size: 1.2rem;
}

/* Container */
.container {
  width: 90%;
  max-width: 1200px;
  margin: 0 auto;
}

/* Card styling */
.card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 2rem;
}

/* Profile card */
.profile-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}
.profile-avatar {
  flex-shrink: 0;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
}
.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.default-avatar {
  font-size: 4rem;
  color: #ccc;
}
.profile-details h2 {
  margin: 0;
  font-size: 2rem;
}
.profile-details .email {
  color: #777;
  margin: 0.25rem 0;
}
.profile-details .bio {
  margin: 0.5rem 0;
  font-size: 1rem;
  line-height: 1.4;
}
.profile-details .skillpoints {
  font-size: 1.1rem;
  margin: 0.5rem 0;
}
.edit-btn {
  background: #2575fc;
  color: #fff;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 0.5rem;
}
.edit-btn:hover {
  background: #1b5bb8;
}

/* Edit profile card */
.edit-profile-card h3 {
  margin-top: 0;
}
.form-group {
  margin-bottom: 1rem;
}
.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
}
.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}
.save-btn {
  background: #6a11cb;
  color: #fff;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
}
.save-btn:hover {
  background: #561b9f;
}

/* Skills card */
.skills-card h3 {
  margin-top: 0;
}
.my-skills {
  /* Wrapper for skills section to satisfy test selectors */
}
.skills-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}
.skill-item {
  background: #f9f9f9;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 1rem;
  width: calc(50% - 1rem);
  display: flex;
  align-items: center;
  gap: 1rem;
}
.skill-info h4 {
  margin: 0;
  font-size: 1.25rem;
}
.skill-info p {
  margin: 0.5rem 0;
  font-size: 0.9rem;
  color: #555;
}
.add-skill-btn {
  background: #2575fc;
  color: #fff;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 1rem;
}
.add-skill-btn:hover {
  background: #1b5bb8;
}

/* Schedules card */
.schedules-card h3 {
  margin-top: 0;
}
.schedule-form-wrapper {
  margin-top: 1rem;
}
.schedule-form {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}
.schedule-form .form-group {
  flex: 1;
  min-width: 200px;
}
.schedule-btn {
  background: #6a11cb;
  color: #fff;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  align-self: flex-end;
}
.schedule-btn:hover {
  background: #561b9f;
}
.loading-message {
  font-style: italic;
  margin-top: 1rem;
}
.error-message {
  color: red;
  margin-top: 1rem;
}
.schedule-list {
  list-style: none;
  padding: 0;
  margin-top: 1rem;
}
.schedule-list li {
  padding: 0.75rem;
  border-bottom: 1px solid #eee;
}
.no-schedules {
  margin-top: 1rem;
  color: #777;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .profile-card {
    flex-direction: column;
    text-align: center;
  }
  .skill-item {
    width: 100%;
  }
}
</style>
