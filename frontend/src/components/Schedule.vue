<template>
  <div class="schedule-container">
    <h2>My Schedules</h2>
    <form @submit.prevent="createSchedule" class="schedule-form">
      <label>
        Start Time:
        <input type="datetime-local" v-model="newSchedule.startTime" required />
      </label>
      <label>
        End Time:
        <input type="datetime-local" v-model="newSchedule.endTime" required />
      </label>
      <button type="submit">Schedule Session</button>
    </form>
    <div v-if="scheduleError" class="error">{{ scheduleError }}</div>
    <div v-if="scheduleLoading" class="loading">Processing...</div>
    <div v-if="schedules.length">
      <ul class="schedule-list">
        <li v-for="(schedule, index) in schedules" :key="index">
          Session on Skill ID: {{ schedule.skill_id }} from
          {{ formatDate(schedule.startTime) }} to
          {{ formatDate(schedule.endTime) }}
        </li>
      </ul>
    </div>
    <div v-else-if="!scheduleLoading && schedulesFetched">
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
        skill_id: 1, // Default value for demonstration
        startTime: "",
        endTime: "",
      },
      schedules: [],
      scheduleError: null,
      scheduleLoading: false,
      schedulesFetched: false,
    };
  },
  created() {
    this.fetchSchedules();
  },
  methods: {
    async createSchedule() {
      this.scheduleError = null;
      this.scheduleLoading = true;
      try {
        const response = await axios.post("/api/schedule", this.newSchedule);
        this.schedules.push(response.data);
      } catch (error) {
        console.error("Error creating schedule:", error);
        this.scheduleError =
          "Failed to create schedule. Please ensure the session is in the future and try again.";
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
        console.error("Error fetching schedules:", error);
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
};
</script>

<style scoped>
.schedule-container {
  padding: 2rem;
}
.schedule-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1rem;
}
.schedule-form label {
  display: flex;
  flex-direction: column;
  font-size: 1rem;
}
.schedule-form input {
  padding: 0.5rem;
  font-size: 1rem;
}
button {
  align-self: flex-start;
  padding: 0.5rem 1rem;
  font-size: 1rem;
}
.loading {
  font-style: italic;
  margin-bottom: 1rem;
}
.error {
  color: red;
  margin-bottom: 1rem;
}
.schedule-list {
  list-style: none;
  padding: 0;
}
.schedule-list li {
  padding: 0.5rem;
  border-bottom: 1px solid #ccc;
}
</style>
