<template>
  <div class="dashboard" :class="{ 'dark-theme': isDark }">
    <!-- Hero Section with Enhanced Animation -->
    <section class="hero-section">
      <div class="container">
        <div class="hero-content">
          <h1 class="hero-title">
            Welcome to <span class="gradient-text">SkillSwap</span
            ><span class="welcome-name">, {{ user.name }}</span>
          </h1>
          <p class="hero-subtitle">
            Discover, learn, and share skills with our global community
          </p>

          <div class="hero-stats">
            <div class="stat-card points-card">
              <div class="stat-icon">
                <font-awesome-icon icon="coins" class="fa-bounce" />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ user.skillPoints || 0 }}</span>
                <span class="stat-label">SkillPoints</span>
              </div>
              <button class="stat-action" @click="openSendPointsModal">
                <font-awesome-icon icon="paper-plane" /> Send
              </button>
            </div>

            <div class="stat-card sessions-card">
              <div class="stat-icon">
                <font-awesome-icon icon="calendar-check" class="fa-beat-fade" />
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
          <img src="@/assets/images/skill-sharing.svg" alt="Skill sharing" />
          <div class="hero-shapes">
            <div class="shape shape-1"></div>
            <div class="shape shape-2"></div>
            <div class="shape shape-3"></div>
            <div class="shape shape-4"></div>
            <div class="shape shape-5"></div>
          </div>
        </div>
      </div>
      
      <!-- Animated Wave Divider -->
      <div class="wave-divider">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path fill="#f8fafc" fill-opacity="1" d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,149.3C960,160,1056,160,1152,138.7C1248,117,1344,75,1392,53.3L1440,32L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path>
        </svg>
      </div>
    </section>

    <!-- Quick Actions with Improved Hover Effects -->
    <section class="quick-actions">
      <div class="container">
        <div class="actions-grid">
          <button
            class="action-button find-button"
            @click="$router.push('/search')"
          >
            <div class="action-icon">
              <font-awesome-icon icon="search" />
            </div>
            <span>Find Skills</span>
            <div class="action-shine"></div>
          </button>

          <button
            class="action-button share-button"
            @click="$router.push('/upload-video')"
          >
            <div class="action-icon">
              <font-awesome-icon icon="video" />
            </div>
            <span>Share Video</span>
            <div class="action-shine"></div>
          </button>

          <button
            class="action-button jobs-button"
            @click="$router.push('/jobs')"
          >
            <div class="action-icon">
              <font-awesome-icon icon="briefcase" />
            </div>
            <span>Browse Jobs</span>
            <div class="action-shine"></div>
          </button>

          <button
            class="action-button chat-button"
            @click="$router.push('/chat')"
          >
            <div class="action-icon">
              <font-awesome-icon icon="comments" />
            </div>
            <span>Messages</span>
            <div class="action-shine"></div>
          </button>
        </div>
      </div>
    </section>

    <!-- Main Dashboard Content with Glass Morphism -->
    <div class="container">
      <div class="theme-toggle" @click="toggleDark">
        <font-awesome-icon :icon="isDark ? 'sun' : 'moon'" />
      </div>
      
      <div class="dashboard-grid">
        <!-- Featured Skills Section -->
        <section class="dashboard-section skills-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="star" class="fa-spin-pulse" /> Featured Skills</h2>
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
              :style="{ animationDelay: index * 0.1 + 's' }"
            >
              <div
                class="skill-header"
                :style="{
                  backgroundImage: `linear-gradient(135deg, ${getColorBySkillCategory(skill.category || 'programming')} 0%, ${getColorGradient(skill.category || 'programming')} 100%)`,
                }"
              >
                <div class="skill-icon">
                  <font-awesome-icon :icon="getIconForSkill(skill.name)" />
                </div>
                <div class="skill-level">{{ skill.level || "Beginner" }}</div>
              </div>
              <div class="skill-content">
                <h3>{{ skill.name }}</h3>
                <p>{{ skill.description }}</p>
                <div class="skill-meta">
                  <span class="skill-duration">
                    <font-awesome-icon icon="clock" />
                    {{ skill.duration || "4 weeks" }}
                  </span>
                  <span class="skill-users">
                    <font-awesome-icon icon="users" />
                    {{ getRandomUserCount(skill) }} users
                  </span>
                </div>
              </div>
              <div class="skill-card-shine"></div>
            </div>
          </div>
        </section>

        <!-- Recent Activity Section with Improved Animations -->
        <section class="dashboard-section activity-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="history" /> Recent Activity</h2>
          </div>

          <div class="activity-list">
            <div
              class="activity-item"
              v-for="(activity, index) in recentActivities"
              :key="index"
              :class="{ highlight: index === 0 }"
              :style="{ animationDelay: index * 0.1 + 's' }"
            >
              <div
                class="activity-icon"
                :class="`activity-type-${getActivityType(activity)}`"
              >
                <font-awesome-icon :icon="getActivityIcon(activity)" />
              </div>
              <div class="activity-content">
                <p v-html="formatActivityText(activity)"></p>
                <span class="activity-time">{{ randomTimeAgo() }}</span>
              </div>
            </div>
          </div>
        </section>

        <!-- Announcements Section with Card Flip Effect -->
        <section class="dashboard-section announcements-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="bullhorn" /> Announcements</h2>
          </div>

          <div class="announcements-list">
            <div
              class="announcement-item"
              v-for="(announcement, index) in announcements"
              :key="index"
              :style="{ animationDelay: index * 0.1 + 's' }"
            >
              <div class="announcement-badge" :class="{ new: index === 0 }">
                {{ index === 0 ? "NEW" : "" }}
              </div>
              <h3>{{ announcement.title }}</h3>
              <p>{{ announcement.message }}</p>
              <div class="announcement-footer">
                <button class="read-more-btn">
                  Learn more <font-awesome-icon icon="chevron-right" />
                </button>
                <span class="announcement-date">{{
                  getAnnouncementDate(index)
                }}</span>
              </div>
              <div class="announcement-shine"></div>
            </div>
          </div>
        </section>

        <!-- Trending Skills with Progress Bars -->
        <section class="dashboard-section trending-section">
          <div class="section-header">
            <h2><font-awesome-icon icon="fire" class="fa-shake" /> Trending Now</h2>
          </div>

          <div class="trending-list">
            <div
              class="trending-item"
              v-for="(item, index) in trendingSkills"
              :key="index"
              @click="viewSkillDetails(item)"
              :style="{ animationDelay: index * 0.1 + 's' }"
            >
              <div class="trending-rank">{{ index + 1 }}</div>
              <div class="trending-name">{{ item.name }}</div>
              <div class="trending-progress">
                <div class="progress-bar" :style="{ width: getProgressWidth(item.users) + '%' }"></div>
              </div>
              <div class="trending-stats">
                <div class="trending-users">
                  <font-awesome-icon icon="user" /> {{ item.users }}
                </div>
                <div
                  class="trending-growth"
                  :class="{ positive: item.growth > 0 }"
                >
                  <font-awesome-icon
                    :icon="item.growth > 0 ? 'arrow-up' : 'arrow-down'"
                  />
                  {{ Math.abs(item.growth) }}%
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script>
import { useDark, useToggle } from '@vueuse/core';

export default {
  name: "Dashboard",
  setup() {
    const isDark = useDark();
    const toggleDark = useToggle(isDark);
    
    return {
      isDark,
      toggleDark
    };
  },
  data() {
    return {
      upcomingSessions: 2, // This would come from your API in a real app
      featuredSkills: [
        {
          name: "Go Programming",
          description:
            "Learn the basics of Go programming language and build efficient applications.",
          category: "programming",
          level: "Beginner",
          duration: "4 weeks",
        },
        {
          name: "Vue.js",
          description:
            "Frontend development with Vue framework, components, and state management.",
          category: "programming",
          level: "Intermediate",
          duration: "6 weeks",
        },
        {
          name: "Guitar Lessons",
          description:
            "Play your favorite tunes on guitar with proper technique and music theory.",
          category: "music",
          level: "Beginner",
          duration: "8 weeks",
        },
        {
          name: "Spanish Language",
          description:
            "Learn conversational Spanish for travel, business, or personal enrichment.",
          category: "language",
          level: "Beginner",
          duration: "10 weeks",
        },
      ],
      recentActivities: [
        "Alice taught <strong>Spanish</strong> to Bob",
        "Charlie added a new skill: <strong>Italian Cooking</strong>",
        "You updated your profile information",
        "David completed a <strong>Python</strong> session with Emma",
        "You earned <strong>15 SkillPoints</strong> for your JavaScript tutorial",
        "Maria left a <strong>5-star review</strong> on your teaching session",
      ],
      announcements: [
        {
          title: "New Feature: Live Sessions",
          message:
            "We are excited to announce live skill sessions are now available! Connect with other users in real-time for interactive learning.",
        },
        {
          title: "Scheduled Maintenance",
          message:
            "Scheduled maintenance will occur on Saturday at 2 PM. The platform may be unavailable for a short period.",
        },
        {
          title: "Community Spotlight",
          message:
            "Every month we highlight exceptional community members. Nominate yourself or others for the next spotlight!",
        },
      ],
      trendingSkills: [
        { name: "React.js", users: 2450, growth: 15 },
        { name: "Digital Marketing", users: 1830, growth: 12 },
        { name: "UI/UX Design", users: 1650, growth: 8 },
        { name: "Machine Learning", users: 1425, growth: 20 },
        { name: "Video Editing", users: 1380, growth: -5 },
      ],
    };
  },
  computed: {
    user() {
      return (
        this.$store.state.user || {
          name: "Guest",
          skillPoints: 0,
        }
      );
    },
  },
  methods: {
    viewSkillDetails(skill) {
      // In a real app, you would navigate to a skill details page
      this.$router.push({
        name: "Search",
        query: { q: skill.name },
      });
    },
    randomTimeAgo() {
      const times = [
        "Just now",
        "5 minutes ago",
        "2 hours ago",
        "Yesterday",
        "2 days ago",
      ];
      return times[Math.floor(Math.random() * times.length)];
    },
    getColorBySkillCategory(category) {
      const categoryColors = {
        programming: "#4361ee",
        design: "#3a86ff",
        language: "#ff9f1c",
        music: "#7209b7",
        cooking: "#e63946",
        business: "#2a9d8f",
        default: "#4361ee",
      };

      return categoryColors[category.toLowerCase()] || categoryColors.default;
    },
    getColorGradient(category) {
      const gradients = {
        programming: "#3a0ca3",
        design: "#4895ef",
        language: "#f3722c",
        music: "#560bad",
        cooking: "#ff595e",
        business: "#1a759f",
        default: "#3a0ca3",
      };

      return gradients[category.toLowerCase()] || gradients.default;
    },
    getIconForSkill(skillName) {
      const skillIcons = {
        "Go Programming": ["fab", "golang"],
        "Vue.js": ["fab", "vuejs"],
        "Guitar Lessons": "guitar",
        Python: ["fab", "python"],
        Cooking: "utensils",
        "Spanish Language": "language",
        "React.js": ["fab", "react"],
        "Digital Marketing": "bullhorn",
        "UI/UX Design": "palette",
        "Machine Learning": "brain",
        "Video Editing": "film",
      };

      return skillIcons[skillName] || "graduation-cap"; // Default icon
    },
    getRandomUserCount(skill) {
      // In a real app, this would come from the backend
      // Make random number deterministic based on skill name for consistency
      const nameHash = skill.name
        .split("")
        .reduce((acc, char) => acc + char.charCodeAt(0), 0);
      return (nameHash % 900) + 100; // Return number between 100-999
    },
    getActivityType(activity) {
      if (activity.includes("taught")) return "teach";
      if (activity.includes("added")) return "add";
      if (activity.includes("updated")) return "update";
      if (activity.includes("completed")) return "complete";
      if (activity.includes("earned")) return "earn";
      if (activity.includes("review")) return "review";
      return "default";
    },
    getActivityIcon(activity) {
      const activityType = this.getActivityType(activity);
      const icons = {
        teach: "chalkboard-teacher",
        add: "plus-circle",
        update: "user-edit",
        complete: "check-circle",
        earn: "coins",
        review: "star",
        default: "circle-dot",
      };

      return icons[activityType] || icons.default;
    },
    formatActivityText(activity) {
      return activity;
    },
    getAnnouncementDate(index) {
      const days = ["Today", "3 days ago", "1 week ago", "2 weeks ago"];
      return days[index] || days[0];
    },
    openSendPointsModal() {
      // You could navigate to the transactions page
      this.$router.push("/transactions");
      // Or implement a modal directly in this component
    },
    getProgressWidth(users) {
      // Calculate width percentage based on highest user count
      const maxUsers = Math.max(...this.trendingSkills.map(skill => skill.users));
      return (users / maxUsers) * 100;
    }
  },
};
</script>

<style scoped>
/* Base Styles */
.dashboard {
  font-family: var(--font-family-sans, 'Inter', sans-serif);
  color: var(--dark, #1e293b);
  background-color: #f8fafc;
  overflow-x: hidden;
  transition: background-color 0.3s ease, color 0.3s ease;
  position: relative;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
  position: relative;
}

/* Dark Theme */
.dark-theme {
  background-color: #0f172a;
  color: #e2e8f0;
}

.dark-theme .dashboard-section {
  background-color: #1e293b;
  border-color: #334155;
}

.dark-theme .action-button {
  background-color: #1e293b;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
}

.dark-theme .skill-card {
  background-color: #1e293b;
  border-color: #334155;
}

.dark-theme .skill-content h3 {
  color: #e2e8f0;
}

.dark-theme .skill-content p {
  color: #94a3b8;
}

.dark-theme .activity-item {
  background-color: #334155;
}

.dark-theme .announcement-item {
  background-color: #334155;
}

.dark-theme .trending-item {
  background-color: #334155;
}

.dark-theme .section-header {
  border-bottom-color: #334155;
}

.dark-theme .wave-divider svg path {
  fill: #0f172a;
}

/* Theme Toggle */
.theme-toggle {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background: linear-gradient(135deg, #4361ee, #3a0ca3);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 100;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  font-size: 1.25rem;
}

.theme-toggle:hover {
  transform: translateY(-5px) rotate(15deg);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
}

/* Hero Section with Enhanced Animation */
.hero-section {
  background: linear-gradient(135deg, #4361ee 0%, #3a0ca3 100%);
  padding: 5rem 0 8rem;
  color: white;
  position: relative;
  overflow: hidden;
}

.hero-section .container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  z-index: 2;
}

.hero-content {
  max-width: 60%;
  position: relative;
  z-index: 2;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 1rem;
  animation: fadeInUp 0.8s ease forwards;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.welcome-name {
  opacity: 0;
  animation: fadeInRight 0.8s ease forwards 0.4s;
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
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.8s ease forwards 0.2s;
  max-width: 80%;
  margin-bottom: 2rem;
  text-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
}

.hero-stats {
  display: flex;
  gap: 1.5rem;
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.8s ease forwards 0.4s;
}

.stat-card {
  background-color: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 1rem;
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.stat-card:hover {
  background-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-5px) scale(1.02);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.15);
}

.points-card {
  background-color: rgba(249, 200, 14, 0.15);
}

.sessions-card {
  background-color: rgba(58, 12, 163, 0.3);
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
  transition: all 0.3s ease;
}

.points-card .stat-icon {
  background-color: rgba(249, 200, 14, 0.6);
  color: #000;
}

.stat-card:hover .stat-icon {
  transform: scale(1.1) rotate(5deg);
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
  transition: all 0.3s ease;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}

.stat-action:hover {
  background-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.stat-action::after {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(
    to right,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.3) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  transform: rotate(45deg);
  transition: all 0.5s ease;
  opacity: 0;
}

.stat-action:hover::after {
  animation: shine 1.5s ease;
  opacity: 1;
}

.hero-illustration {
  max-width: 40%;
  position: relative;
  z-index: 1;
  opacity: 0;
  transform: translateX(50px);
  animation: fadeInRight 0.8s ease forwards 0.3s;
}

.hero-illustration img {
  max-width: 100%;
  height: auto;
  animation: float 6s ease-in-out infinite;
  position: relative;
  z-index: 2;
  filter: drop-shadow(0 10px 15px rgba(0, 0, 0, 0.2));
}

.hero-shapes {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.2;
}

.shape-1 {
  width: 300px;
  height: 300px;
  background-color: rgba(255, 255, 255, 0.1);
  top: -100px;
  right: -50px;
  animation: pulse 8s ease-in-out infinite alternate;
}

.shape-2 {
  width: 200px;
  height: 200px;
  background-color: rgba(249, 200, 14, 0.2);
  bottom: -80px;
  right: 10%;
  animation: pulse 12s ease-in-out infinite alternate-reverse;
}

.shape-3 {
  width: 150px;
  height: 150px;
  background-color: rgba(58, 134, 255, 0.2);
  bottom: 20%;
  left: 5%;
  animation: pulse 10s ease-in-out infinite alternate;
}

.shape-4 {
  width: 100px;
  height: 100px;
  background-color: rgba(255, 255, 255, 0.15);
  top: 30%;
  left: 20%;
  animation: float 15s ease-in-out infinite;
}

.shape-5 {
  width: 80px;
  height: 80px;
  background-color: rgba(249, 200, 14, 0.15);
  top: 20%;
  right: 25%;
  animation: float 12s ease-in-out infinite reverse;
}

/* Wave Divider */
.wave-divider {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  line-height: 0;
  direction: ltr;
  overflow: hidden;
}

.wave-divider svg {
  display: block;
  width: 100%;
  height: 70px;
  transform-origin: bottom;
  animation: wave 20s linear infinite;
}

/* Quick Actions with Improved Hover Effects */
.quick-actions {
  margin-top: -2rem;
  margin-bottom: 3rem;
  position: relative;
  z-index: 5;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

.action-button {
  height: 100%;
  background-color: white;
  border: none;
  border-radius: 1rem;
  padding: 1.5rem 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05);
  text-align: center;
  position: relative;
  overflow: hidden;
}

.action-icon {
  width: 60px;
  height: 60px;
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  margin-bottom: 0.5rem;
  transition: all 0.3s ease;
  position: relative;
  z-index: 2;
}

.find-button .action-icon {
  background-color: #e0f2fe;
  color: #0284c7;
}

.share-button .action-icon {
  background-color: #fef3c7;
  color: #d97706;
}

.jobs-button .action-icon {
  background-color: #f3e8ff;
  color: #7c3aed;
}

.chat-button .action-icon {
  background-color: #dcfce7;
  color: #16a34a;
}

.action-button span {
  font-weight: 600;
  font-size: 1.1rem;
  position: relative;
  z-index: 2;
}

.action-button:hover {
  transform: translateY(-10px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
}

.action-button:hover .action-icon {
  transform: scale(1.1) rotate(10deg);
}

.find-button:hover .action-icon {
  background-color: #0284c7;
  color: white;
}

.share-button:hover .action-icon {
  background-color: #d97706;
  color: white;
}

.jobs-button:hover .action-icon {
  background-color: #7c3aed;
  color: white;
}

.chat-button:hover .action-icon {
  background-color: #16a34a;
  color: white;
}

/* Shine effect for action buttons */
.action-shine {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.4) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  z-index: 1;
  transform: translateX(-100%) rotate(45deg);
  animation: none;
  pointer-events: none;
}

.action-button:hover .action-shine {
  animation: shine 1.5s ease-in-out;
}

/* Dashboard Grid Layout with Glass Morphism */
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 1.5rem;
  margin-bottom: 4rem;
}

.dashboard-section {
  background-color: white;
  border-radius: 1rem;
  padding: 1.5rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  height: 100%;
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.03);
  overflow: hidden;
  animation: fadeInUp 0.5s ease forwards;
  opacity: 0;
  transform: translateY(20px);
  backdrop-filter: blur(10px);
}

.dashboard-section:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.skills-section {
  grid-column: span 8;
}

.activity-section {
  grid-column: span 4;
  grid-row: span 2;
}

.announcements-section {
  grid-column: span 8;
}

.trending-section {
  grid-column: span 4;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid #f0f0f0;
}

.section-header h2 {
  margin-bottom: 0;
  font-size: 1.25rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--dark);
}

.section-header h2 svg {
  color: var(--primary-color, #4361ee);
}

.view-all-btn {
  background-color: transparent;
  color: var(--primary-color, #4361ee);
  border: none;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
}

.view-all-btn:hover {
  color: var(--primary-dark, #3a0ca3);
  transform: translateX(3px);
}

.view-all-btn::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background-color: var(--primary-color, #4361ee);
  transition: width 0.3s ease;
}

.view-all-btn:hover::after {
  width: 100%;
}

/* Skills Grid with Enhanced Cards */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.25rem;
}

.skill-card {
  background-color: white;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  cursor: pointer;
  border: 1px solid #f0f0f0;
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  animation: fadeInUp 0.5s ease forwards;
  opacity: 0;
  transform: translateY(20px);
}

.skill-card:hover {
  transform: translateY(-5px) scale(1.02);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
}

.skill-header {
  padding: 1.5rem;
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
  transition: all 0.5s ease;
}

.skill-card:hover .skill-icon {
  transform: rotate(360deg);
}

.skill-level {
  background-color: rgba(0, 0, 0, 0.2);
  font-size: 0.75rem;
  padding: 0.25rem 0.75rem;
  border-radius: 50px;
  font-weight: 500;
}

.skill-content {
  padding: 1.25rem;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.skill-content h3 {
  margin: 0 0 0.75rem 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--dark);
}

.skill-content p {
  margin: 0 0 1rem 0;
  color: var(--medium, #64748b);
  font-size: 0.9rem;
  line-height: 1.5;
  flex: 1;
}

.skill-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  color: var(--medium, #64748b);
}

.skill-duration,
.skill-users {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

/* Shine effect for skill cards */
.skill-card-shine {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.4) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  z-index: 1;
  transform: translateX(-100%) rotate(45deg);
  pointer-events: none;
}

.skill-card:hover .skill-card-shine {
  animation: shine 1.5s ease-in-out;
}

/* Activity Section with Improved Items */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.activity-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background-color: #f9fafb;
  border-radius: 1rem;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
  animation: fadeInRight 0.5s ease forwards;
  opacity: 0;
  transform: translateX(20px);
}

.activity-item:hover {
  background-color: #f3f4f6;
  transform: translateX(5px);
}

.activity-item.highlight {
  border-left-color: #4361ee;
}

.activity-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  flex-shrink: 0;
  color: white;
  transition: all 0.3s ease;
}

.activity-item:hover .activity-icon {
  transform: scale(1.1);
}

.activity-type-teach {
  background-color: #4361ee;
}

.activity-type-add {
  background-color: #10b981;
}

.activity-type-update {
  background-color: #3b82f6;
}

.activity-type-complete {
  background-color: #8b5cf6;
}

.activity-type-earn {
  background-color: #f59e0b;
}

.activity-type-review {
  background-color: #f43f5e;
}

.activity-type-default {
  background-color: #6b7280;
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-content p {
  margin: 0 0 0.25rem 0;
  font-size: 0.9rem;
  line-height: 1.5;
}

.activity-time {
  font-size: 0.75rem;
  color: var(--medium, #64748b);
}

/* Announcements Section with Card Flip Effect */
.announcements-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.25rem;
}

.announcement-item {
  background-color: #f9fafb;
  border-radius: 1rem;
  padding: 1.25rem;
  position: relative;
  transition: all 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
  animation: fadeInUp 0.5s ease forwards;
  opacity: 0;
  transform: translateY(20px);
  overflow: hidden;
}

.announcement-item:hover {
  background-color: #f3f4f6;
  transform: translateY(-5px) scale(1.02);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05);
}

.announcement-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background-color: transparent;
  color: transparent;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  z-index: 2;
}

.announcement-badge.new {
  background-color: #ef4444;
  color: white;
  animation: pulse 2s infinite;
}

.announcement-item h3 {
  margin: 0 0 0.75rem 0;
  font-size: 1.1rem;
  color: var(--dark);
  padding-right: 3rem;
}

.announcement-item p {
  margin: 0;
  font-size: 0.9rem;
  color: var(--medium, #64748b);
  line-height: 1.5;
  flex: 1;
}

.announcement-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1rem;
  padding-top: 0.75rem;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.read-more-btn {
  background: none;
  border: none;
  color: var(--primary-color, #4361ee);
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  gap: 0.3rem;
  transition: all 0.3s ease;
}

.read-more-btn:hover {
  color: var(--primary-dark, #3a0ca3);
  transform: translateX(3px);
}

.announcement-date {
  font-size: 0.75rem;
  color: var(--medium, #64748b);
}

/* Shine effect for announcements */
.announcement-shine {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.4) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  z-index: 1;
  transform: translateX(-100%) rotate(45deg);
  pointer-events: none;
}

.announcement-item:hover .announcement-shine {
  animation: shine 1.5s ease-in-out;
}

/* Trending Skills Section with Progress Bars */
.trending-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.trending-item {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  padding: 0.75rem 1rem;
  background-color: #f9fafb;
  border-radius: 0.75rem;
  transition: all 0.3s ease;
  cursor: pointer;
  animation: fadeInRight 0.5s ease forwards;
  opacity: 0;
  transform: translateX(20px);
}

.trending-item:hover {
  background-color: #f3f4f6;
  transform: translateX(5px);
}

.trending-rank {
  width: 28px;
  height: 28px;
  background-color: var(--primary-light, #eef2ff);
  color: var(--primary-color, #4361ee);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.85rem;
  margin-right: 1rem;
}

.trending-name {
  flex: 1;
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 0.5rem;
}

.trending-progress {
  flex-basis: 100%;
  height: 6px;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
  margin: 0.5rem 0;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #4361ee, #3a0ca3);
  border-radius: 3px;
  transition: width 1s ease;
}

.trending-stats {
  display: flex;
  gap: 1rem;
  width: 100%;
  justify-content: flex-end;
}

.trending-users,
.trending-growth {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.8rem;
  color: var(--medium, #64748b);
}

.trending-growth {
  color: #ef4444;
}

.trending-growth.positive {
  color: #10b981;
}

/* Animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInRight {
  from {
    opacity: 0;
    transform: translateX(30px);
  }

  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes float {
  0% {
    transform: translateY(0px);
  }

  50% {
    transform: translateY(-20px);
  }

  100% {
    transform: translateY(0px);
  }
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 0.8;
  }

  50% {
    transform: scale(1.05);
    opacity: 1;
  }

  100% {
    transform: scale(1);
    opacity: 0.8;
  }
}

@keyframes shine {
  0% {
    transform: translateX(-100%) rotate(45deg);
  }
  
  100% {
    transform: translateX(100%) rotate(45deg);
  }
}

@keyframes wave {
  0% {
    transform: translateX(0) scaleY(1);
  }
  50% {
    transform: translateX(-2%) scaleY(0.95);
  }
  100% {
    transform: translateX(0) scaleY(1);
  }
}

/* Responsive adjustments */
@media (max-width: 1200px) {
  .dashboard-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .skills-section,
  .activity-section,
  .announcements-section,
  .trending-section {
    grid-column: span 1;
  }

  .activity-section {
    grid-row: auto;
  }
}

@media (max-width: 992px) {
  .hero-section .container {
    flex-direction: column;
    text-align: center;
  }

  .hero-content {
    max-width: 100%;
    margin-bottom: 3rem;
  }

  .hero-subtitle {
    margin-left: auto;
    margin-right: auto;
  }

  .hero-illustration {
    max-width: 70%;
  }

  .hero-stats {
    justify-content: center;
  }

  .actions-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .hero-title {
    font-size: 2.25rem;
  }

  .hero-stats {
    flex-direction: column;
  }

  .announcements-list {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 576px) {
  .actions-grid {
    grid-template-columns: 1fr;
  }

  .hero-section {
    padding: 3rem 0;
  }

  .hero-title {
    font-size: 1.75rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .skills-grid {
    grid-template-columns: 1fr;
  }
  
  .theme-toggle {
    bottom: 1rem;
    right: 1rem;
    width: 2.5rem;
    height: 2.5rem;
  }
}
</style>