<template>
  <div class="schedule-container">
    <header class="schedule-header">
      <div class="container">
        <h1>My Schedule</h1>
        <p>Manage your skill exchange sessions and availability</p>
      </div>
    </header>

    <div class="container schedule-content">
      <div class="card schedule-form-card">
        <h2>Create New Session</h2>
        <form @submit.prevent="createSchedule" class="schedule-form">
          <div class="form-group">
            <label for="startTime">Start Time</label>
            <input id="startTime" type="datetime-local" v-model="newSchedule.startTime" required />
          </div>
          <div class="form-group">
            <label for="endTime">End Time</label>
            <input id="endTime" type="datetime-local" v-model="newSchedule.endTime" required />
          </div>
          <button type="submit" class="schedule-btn">
            <span class="btn-icon">📅</span>
            Schedule Session
          </button>
        </form>
      </div>

      <div v-if="scheduleError" class="error-message">
        {{ scheduleError }}
      </div>

      <div v-if="scheduleLoading" class="loading-indicator">
        <div class="loading-spinner"></div>
        <span>Processing your request...</span>
      </div>

      <div class="card sessions-card" v-if="schedules.length">
        <h2>Upcoming Sessions</h2>
        <ul class="schedule-list">
          <li v-for="(schedule, index) in schedules" :key="index" class="schedule-item">
            <div class="session-icon">📚</div>
            <div class="session-details">
              <h3>Skill Exchange Session</h3>
              <p class="session-skill">Skill ID: {{ schedule.skill_id }}</p>
              <div class="session-time">
                <div class="time-block">
                  <span class="time-label">Starts:</span>
                  <span class="time-value">{{
                    formatDate(schedule.startTime)
                  }}</span>
                </div>
                <div class="time-block">
                  <span class="time-label">Ends:</span>
                  <span class="time-value">{{
                    formatDate(schedule.endTime)
                  }}</span>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>

      <div v-else-if="!scheduleLoading && schedulesFetched" class="card empty-card">
        <div class="empty-state">
          <div class="empty-icon">📅</div>
          <h3>No Sessions Scheduled</h3>
          <p>You don't have any upcoming skill exchange sessions.</p>
          <p>Use the form above to schedule your first session!</p>
        </div>
      </div>
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
      // Validate date inputs before submission
      if (!this.newSchedule.startTime || !this.newSchedule.endTime) {
        this.scheduleError = "Please set both start and end times.";
        return;
      }

      const startTime = new Date(this.newSchedule.startTime);
      const endTime = new Date(this.newSchedule.endTime);
      const now = new Date();

      if (startTime <= now) {
        this.scheduleError = "Start time must be in the future.";
        return;
      }

      if (endTime <= startTime) {
        this.scheduleError = "End time must be after start time.";
        return;
      }

      // Clear previous errors before submission
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
/* General page styling */
.schedule-container {
  font-family: "Helvetica Neue", Arial, sans-serif;
  color: var(--dark);
  background: linear-gradient(to bottom, #f4f7f9, #e8eef2);
  min-height: 100vh;
  padding-bottom: var(--space-8);
}

.container {
  width: 90%;
  max-width: 1000px;
  margin: 0 auto;
}

/* Header section */
.schedule-header {
  background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
  color: #fff;
  padding: 3rem 0;
  text-align: center;
  margin-bottom: var(--space-6);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.schedule-header h1 {
  margin: 0;
  font-size: 2.5rem;
  letter-spacing: 0.5px;
  font-weight: 700;
}

.schedule-header p {
  margin: 0.5rem 0 0;
  font-size: 1.2rem;
  opacity: 0.9;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

.schedule-content {
  position: relative;
  margin-top: -2rem;
}

/* Card styling */
.card {
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  padding: var(--space-6);
  margin-bottom: var(--space-6);
  border: 1px solid #f0f0f0;
  position: relative;
  overflow: hidden;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.08);
}

.card h2 {
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

.schedule-form-card h2::before {
  content: "📝";
  margin-right: var(--space-2);
}

/* Form styling */
.schedule-form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: var(--space-4);
}

.form-group {
  margin-bottom: var(--space-4);
}

.form-group label {
  display: block;
  font-weight: 600;
  margin-bottom: var(--space-2);
  color: var(--dark);
}

.form-group input {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid #ddd;
  border-radius: var(--radius-md);
  font-size: 1rem;
  transition:
    border 0.2s ease,
    box-shadow 0.2s ease;
}

.form-group input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
  outline: none;
}

.schedule-btn {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-full);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  height: fit-content;
  align-self: flex-end;
}

.btn-icon {
  font-size: 1.1rem;
}

.schedule-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background: var(--primary-dark);
}

/* Sessions list */
.sessions-card h2::before {
  content: "🗓️";
  margin-right: var(--space-2);
}

.schedule-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.schedule-item {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4);
  background: #f8fafc;
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-3);
  border: 1px solid #edf2f7;
  transition: all 0.2s ease;
}

.schedule-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
  border-color: var(--primary-light);
}

.session-icon {
  font-size: 2rem;
  color: var(--primary-color);
  background: rgba(79, 70, 229, 0.1);
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  flex-shrink: 0;
}

.session-details {
  flex: 1;
}

.session-details h3 {
  margin: 0 0 var(--space-1) 0;
  font-size: 1.2rem;
  color: var(--dark);
}

.session-skill {
  color: var(--medium);
  margin: 0 0 var(--space-2) 0;
  font-size: 0.9rem;
}

.session-time {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
}

.time-block {
  background: white;
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
  border: 1px solid #edf2f7;
}

.time-label {
  font-weight: 600;
  color: var(--primary-color);
  margin-right: var(--space-1);
  font-size: 0.85rem;
}

.time-value {
  color: var(--dark);
}

/* Error and loading states */
.error-message {
  background-color: rgba(239, 68, 68, 0.08);
  color: var(--error-color);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  border-left: 4px solid var(--error-color);
  margin-bottom: var(--space-4);
  font-size: 0.95rem;
  display: flex;
  align-items: center;
}

.error-message::before {
  content: "⚠️";
  margin-right: var(--space-2);
}

.loading-indicator {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  background: #f8fafc;
  border-radius: var(--radius-md);
  margin-bottom: var(--space-4);
  font-style: italic;
  color: var(--medium);
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(79, 70, 229, 0.3);
  border-radius: 50%;
  border-top-color: var(--primary-color);
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Empty state */
.empty-card {
  padding: var(--space-8) var(--space-6);
}

.empty-state {
  text-align: center;
  color: var(--medium);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: var(--space-4);
  opacity: 0.7;
}

.empty-state h3 {
  font-size: 1.5rem;
  color: var(--dark);
  margin-top: 0;
  margin-bottom: var(--space-2);
}

.empty-state p {
  margin: var(--space-1) 0;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

/* Responsive styles */
@media (max-width: 768px) {
  .schedule-form {
    grid-template-columns: 1fr;
  }

  .session-time {
    flex-direction: column;
    gap: var(--space-2);
  }

  .schedule-item {
    flex-direction: column;
    text-align: center;
  }

  .session-icon {
    margin-bottom: var(--space-2);
  }
}
</style>
