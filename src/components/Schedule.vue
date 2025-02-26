<template>
  <div class="schedule-container">
    <h2>My Schedules</h2>
    <form @submit.prevent="createSchedule">
      <input type="datetime-local" v-model="newSchedule.startTime" required />
      <input type="datetime-local" v-model="newSchedule.endTime" required />
      <button type="submit">Schedule Session</button>
    </form>
    <div v-if="schedules.length">
      <ul>
        <li v-for="(schedule, index) in schedules" :key="index">
          Session on Skill ID: {{ schedule.skill_id }} from {{ schedule.startTime }} to {{ schedule.endTime }}
        </li>
      </ul>
    </div>
    <div v-else>
      <p>No scheduled sessions found.</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Schedule",
  data() {
    return {
      newSchedule: {
        skill_id: 1, // For demonstration, a default skill ID.
        startTime: "",
        endTime: ""
      },
      schedules: []
    };
  },
  created() {
    this.fetchSchedules();
  },
  methods: {
    async createSchedule() {
      try {
        const response = await axios.post("/api/schedule", this.newSchedule);
        this.schedules.push(response.data);
      } catch (error) {
        console.error("Error creating schedule:", error);
        if (this.$toast) {
          this.$toast.error("Failed to create schedule.");
        }
      }
    },
    async fetchSchedules() {
      try {
        const response = await axios.get("/api/schedule");
        this.schedules = response.data;
      } catch (error) {
        console.error("Error fetching schedules:", error);
      }
    }
  }
};
</script>

<style scoped>
.schedule-container {
  padding: 2rem;
}
form {
  margin-bottom: 1rem;
  display: flex;
  gap: 1rem;
}
</style>

