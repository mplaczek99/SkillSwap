<template>
  <div class="dashboard">
    <div class="container">
      <!-- Hero Section -->
      <section class="hero-section">
        <div class="hero-content">
          <h1>Welcome to <span class="gradient-text">SkillSwap</span>, {{ user.name }}!</h1>
          <p class="hero-subtitle">Discover, learn, and share skills with our global community</p>
          
          <div class="hero-stats">
            <div class="stat-card">
              <div class="stat-icon">
                <font-awesome-icon icon="coins" />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ user.skillPoints || 0 }}</span>
                <span class="stat-label">SkillPoints</span>
              </div>
              <button class="stat-action" @click="openSendPointsModal">
                <font-awesome-icon icon="paper-plane" /> Send
              </button>
            </div>
            
            <div class="stat-card">
              <div class="stat-icon">
                <font-awesome-icon icon="calendar-check" />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ upcomingSessions }}</span>
                <span class="stat-label">Upcoming Sessions</span>
              </div>
              <button class="stat-action" @click="$router.push('/schedule')">
                <font-awesome-icon icon="calendar-plus" /> Schedule
              </button>
            </div>
          </div>
        </div>
        
        <div class="hero-illustration">
          <img src="@/assets/images/skill-sharing.svg" alt="Skill sharing illustration" />
        </div>
      </section>

      <!-- Quick Actions -->
      <section class="quick-actions">
        <button class="action-button" @click="$router.push('/search')">
          <font-awesome-icon icon="search" />
          <span>Find Skills</span>
        </button>
        <button class="action-button" @click="$router.push('/upload-video')">
          <font-awesome-icon icon="video" />
          <span>Share Video</span>
        </button>
        <button class="action-button" @click="$router.push('/jobs')">
          <font-awesome-icon icon="briefcase" />
          <span>Browse Jobs</span>
        </button>
        <button class="action-button" @click="$router.push('/chat')">
          <font-awesome-icon icon="comments" />
          <span>Messages</span>
        </button>
      </section>

      <!-- Main Dashboard Sections -->
      <div class="dashboard-grid">
        <!-- Featured Skills Section -->
        <section class="dashboard-section skills-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="star" /> Featured Skills</h2>
            <button class="view-all-btn" @click="$router.push('/search')">
              View All <font-awesome-icon icon="arrow-right" />
            </button>
          </div>

          <div class="skills-grid">
            <div
              class="skill-card"
              v-for="(skill, index) in featuredSkills"
              :key="index"
              @click="viewSkillDetails(skill)"
            >
              <div class="skill-header" :style="{ backgroundColor: getColorBySkillCategory(skill.category || 'programming') }">
                <div class="skill-icon">
                  <font-awesome-icon :icon="getIconForSkill(skill.name)" />
                </div>
                <div class="skill-level">{{ skill.level || 'Beginner' }}</div>
              </div>
              <div class="skill-content">
                <h3>{{ skill.name }}</h3>
                <p>{{ skill.description }}</p>
                <div class="skill-meta">
                  <span class="skill-duration">
                    <font-awesome-icon icon="clock" /> {{ skill.duration || '4 weeks' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Recent Activity Section -->
        <section class="dashboard-section activity-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="history" /> Recent Activity</h2>
          </div>

          <div class="activity-list">
            <div
              class="activity-item"
              v-for="(activity, index) in recentActivities"
              :key="index"
            >
              <div class="activity-icon">
                <font-awesome-icon icon="circle-dot" />
              </div>
              <div class="activity-content">
                <p>{{ activity }}</p>
                <span class="activity-time">{{ randomTimeAgo() }}</span>
              </div>
            </div>
          </div>
        </section>

        <!-- Announcements Section -->
        <section class="dashboard-section announcements-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="bullhorn" /> Announcements</h2>
          </div>

          <div class="announcements-list">
            <div
              class="announcement-item"
              v-for="(announcement, index) in announcements"
              :key="index"
            >
              <div class="announcement-badge" :class="{ 'new': index === 0 }">
                {{ index === 0 ? 'NEW' : '' }}
              </div>
              <h3>{{ announcement.title }}</h3>
              <p>{{ announcement.message }}</p>
            </div>
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
      upcomingSessions: 2, // This would come from your API in a real app
      featuredSkills: [
        {
          name: "Go Programming",
          description: "Learn the basics of Go programming language",
          category: "programming",
          level: "Beginner",
          duration: "4 weeks"
        },
        {
          name: "Vue.js",
          description: "Frontend development with Vue framework",
          category: "programming",
          level: "Intermediate",
          duration: "6 weeks"
        },
        {
          name: "Guitar Lessons",
          description: "Play your favorite tunes on guitar",
          category: "music",
          level: "Beginner",
          duration: "8 weeks"
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
      return this.$store.state.user || {
        name: "Guest",
        skillPoints: 0
      };
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
    getColorBySkillCategory(category) {
      const categoryColors = {
        programming: "#4361ee",
        design: "#3a86ff",
        language: "#ff9f1c",
        music: "#7209b7",
        cooking: "#e63946",
        business: "#2a9d8f"
      };
      
      return categoryColors[category] || "#4361ee"; // Default to blue if category not found
    },
    getIconForSkill(skillName) {
      const skillIcons = {
        "Go Programming": "code",
        "Vue.js": "code",
        "Guitar Lessons": "guitar",
        "Python": "code",
        "Cooking": "utensils",
        "Spanish": "language"
      };
      
      return skillIcons[skillName] || "graduation-cap"; // Default icon
    },
    openSendPointsModal() {
      // You could navigate to the transactions page
      this.$router.push('/transactions');
      // Or implement a modal directly in this component
    }
  },
};
</script>

<style scoped>
.dashboard {
  padding-bottom: 4rem;
}

/* Hero Section */
.hero-section {
  background: linear-gradient(135deg, #4361ee 0%, #3a0ca3 100%);
  border-radius: 16px;
  padding: 3rem;
  color: white;
  margin-bottom: 2rem;
  box-shadow: 0 10px 25px rgba(67, 97, 238, 0.3);
  display: flex;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  position: relative;
}

.hero-content {
  max-width: 60%;
  position: relative;
  z-index: 2;
}

.hero-illustration {
  max-width: 40%;
  position: relative;
  z-index: 1;
}

.hero-illustration img {
  max-width: 100%;
  height: auto;
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
  100% { transform: translateY(0px); }
}

.hero-section h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  font-weight: 800;
  line-height: 1.2;
  color: white;
}

.gradient-text {
  background: linear-gradient(90deg, #f9c80e, #ff9f1c);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  font-weight: 900;
}

.hero-subtitle {
  font-size: 1.25rem;
  margin-bottom: 2rem;
  opacity: 0.9;
  max-width: 80%;
}

.hero-stats {
  display: flex;
  gap: 1.5rem;
  margin-top: 2rem;
}

.stat-card {
  background-color: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
  transition: all 0.3s ease;
}

.stat-card:hover {
  background-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-5px);
}

.stat-icon {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.stat-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1;
}

.stat-label {
  font-size: 0.875rem;
  opacity: 0.8;
  margin-top: 0.25rem;
}

.stat-action {
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.stat-action:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

/* Quick Actions */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  margin-bottom: 2rem;
}

.action-button {
  background-color: white;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.action-button:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
  border-color: #4361ee;
}

.action-button i {
  font-size: 1.5rem;
  color: #4361ee;
}

.action-button span {
  font-weight: 600;
}

/* Dashboard Grid Layout */
.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

.dashboard-section {
  background-color: white;
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  height: 100%;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f0f0f0;
}

.section-header h2 {
  margin-bottom: 0;
  font-size: 1.25rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.section-header h2 svg {
  color: #4361ee;
}

.view-all-btn {
  background-color: transparent;
  color: #4361ee;
  border: none;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 0.2s ease;
}

.view-all-btn:hover {
  color: #3a0ca3;
}

/* Skills Grid */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1rem;
}

.skill-card {
  background-color: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  cursor: pointer;
  border: 1px solid #f0f0f0;
}

.skill-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.skill-header {
  padding: 1.25rem;
  color: white;
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.skill-icon {
  background-color: rgba(255, 255, 255, 0.2);
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}

.skill-level {
  background-color: rgba(0, 0, 0, 0.2);
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

.skill-content {
  padding: 1.25rem;
}

.skill-content h3 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-size: 1.125rem;
  font-weight: 700;
}

.skill-content p {
  color: #666;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.skill-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.75rem;
  color: #666;
}

.skill-duration {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* Activity Section */
.activity-section {
  grid-column: 2;
  grid-row: 1 / span 2;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.activity-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 10px;
  transition: all 0.2s ease;
}

.activity-item:hover {
  background-color: #f0f0f0;
}

.activity-icon {
  color: #4361ee;
  font-size: 0.75rem;
  padding-top: 0.25rem;
}

.activity-content {
  flex: 1;
}

.activity-content p {
  margin: 0 0 0.25rem 0;
  font-size: 0.875rem;
}

.activity-time {
  font-size: 0.75rem;
  color: #888;
}

/* Announcements Section */
.announcements-section {
  grid-column: 1;
  grid-row: 2;
}

.announcements-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.announcement-item {
  background-color: #f9f9f9;
  border-radius: 10px;
  padding: 1.25rem;
  position: relative;
  border-left: 4px solid #4361ee;
}

.announcement-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background-color: #4361ee;
  color: white;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  opacity: 0;
}

.announcement-badge.new {
  opacity: 1;
}

.announcement-item h3 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-size: 1rem;
  font-weight: 700;
}

.announcement-item p {
  margin: 0;
  font-size: 0.875rem;
  color: #666;
  line-height: 1.5;
}

/* Responsive adjustments */
@media (max-width: 1024px) {
  .hero-section {
    flex-direction: column;
    padding: 2rem;
  }
  
  .hero-content {
    max-width: 100%;
    margin-bottom: 2rem;
  }
  
  .hero-illustration {
    max-width: 60%;
  }
  
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  
  .activity-section, 
  .announcements-section {
    grid-column: 1;
    grid-row: auto;
  }
  
  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .hero-stats {
    flex-direction: column;
  }
  
  .hero-section h1 {
    font-size: 2rem;
  }
  
  .hero-subtitle {
    font-size: 1rem;
    max-width: 100%;
  }
  
  .skills-grid {
    grid-template-columns: 1fr;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
}
</style>
