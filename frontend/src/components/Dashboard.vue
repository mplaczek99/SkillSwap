<template>
  <div class="dashboard">
    <div class="container">
      <!-- Hero Section -->
      <section class="hero-section">
        <div class="hero-content">
          <h1>Welcome to SkillSwap, {{ user.name || "Guest" }}!</h1>
          <p class="hero-subtitle">
            Your platform for sharing skills and knowledge
          </p>

          <div class="skill-points-card">
            <div class="skill-points-icon">
              <font-awesome-icon icon="star" />
            </div>
            <div class="skill-points-content">
              <p class="skill-points-label">Your SkillPoints</p>
              <p class="skill-points-value">{{ user.skillPoints || 0 }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Main Dashboard Sections -->
      <div class="dashboard-grid">
        <!-- Featured Skills Section -->
        <section class="dashboard-section">
          <div class="section-header">
            <h2>Featured Skills</h2>
            <router-link to="/search" class="btn btn-outline btn-sm">
              View All
              <font-awesome-icon icon="arrow-right" />
            </router-link>
          </div>

          <div class="skills-grid">
            <div
              class="skill-card"
              v-for="(skill, index) in featuredSkills"
              :key="index"
            >
              <div class="skill-card-header">
                <img
                  :src="getSkillImage(skill)"
                  alt="Skill Image"
                  class="skill-image"
                />
              </div>
              <div class="skill-card-body">
                <h3 class="skill-card-title">{{ skill.name }}</h3>
                <p class="skill-card-description">{{ skill.description }}</p>
                <button
                  class="btn btn-primary btn-sm"
                  @click="viewSkillDetails(skill)"
                >
                  Learn More
                </button>
              </div>
            </div>
          </div>
        </section>

        <!-- Recent Activity Section -->
        <section class="dashboard-section">
          <div class="section-header">
            <h2>Recent Activity</h2>
          </div>

          <div class="activity-list">
            <div
              class="activity-item"
              v-for="(activity, index) in recentActivities"
              :key="index"
            >
              <div class="activity-icon">
                <font-awesome-icon icon="history" />
              </div>
              <div class="activity-content">
                <p>{{ activity }}</p>
                <span class="activity-time">{{ randomTimeAgo() }}</span>
              </div>
            </div>
          </div>
        </section>

        <!-- Announcements Section -->
        <section class="dashboard-section">
          <div class="section-header">
            <h2>Announcements</h2>
          </div>

          <div v-if="announcements.length">
            <div
              class="announcement-card"
              v-for="(announcement, index) in announcements"
              :key="index"
            >
              <div class="announcement-header">
                <h3>{{ announcement.title }}</h3>
                <span class="badge badge-primary">New</span>
              </div>
              <p class="announcement-message">{{ announcement.message }}</p>
            </div>
          </div>
          <div v-else class="empty-state">
            <font-awesome-icon icon="bell-slash" class="empty-icon" />
            <p>No announcements at this time.</p>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script>
import { fetchDynamicIcon } from "@/services/iconService";

export default {
  name: "Dashboard",
  data() {
    return {
      featuredSkills: [
        {
          name: "Go Programming",
          description: "Learn the basics of Go programming language",
          image: "",
        },
        {
          name: "Vue.js",
          description: "Frontend development with Vue framework",
          image: "",
        },
        {
          name: "Guitar Lessons",
          description: "Play your favorite tunes on guitar",
          image: "",
        },
      ],
      recentActivities: [
        "Alice taught Spanish to Bob",
        "Charlie added a new skill: Cooking",
        "Test User updated their profile",
        "David completed a Python session with Emma",
      ],
      announcements: [
        {
          title: "New Feature",
          message:
            "We are excited to announce live chat is now available! Connect with other users in real-time.",
        },
        {
          title: "Scheduled Maintenance",
          message:
            "Scheduled maintenance will occur on Saturday at 2 PM. The platform may be unavailable for a short period.",
        },
      ],
    };
  },
  computed: {
    user() {
      return this.$store.state.user || {};
    },
    defaultSkillImage() {
      return "https://api.iconify.design/fa-solid/cog.svg";
    },
  },
  async created() {
    // Pre-fetch dynamic icons for featured skills.
    try {
      await Promise.all(
        this.featuredSkills.map(async (skill) => {
          if (!skill.image || skill.image.trim() === "") {
            skill.dynamicIcon = await fetchDynamicIcon(skill.name);
          }
        }),
      );
    } catch (error) {
      console.error(
        "Failed to fetch dynamic icons for featured skills:",
        error,
      );
    }
  },
  methods: {
    getSkillImage(skill) {
      if (skill.image && skill.image.trim() !== "") {
        return skill.image;
      }
      if (skill.dynamicIcon) {
        // Fixed template string syntax
        return `https://api.iconify.design/fa-solid/${skill.dynamicIcon}.svg`;
      }
      return this.defaultSkillImage;
    },
    viewSkillDetails(skill) {
      // In a real app, you would navigate to a skill details page
      this.$router.push({
        name: "Search",
        query: { q: skill.name },
      });
    },
    randomTimeAgo() {
      const times = ["Just now", "5 minutes ago", "2 hours ago", "Yesterday"];
      return times[Math.floor(Math.random() * times.length)];
    },
  },
};
</script>

<style scoped>
.dashboard {
  padding-bottom: var(--space-12);
}

/* Hero Section */
.hero-section {
  background: linear-gradient(
    135deg,
    var(--primary-color) 0%,
    var(--secondary-color) 100%
  );
  border-radius: var(--radius-lg);
  padding: var(--space-8);
  color: white;
  margin-bottom: var(--space-8);
  box-shadow: var(--shadow-lg);
}

.hero-content {
  max-width: 700px;
}

.hero-section h1 {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
  color: white;
}

.hero-subtitle {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-6);
  opacity: 0.9;
}

.skill-points-card {
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  display: flex;
  align-items: center;
  backdrop-filter: blur(10px);
  max-width: 250px;
}

.skill-points-icon {
  background-color: var(--warning-color);
  color: white;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xl);
  margin-right: var(--space-4);
}

.skill-points-label {
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-1);
  opacity: 0.8;
}

.skill-points-value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  margin: 0;
}

/* Dashboard Grid Layout */
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: var(--space-6);
}

.dashboard-section {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  box-shadow: var(--shadow-md);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-6);
}

.section-header h2 {
  margin-bottom: 0;
  font-size: var(--font-size-xl);
}

/* Skills Grid */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: var(--space-4);
}

.skill-card {
  background-color: var(--light);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition:
    transform var(--transition-normal) ease,
    box-shadow var(--transition-normal) ease;
}

.skill-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.skill-card-header {
  height: 120px;
  background-color: var(--primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
}

.skill-card-header img {
  width: 64px;
  height: 64px;
}

.skill-card-body {
  padding: var(--space-4);
}

.skill-card-title {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-2);
}

.skill-card-description {
  color: var(--medium);
  margin-bottom: var(--space-4);
  font-size: var(--font-size-sm);
  height: 40px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* Activity List */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.activity-item {
  display: flex;
  align-items: center;
  padding: var(--space-3);
  background-color: var(--light);
  border-radius: var(--radius-md);
  transition: background-color var(--transition-fast) ease;
}

.activity-item:hover {
  background-color: var(--primary-light);
}

.activity-icon {
  background-color: var(--secondary-light);
  color: var(--secondary-color);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: var(--space-3);
}

.activity-content {
  flex: 1;
}

.activity-content p {
  margin-bottom: 0;
  font-size: var(--font-size-sm);
}

.activity-time {
  font-size: var(--font-size-xs);
  color: var(--medium);
}

/* Announcements */
.announcement-card {
  background-color: var(--light);
  border-radius: var(--radius-md);
  padding: var(--space-4);
  margin-bottom: var(--space-4);
  border-left: 4px solid var(--primary-color);
}

.announcement-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-2);
}

.announcement-header h3 {
  font-size: var(--font-size-md);
  margin-bottom: 0;
  font-weight: var(--font-weight-semibold);
}

.announcement-message {
  color: var(--medium);
  font-size: var(--font-size-sm);
  margin-bottom: 0;
}

/* Empty state */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-8) 0;
  color: var(--medium);
}

.empty-icon {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
  opacity: 0.5;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .hero-section {
    padding: var(--space-6);
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .skills-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  }
}

@media (max-width: 576px) {
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-2);
  }

  .skill-points-card {
    max-width: 100%;
  }
}
</style>
