<!-- File: ./src/components/Dashboard.vue -->
<template>
  <div class="dashboard">
    <h2>Welcome to SkillSwap, {{ user.name || "Guest" }}!</h2>
    <p>Your SkillPoints: {{ user.skillPoints || 0 }}</p>

    <!-- Featured Skills Section -->
    <section class="featured-skills">
      <h3>Featured Skills</h3>
      <div class="skills-list">
        <div
          class="skill-card"
          v-for="(skill, index) in featuredSkills"
          :key="index"
        >
          <img
            :src="skill.image || defaultSkillImage"
            alt="Skill Image"
            class="skill-image"
          />
          <div class="skill-info">
            <h4>{{ skill.name }}</h4>
            <p>{{ skill.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Recent Activity Section -->
    <section class="recent-activity">
      <h3>Recent Activity</h3>
      <ul>
        <li v-for="(activity, index) in recentActivities" :key="index">
          {{ activity }}
        </li>
      </ul>
    </section>

    <!-- Announcements Section -->
    <section class="announcements">
      <h3>Announcements</h3>
      <div v-if="announcements.length">
        <div
          v-for="(announcement, index) in announcements"
          :key="index"
          class="announcement"
        >
          <h4>{{ announcement.title }}</h4>
          <p>{{ announcement.message }}</p>
        </div>
      </div>
      <div v-else>
        <p>No announcements at this time.</p>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: "Dashboard",
  data() {
    return {
      featuredSkills: [
        {
          name: "Go Programming",
          description: "Learn the basics of Go",
          image: "https://via.placeholder.com/150",
        },
        {
          name: "Vue.js",
          description: "Frontend development with Vue",
          image: "https://via.placeholder.com/150",
        },
        {
          name: "Guitar Lessons",
          description: "Play your favorite tunes",
          image: "https://via.placeholder.com/150",
        },
      ],
      recentActivities: [
        "Alice taught Spanish to Bob",
        "Charlie added a new skill: Cooking",
        "Test User updated their profile",
      ],
      announcements: [
        {
          title: "New Feature",
          message: "We are excited to announce a new feature: live chat!",
        },
        {
          title: "Maintenance",
          message: "Scheduled maintenance on Saturday at 2 PM.",
        },
      ],
      defaultSkillImage: "https://via.placeholder.com/150",
    };
  },
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
  max-width: 1200px;
  margin: 0 auto;
}

.featured-skills,
.recent-activity,
.announcements {
  margin-top: 2rem;
}

.skills-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.skill-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  width: calc(33.33% - 1rem);
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.skill-image {
  width: 100%;
  height: auto;
  border-radius: 4px;
}

.skill-info {
  margin-top: 0.5rem;
}

.recent-activity ul {
  list-style: none;
  padding: 0;
}

.recent-activity li {
  background: #f9f9f9;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  border-radius: 4px;
}

.announcements .announcement {
  border: 1px solid #eee;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  background: #fdfdfd;
}
</style>
