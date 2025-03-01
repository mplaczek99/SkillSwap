<template>
  <div class="feedback-page">
    <div class="container">
      <h2>Ratings & Feedback</h2>
      <p class="subtitle">
        Build trust in the community by providing and receiving feedback.
      </p>

      <!-- Pending Feedback Section -->
      <section
        class="pending-feedback-section"
        v-if="pendingFeedbacks.length > 0"
      >
        <h3>Pending Feedback</h3>
        <p class="section-desc">
          Please rate these recently completed sessions:
        </p>

        <div class="feedback-list">
          <div
            v-for="(session, index) in pendingFeedbacks"
            :key="index"
            class="feedback-card"
          >
            <div class="session-info">
              <div class="user-avatar">
                <img
                  :src="session.partnerAvatar || '/default-avatar.svg'"
                  :alt="session.partnerName"
                />
              </div>
              <div class="session-details">
                <h4>{{ session.skillName }} with {{ session.partnerName }}</h4>
                <p class="session-date">
                  {{ formatDate(session.date) }} ·
                  {{ session.duration }} minutes
                </p>
                <p class="session-meta">
                  {{
                    session.type === "taught"
                      ? "You taught this session"
                      : "You learned in this session"
                  }}
                </p>
              </div>
            </div>

            <div class="feedback-form">
              <RatingComponent
                :title="
                  session.type === 'taught'
                    ? 'Rate your student'
                    : 'Rate your teacher'
                "
                @submit="submitFeedback(session.id, $event)"
                @cancel="skipFeedback(session.id)"
              />
            </div>
          </div>
        </div>
      </section>

      <!-- Feedback Received -->
      <section class="received-feedback-section">
        <h3>Feedback Received</h3>
        <div v-if="receivedFeedbacks.length === 0" class="empty-feedback">
          <font-awesome-icon icon="comment-alt" class="empty-icon" />
          <p>You haven't received any feedback yet.</p>
          <p class="empty-hint">
            Complete skills exchanges to receive feedback from other users.
          </p>
        </div>

        <div v-else class="reviews-grid">
          <div
            v-for="(review, index) in receivedFeedbacks"
            :key="index"
            class="review-card"
          >
            <div class="review-header">
              <div class="reviewer-info">
                <img
                  :src="review.reviewerAvatar || '/default-avatar.svg'"
                  :alt="review.reviewerName"
                  class="reviewer-avatar"
                />
                <div>
                  <h4 class="reviewer-name">{{ review.reviewerName }}</h4>
                  <p class="review-date">{{ formatDate(review.date) }}</p>
                </div>
              </div>
              <div class="review-rating">
                <RatingComponent
                  :value="review.rating"
                  :viewOnly="true"
                  :showCount="false"
                />
              </div>
            </div>
            <div class="review-content">
              <p>{{ review.feedback }}</p>
            </div>
            <div class="review-meta">
              <span class="review-skill">
                <font-awesome-icon icon="graduation-cap" />
                {{ review.skillName }}
              </span>
              <span class="review-type">
                <font-awesome-icon icon="exchange-alt" />
                {{ review.type === "teacher" ? "You taught" : "You learned" }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- Feedback Given -->
      <section class="given-feedback-section">
        <h3>Feedback Given</h3>
        <div v-if="givenFeedbacks.length === 0" class="empty-feedback">
          <font-awesome-icon icon="comment" class="empty-icon" />
          <p>You haven't given any feedback yet.</p>
        </div>

        <div v-else class="reviews-grid">
          <div
            v-for="(review, index) in givenFeedbacks"
            :key="index"
            class="review-card"
          >
            <div class="review-header">
              <div class="reviewer-info">
                <img
                  :src="review.recipientAvatar || '/default-avatar.svg'"
                  :alt="review.recipientName"
                  class="reviewer-avatar"
                />
                <div>
                  <h4 class="reviewer-name">{{ review.recipientName }}</h4>
                  <p class="review-date">{{ formatDate(review.date) }}</p>
                </div>
              </div>
              <div class="review-rating">
                <RatingComponent
                  :value="review.rating"
                  :viewOnly="true"
                  :showCount="false"
                />
              </div>
            </div>
            <div class="review-content">
              <p>{{ review.feedback }}</p>
            </div>
            <div class="review-meta">
              <span class="review-skill">
                <font-awesome-icon icon="graduation-cap" />
                {{ review.skillName }}
              </span>
              <span class="review-type">
                <font-awesome-icon icon="exchange-alt" />
                {{ review.type === "teacher" ? "They taught" : "They learned" }}
              </span>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import RatingComponent from "./RatingComponent.vue";

export default {
  name: "FeedbackSystem",
  components: {
    RatingComponent,
  },
  data() {
    return {
      pendingFeedbacks: [
        {
          id: 1,
          partnerName: "Alice Smith",
          partnerAvatar: null,
          skillName: "JavaScript Fundamentals",
          date: new Date(Date.now() - 86400000), // 1 day ago
          duration: 60,
          type: "taught", // 'taught' or 'learned'
        },
        {
          id: 2,
          partnerName: "Bob Johnson",
          partnerAvatar: null,
          skillName: "Italian Cooking",
          date: new Date(Date.now() - 172800000), // 2 days ago
          duration: 90,
          type: "learned",
        },
      ],
      receivedFeedbacks: [
        {
          id: 1,
          reviewerName: "Carol Williams",
          reviewerAvatar: null,
          rating: 5,
          feedback:
            "Excellent teacher! Explained complex concepts in a very accessible way. I learned a lot and would definitely recommend!",
          date: new Date(Date.now() - 604800000), // 1 week ago
          skillName: "Python Programming",
          type: "teacher",
        },
        {
          id: 2,
          reviewerName: "David Brown",
          reviewerAvatar: null,
          rating: 4.5,
          feedback:
            "Very knowledgeable and patient. The session was well structured and I came away with a much better understanding.",
          date: new Date(Date.now() - 1209600000), // 2 weeks ago
          skillName: "Web Development",
          type: "teacher",
        },
      ],
      givenFeedbacks: [
        {
          id: 1,
          recipientName: "Eve Jones",
          recipientAvatar: null,
          rating: 5,
          feedback:
            "Eve is a fantastic teacher! She made learning Spanish fun and engaging. Would definitely recommend!",
          date: new Date(Date.now() - 518400000), // 6 days ago
          skillName: "Spanish Language",
          type: "teacher",
        },
      ],
    };
  },
  methods: {
    formatDate(date) {
      return new Date(date).toLocaleDateString(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    },
    submitFeedback(sessionId, feedback) {
      console.log(`Submitting feedback for session ${sessionId}:`, feedback);

      // In a real app, you would send this to your API
      alert(`Feedback submitted: ${feedback.rating} stars`);

      // Remove from pending list
      this.pendingFeedbacks = this.pendingFeedbacks.filter(
        (session) => session.id !== sessionId,
      );

      // Add to given feedbacks
      const session = this.pendingFeedbacks.find((s) => s.id === sessionId);
      if (session) {
        this.givenFeedbacks.unshift({
          id: Date.now(), // Generate a unique ID
          recipientName: session.partnerName,
          recipientAvatar: session.partnerAvatar,
          rating: feedback.rating,
          feedback: feedback.feedback,
          date: new Date(),
          skillName: session.skillName,
          type: session.type === "taught" ? "learned" : "teacher",
        });
      }
    },
    skipFeedback(sessionId) {
      // Just remove from pending list
      this.pendingFeedbacks = this.pendingFeedbacks.filter(
        (session) => session.id !== sessionId,
      );
    },
  },
};
</script>

<style scoped>
.feedback-page {
  padding-bottom: var(--space-12);
}

h2 {
  text-align: center;
  color: var(--primary-color);
  margin-bottom: var(--space-2);
}

.subtitle {
  text-align: center;
  color: var(--medium);
  margin-bottom: var(--space-8);
}

section {
  margin-bottom: var(--space-8);
}

h3 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-4);
  color: var(--dark);
  border-bottom: 1px solid var(--light);
  padding-bottom: var(--space-2);
}

.section-desc {
  margin-bottom: var(--space-4);
  color: var(--medium);
}

.feedback-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.feedback-card {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-4);
}

.session-info {
  display: flex;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.session-details h4 {
  margin: 0 0 var(--space-1) 0;
  font-size: var(--font-size-md);
}

.session-date {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--medium);
}

.session-meta {
  margin: var(--space-1) 0 0 0;
  font-size: var(--font-size-sm);
  color: var(--primary-color);
}

.empty-feedback {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-8);
  background-color: var(--white);
  border-radius: var(--radius-lg);
  text-align: center;
}

.empty-icon {
  font-size: var(--font-size-3xl);
  color: var(--medium);
  opacity: 0.6;
  margin-bottom: var(--space-4);
}

.empty-hint {
  font-size: var(--font-size-sm);
  color: var(--medium);
  margin-top: var(--space-2);
}

.reviews-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-4);
}

.review-card {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
}

.review-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.reviewer-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.reviewer-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.reviewer-name {
  margin: 0;
  font-size: var(--font-size-md);
}

.review-date {
  margin: var(--space-1) 0 0 0;
  font-size: var(--font-size-xs);
  color: var(--medium);
}

.review-content {
  margin-bottom: var(--space-3);
  flex-grow: 1;
}

.review-content p {
  margin: 0;
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.review-meta {
  display: flex;
  gap: var(--space-3);
  font-size: var(--font-size-xs);
  color: var(--medium);
}

.review-skill,
.review-type {
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

@media (max-width: 768px) {
  .reviews-grid {
    grid-template-columns: 1fr;
  }
}
</style>
