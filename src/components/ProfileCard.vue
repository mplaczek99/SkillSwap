<template>
  <div class="profile-card">
    <img :src="imageSrc" alt="Profile picture" class="profile-card__avatar" />
    <div class="profile-card__info">
      <h2 class="profile-card__title">{{ title }}</h2>
      <p class="profile-card__description">{{ description }}</p>
      <button class="profile-card__button" @click="viewProfile">Learn More</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProfileCard',
  props: {
    title: { type: String, required: true },
    description: { type: String, required: true },
    imageSrc: { type: String, required: true }
  },
  methods: {
    viewProfile() {
      this.$emit('open-profile');
    }
  }
};
</script>

<style scoped>
.profile-card {
  display: flex;
  align-items: center;
  background-color: var(--card-bg, #fff);
  border: 1px solid var(--card-border, #e0e0e0);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.profile-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}
.profile-card__avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 1rem;
}
.profile-card__info {
  flex: 1;
  min-width: 0;
}
.profile-card__title {
  font-size: 1.25rem;
  margin: 0;
  color: var(--text-color, #333);
}
.profile-card__description {
  margin: 0.5rem 0;
  color: var(--text-light, #555);
  line-height: 1.5;
}
.profile-card__button {
  display: inline-block;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  background-color: var(--primary-color, #4CAF50);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}
.profile-card__button:hover {
  background-color: var(--primary-color-dark, #388E3C);
}
@media (max-width: 600px) {
  .profile-card {
    flex-direction: column;
    align-items: flex-start;
  }
  .profile-card__avatar {
    margin: 0 0 0.75rem 0;
  }
}
</style>

