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
      <div class="profile-card">
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
          <div class="profile-actions">
            <button
              class="edit-btn"
              data-test="edit-button"
              @click="toggleEdit"
            >
              {{ editing ? "Cancel Edit" : "Edit Profile" }}
            </button>
            <button v-if="!isOwnProfile" class="message-btn" @click="startChat">
              <font-awesome-icon icon="comment" />
              Message
            </button>
          </div>
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
    isOwnProfile() {
      // Check if this is the current user's profile
      return (
        !this.$route.params.userId ||
        (this.user && this.$route.params.userId == this.user.id)
      );
    },
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
    startChat() {
      // Navigate to chat with this user
      const userId = this.$route.params.userId;
      const userName = this.profileData ? this.profileData.name : "User";

      this.$router.push({
        name: "Chat",
        query: {
          user: userId,
          userName: userName,
        },
      });

      // Show notification
      this.$root.$emit("show-notification", {
        type: "info",
        title: "Starting Chat",
        message: `Starting a conversation with ${userName}`,
        duration: 3000,
      });
    },
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
        name: "Search",
        query: { q: skill.name },
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
  background: linear-gradient(to bottom, #f4f7f9, #e8eef2);
  min-height: 100vh;
  padding-bottom: var(--space-8);
}

/* Header section */
.profile-header {
  background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
  color: #fff;
  padding: 3rem 0;
  text-align: center;
  margin-bottom: var(--space-6);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.profile-header h1 {
  margin: 0;
  font-size: 2.5rem;
  letter-spacing: 0.5px;
}

.profile-header p {
  margin: 0.5rem 0 0;
  font-size: 1.2rem;
  opacity: 0.9;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

/* Profile content container */
.profile-content {
  position: relative;
  margin-top: -2rem;
}

/* Card styling */
.card {
  background: #fff;
  border-radius: var(--radius-lg);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
  padding: var(--space-6);
  margin-bottom: var(--space-6);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.card:hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.card h3 {
  position: relative;
  padding-bottom: var(--space-2);
  margin-bottom: var(--space-4);
  border-bottom: 2px solid var(--primary-light);
}

/* Profile Card Enhancements */
.profile-card {
  display: flex;
  align-items: center;
  gap: var(--space-6);
  padding: var(--space-6);
  background: linear-gradient(135deg, #f0f4ff 0%, #e5eeff 100%);
  border-radius: var(--radius-lg);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.06);
  border-left: 4px solid var(--primary-color);
  margin-bottom: var(--space-6);
  position: relative;
  overflow: hidden;
}

.profile-card::after {
  content: "";
  position: absolute;
  top: 0;
  right: 0;
  width: 150px;
  height: 100%;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.2) 100%
  );
  pointer-events: none;
}

.profile-avatar {
  flex-shrink: 0;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border: 5px solid white;
  position: relative;
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f0f4ff;
  color: var(--primary-color);
  font-size: 5rem;
}

.profile-details h2 {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 var(--space-2) 0;
  color: var(--dark);
  letter-spacing: -0.5px;
}

.profile-details .email {
  color: var(--medium);
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-3);
  font-size: 1rem;
}

.profile-details .email::before {
  content: "✉";
  opacity: 0.7;
}

.profile-details .bio {
  margin: var(--space-3) 0;
  font-size: var(--font-size-md);
  line-height: 1.6;
  color: var(--dark);
  padding-left: var(--space-3);
  border-left: 3px solid var(--primary-color);
  max-width: 550px;
}

.profile-details .skillpoints {
  font-size: var(--font-size-lg);
  color: var(--primary-color);
  font-weight: var(--font-weight-semibold);
  display: inline-block;
  padding: var(--space-2) var(--space-3);
  background-color: rgba(79, 70, 229, 0.1);
  border-radius: var(--radius-full);
  margin-top: var(--space-2);
}

.profile-actions {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-4);
}

.edit-btn {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-full);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.edit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background: var(--primary-dark);
}

.message-btn {
  background-color: var(--secondary-color);
}

.message-btn:hover {
  background-color: var(--secondary-dark);
}

/* Skills Section Improvements */
.card {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  padding: var(--space-6);
  margin-bottom: var(--space-6);
  border: 1px solid #f0f0f0;
  position: relative;
  overflow: hidden;
}

.card h3 {
  color: var(--dark);
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 0;
  margin-bottom: var(--space-4);
  display: flex;
  align-items: center;
  padding-bottom: var(--space-3);
  border-bottom: 2px solid #f0f0f0;
}

.skills-card h3::before {
  content: "★";
  margin-right: var(--space-2);
  color: var(--primary-color);
}

.skills-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: var(--space-4);
  margin-top: var(--space-4);
}

.skill-item {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  display: flex;
  align-items: center;
  gap: var(--space-4);
  transition: all 0.25s ease;
  border: 1px solid #f0f0f0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
  position: relative;
  overflow: hidden;
}

.skill-item::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 5px;
  height: 100%;
  background: var(--primary-color);
  opacity: 0;
  transition: opacity 0.25s ease;
}

.skill-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.08);
  border-color: var(--primary-light);
}

.skill-item:hover::after {
  opacity: 1;
}

.skill-info {
  flex: 1;
}

.skill-info h4 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--dark);
  font-weight: 600;
}

.skill-info p {
  margin: var(--space-1) 0 var(--space-3) 0;
  color: var(--medium);
  font-size: 0.9rem;
  line-height: 1.4;
}

.skill-item button {
  background-color: var(--primary-light);
  color: var(--primary-color);
  border: none;
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-full);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.skill-item button:hover {
  background-color: var(--primary-color);
  color: white;
  transform: translateY(-1px);
}

.add-skill-btn {
  background: var(--secondary-color);
  color: white;
  border: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-full);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: var(--space-4);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  width: auto;
  align-self: flex-start;
}

.add-skill-btn::before {
  content: "+";
  font-weight: bold;
  font-size: 1.2rem;
}

.add-skill-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  background: var(--secondary-dark);
}

/* Schedule Section Improvements */
.schedules-card h3::before {
  content: "📅";
  margin-right: var(--space-2);
  font-size: 1.2rem;
}

.schedule-form {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-4);
  background-color: #f8fafc;
  padding: var(--space-4);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-4);
  border: 1px solid #edf2f7;
}

.schedule-btn {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-full);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.schedule-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background: var(--primary-dark);
}

.error-message {
  background-color: rgba(239, 68, 68, 0.08);
  color: var(--error-color);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  border-left: 4px solid var(--error-color);
  margin: var(--space-3) 0;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
}

.error-message::before {
  content: "⚠️";
  margin-right: var(--space-2);
}

.no-schedules {
  text-align: center;
  padding: var(--space-8);
  background-color: #f8fafc;
  border-radius: var(--radius-lg);
  color: var(--medium);
  border: 1px dashed #e2e8f0;
  font-size: 1rem;
}

.no-schedules::before {
  content: "📆";
  display: block;
  font-size: 2.5rem;
  margin-bottom: var(--space-3);
  opacity: 0.7;
}

/* Responsive improvements */
@media (max-width: 768px) {
  .profile-card {
    flex-direction: column;
    text-align: center;
    padding: var(--space-4);
  }

  .profile-details .bio {
    border-left: none;
    padding-left: 0;
    text-align: center;
  }

  .profile-actions {
    justify-content: center;
  }

  .skills-list {
    grid-template-columns: 1fr;
  }
}
</style>
