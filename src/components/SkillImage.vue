<template>
  <img :src="iconUrl" alt="Skill Icon" class="skill-image" />
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
    };
  },
  computed: {
    iconUrl() {
      if (this.skill.image && this.skill.image.trim() !== "") {
        return this.skill.image;
      }
      if (this.dynamicIcon) {
        return `https://api.iconify.design/fa-solid/${this.dynamicIcon}.svg`;
      }
      return "https://api.iconify.design/fa-solid/cog.svg";
    },
  },
  async created() {
    if (!this.skill.image || this.skill.image.trim() === "") {
      this.dynamicIcon = await fetchDynamicIcon(this.skill.name);
    }
  },
};
</script>

<style scoped>
.skill-image {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
}
</style>
