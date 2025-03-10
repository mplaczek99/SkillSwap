<template>
  <div class="rating-component">
    <!-- View Mode -->
    <div v-if="viewOnly" class="rating-display">
      <div class="stars-display">
        <font-awesome-icon
          v-for="star in 5"
          :key="star"
          :icon="getStar(star)"
          :class="getStarClass(star)"
        />
      </div>
      <span class="rating-count" v-if="showCount">
        ({{ ratingCount || 0 }})
      </span>
    </div>

    <!-- Interactive Rating Mode -->
    <div v-else class="rating-input">
      <div class="rating-title">{{ title }}</div>
      <div class="stars-input">
        <font-awesome-icon
          v-for="star in 5"
          :key="star"
          :icon="getInteractiveStar(star)"
          :class="getInteractiveStarClass(star)"
          @mouseover="hoverRating = star"
          @mouseleave="hoverRating = 0"
          @click="setRating(star)"
        />
      </div>
      <div class="rating-feedback" v-if="feedbackEnabled">
        <textarea
          v-model="feedback"
          :placeholder="feedbackPlaceholder"
          maxlength="500"
          rows="3"
        ></textarea>
        <div class="character-count">{{ feedback.length }}/500</div>
      </div>
      <div class="rating-actions" v-if="!viewOnly">
        <button
          type="button"
          class="btn btn-outline btn-sm"
          @click="$emit('cancel')"
        >
          Cancel
        </button>
        <button
          type="button"
          class="btn btn-primary btn-sm"
          @click="submitRating"
          :disabled="
            !currentRating ||
            (feedbackEnabled && feedbackRequired && !feedback.trim())
          "
        >
          Submit
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "RatingComponent",
  props: {
    value: {
      type: Number,
      default: 0,
    },
    count: {
      type: Number,
      default: 0,
    },
    viewOnly: {
      type: Boolean,
      default: false,
    },
    showCount: {
      type: Boolean,
      default: true,
    },
    title: {
      type: String,
      default: "Rate this user",
    },
    feedbackEnabled: {
      type: Boolean,
      default: true,
    },
    feedbackRequired: {
      type: Boolean,
      default: false,
    },
    feedbackPlaceholder: {
      type: String,
      default: "Share your experience with this user...",
    },
  },
  data() {
    return {
      currentRating: this.value,
      hoverRating: 0,
      feedback: "",
      ratingCount: this.count,
    };
  },
  methods: {
    getStar(position) {
      if (this.value >= position) {
        return "star";
      } else if (this.value >= position - 0.5) {
        return "star-half-alt";
      } else {
        return "star";
      }
    },

    getStarClass(position) {
      if (this.value >= position) {
        return "filled";
      } else if (this.value >= position - 0.5) {
        return "half-filled";
      } else {
        return "empty";
      }
    },

    // Removed parameter since it isn't used.
    getInteractiveStar() {
      return "star";
    },

    getInteractiveStarClass(position) {
      const effectiveRating = this.hoverRating || this.currentRating;
      return {
        filled: position <= effectiveRating,
        empty: position > effectiveRating,
      };
    },

    setRating(rating) {
      this.currentRating = rating;
      this.$emit("input", rating);
    },

    submitRating() {
      if (!this.currentRating) return;
      const ratingData = {
        rating: this.currentRating,
        feedback: this.feedback,
      };
      this.$emit("submit", ratingData);
    },
  },
};
</script>

<style scoped>
.rating-component {
  margin-bottom: var(--space-4);
}

.rating-display {
  display: flex;
  align-items: center;
}

.stars-display {
  display: flex;
  gap: 2px;
}

.stars-display .filled {
  color: #ffd700;
  /* Gold color for filled stars */
}

.stars-display .half-filled {
  color: #ffd700;
  /* Gold color for half-filled stars */
}

.stars-display .empty {
  color: #e0e0e0;
  /* Light gray for empty stars */
}

.rating-count {
  margin-left: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--medium);
}

.rating-input {
  background-color: var(--light);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
}

.rating-title {
  margin-bottom: var(--space-3);
  font-weight: var(--font-weight-medium);
  color: var(--dark);
}

.stars-input {
  display: flex;
  gap: var(--space-2);
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-3);
}

.stars-input .filled {
  color: #ffd700;
  cursor: pointer;
}

.stars-input .empty {
  color: #e0e0e0;
  cursor: pointer;
}

.rating-feedback {
  margin-bottom: var(--space-3);
}

.rating-feedback textarea {
  width: 100%;
  padding: var(--space-2);
  border: 1px solid #ddd;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  resize: vertical;
}

.character-count {
  text-align: right;
  font-size: var(--font-size-xs);
  color: var(--medium);
  margin-top: 4px;
}

.rating-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-2);
}
</style>
