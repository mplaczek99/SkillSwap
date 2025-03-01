<template>
  <div class="skill-image-container">
    <img :src="iconUrl" :alt="skill.name + ' Icon'" class="skill-image" />
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

@media (max-width: 576px) {
  .skill-image {
    width: 48px;
    height: 48px;
  }
}
</style>
