<template>
  <div class="skill-image-container">
    <!-- Loading state -->
    <div v-if="isLoading" class="skill-image loading-icon">
      <font-awesome-icon icon="spinner" spin />
    </div>

    <!-- Error state with fallback icon -->
    <div v-else-if="hasError || imgLoadError" class="skill-image fallback-icon">
      <font-awesome-icon icon="cog" />
    </div>

    <!-- Successful image load -->
    <img
      v-else
      :src="iconUrl"
      :alt="skill.name + ' Icon'"
      class="skill-image"
      @error="handleImageError"
    />
  </div>
</template>

<script>
import { fetchDynamicIcon } from "@/services/iconService";

export default {
  name: "SkillImage",
  props: {
    skill: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      dynamicIcon: null,
      isLoading: true,
      hasError: false,
      imgLoadError: false,
    };
  },
  computed: {
    iconUrl() {
      if (this.skill.image && this.skill.image.trim() !== "") {
        return this.skill.image;
      }
      if (this.dynamicIcon) {
        // Fixed template string syntax
        return `https://api.iconify.design/fa-solid/${this.dynamicIcon}.svg`;
      }
      return "https://api.iconify.design/fa-solid/cog.svg";
    },
  },
  async created() {
    try {
      if (!this.skill.image || this.skill.image.trim() === "") {
        this.dynamicIcon = await fetchDynamicIcon(this.skill.name);
      }
    } catch (error) {
      console.error("Failed to fetch dynamic icon:", error);
      this.hasError = true;
    } finally {
      this.isLoading = false;
    }
  },
  methods: {
    handleImageError() {
      console.warn(`Failed to load image for skill: ${this.skill.name}`);
      this.imgLoadError = true;
    },
  },
};
</script>

<style scoped>
.skill-image-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.skill-image {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  background-color: var(--light);
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition-fast) ease;
}

.skill-image:hover {
  transform: scale(1.05);
}

.loading-icon,
.fallback-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--light);
  color: var(--medium);
  font-size: 1.5rem;
}

@media (max-width: 576px) {
  .skill-image,
  .loading-icon,
  .fallback-icon {
    width: 48px;
    height: 48px;
  }
}
</style>
