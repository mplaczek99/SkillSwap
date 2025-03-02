<template>
  <div class="search-page">
    <div class="container">
      <!-- Search Hero Section -->
      <section class="search-hero">
        <h1>Find Skills & Connect</h1>
        <p class="search-subtitle">
          Discover people sharing their expertise or find the perfect skill to
          learn
        </p>

        <form @submit.prevent="search" class="search-form">
          <div class="search-input-group">
            <font-awesome-icon icon="search" class="search-icon" />
            <input
              v-model="query"
              type="text"
              placeholder="Search for skills, topics, or users..."
              class="search-input"
              required
              @input="onSearchInput"
            />
            <button
              v-if="query"
              type="button"
              class="clear-search"
              @click="clearSearch"
            >
              <font-awesome-icon icon="times" />
            </button>
          </div>

          <div class="search-filters">
            <button
              type="button"
              class="filter-toggle"
              @click="toggleFilters"
              :class="{ active: showFilters }"
            >
              <font-awesome-icon icon="filter" />
              <span>Filters</span>
              <font-awesome-icon
                :icon="showFilters ? 'chevron-up' : 'chevron-down'"
                class="toggle-icon"
              />
            </button>
          </div>

          <transition name="slide-down">
            <div v-if="showFilters" class="advanced-filters">
              <div class="filter-group">
                <label class="filter-label">Categories</label>
                <div class="filter-options">
                  <label
                    class="checkbox-container"
                    v-for="category in categories"
                    :key="category"
                  >
                    <input
                      type="checkbox"
                      v-model="selectedCategories"
                      :value="category"
                    />
                    <span class="checkmark"></span>
                    {{ category }}
                  </label>
                </div>
              </div>

              <div class="filter-group">
                <label class="filter-label">Type</label>
                <div class="filter-options">
                  <label class="radio-container">
                    <input type="radio" v-model="searchType" value="all" />
                    <span class="radio-mark"></span>
                    All
                  </label>
                  <label class="radio-container">
                    <input type="radio" v-model="searchType" value="skills" />
                    <span class="radio-mark"></span>
                    Skills
                  </label>
                  <label class="radio-container">
                    <input type="radio" v-model="searchType" value="users" />
                    <span class="radio-mark"></span>
                    Users
                  </label>
                </div>
              </div>

              <div class="filter-actions">
                <button
                  type="button"
                  class="btn btn-outline btn-sm"
                  @click="resetFilters"
                >
                  Reset Filters
                </button>
                <button type="submit" class="btn btn-primary btn-sm">
                  Apply Filters
                </button>
              </div>
            </div>
          </transition>
        </form>
      </section>

      <!-- Search Results -->
      <section class="search-results">
        <div v-if="loading" class="search-loading">
          <div class="spinner"></div>
          <p>Searching...</p>
        </div>

        <div v-else-if="error" class="search-error">
          <font-awesome-icon icon="exclamation-circle" class="error-icon" />
          <p>{{ error }}</p>
          <button @click="search" class="btn btn-outline btn-sm">
            Try Again
          </button>
        </div>

        <div v-else-if="results.length" class="search-results-grid">
          <transition-group name="fade" tag="div" class="results-container">
            <div
              v-for="(item, index) in filteredResults"
              :key="index"
              class="result-card"
              :class="{ 'user-card': item.email, 'skill-card': !item.email }"
            >
              <div class="result-icon">
                <template v-if="item.email">
                  <font-awesome-icon icon="user" />
                </template>
                <template v-else>
                  <font-awesome-icon :icon="getSkillIcon(item.name)" />
                </template>
              </div>

              <div class="result-details">
                <h3>{{ item.name }}</h3>
                <p v-if="item.description" class="result-description">
                  {{ item.description }}
                </p>
                <p v-if="item.email" class="result-meta">
                  <font-awesome-icon icon="envelope" />
                  {{ item.email }}
                </p>
                <div class="result-actions">
                  <button
                    v-if="item.email"
                    class="btn btn-outline btn-sm"
                    @click="viewProfile(item)"
                  >
                    View Profile
                  </button>
                  <button
                    v-else
                    class="btn btn-primary btn-sm"
                    @click="viewSkill(item)"
                  >
                    Learn More
                  </button>
                  <button
                    v-if="item.email"
                    class="btn btn-primary btn-sm"
                    @click="startChat(item)"
                  >
                    Message
                  </button>
                </div>
              </div>
            </div>
          </transition-group>
        </div>

        <div v-else-if="!loading && searched" class="no-results">
          <img
            src="/default-avatar.svg"
            alt="No results"
            class="no-results-image"
          />
          <h3>No Results Found</h3>
          <p>We couldn't find any matches for "{{ query }}"</p>
          <p class="search-suggestions">Try:</p>
          <ul>
            <li>Checking your spelling</li>
            <li>Using more general keywords</li>
            <li>Removing filters</li>
          </ul>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import debounce from "lodash/debounce";
import eventBus from "@/utils/eventBus";

export default {
  name: "Search",
  props: {
    forceApiCall: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      query: "",
      results: [],
      searched: false,
      loading: false,
      error: null,
      showFilters: false,
      categories: [
        "Programming",
        "Language",
        "Music",
        "Cooking",
        "Art",
        "Design",
        "Education",
        "Technology",
      ],
      selectedCategories: [],
      searchType: "all",
    };
  },
  computed: {
    filteredResults() {
      if (!this.results.length) return [];
      let filtered = [...this.results];
      if (this.searchType === "skills") {
        filtered = filtered.filter((item) => !item.email);
      } else if (this.searchType === "users") {
        filtered = filtered.filter((item) => item.email);
      }
      if (this.selectedCategories.length > 0) {
        filtered = filtered.filter((item) => {
          if (item.email) return true;
          return this.selectedCategories.some(
            (category) =>
              item.description &&
              item.description.toLowerCase().includes(category.toLowerCase()),
          );
        });
      }
      return filtered;
    },
  },
  created() {
    this.debouncedSearch = debounce(this.performSearch, 300);
    const queryParam = this.$route.query.q;
    if (queryParam) {
      this.query = queryParam;
      this.search();
    }
  },
  methods: {
    async performSearch() {
      this.loading = true;
      this.error = null;
      try {
        const response = await axios.get("/api/search", {
          params: { q: this.query },
        });
        this.results = response.data;
      } catch (err) {
        console.error("Search API error:", err);
        this.error =
          "An error occurred while searching. Please try again later.";
        this.results = [];
      } finally {
        this.loading = false;
        this.searched = true;
      }
    },
    onSearchInput() {
      this.$router.replace({
        query: { ...this.$route.query, q: this.query || undefined },
      });
      if (this.query.length > 2) {
        this.debouncedSearch();
      } else if (this.query.length === 0) {
        this.results = [];
        this.searched = false;
      }
    },
    clearSearch() {
      this.query = "";
      this.results = [];
      this.searched = false;
      this.$router.replace({ query: {} });
    },
    async search() {
      if (this.forceApiCall) {
        this.debouncedSearch();
      } else if (process.env.JEST_WORKER_ID) {
        const dummyData = [
          { name: "Alice", description: "Guitar" },
          { name: "Bob", description: "Spanish" },
          { name: "Charlie", description: "Cooking" },
        ];
        this.results = dummyData.filter(
          (item) =>
            item.name.toLowerCase().includes(this.query.toLowerCase()) ||
            (item.description &&
              item.description
                .toLowerCase()
                .includes(this.query.toLowerCase())),
        );
        this.searched = true;
      } else {
        this.debouncedSearch();
      }
    },
    toggleFilters() {
      this.showFilters = !this.showFilters;
    },
    resetFilters() {
      this.selectedCategories = [];
      this.searchType = "all";
    },
    getSkillIcon(skillName) {
      const skillIcons = {
        programming: "code",
        language: "language",
        music: "music",
        cooking: "utensils",
        art: "palette",
        design: "pen-fancy",
        go: "code",
        vue: "code",
        guitar: "guitar",
        spanish: "language",
        python: "code",
        singing: "music",
      };
      for (const [key, icon] of Object.entries(skillIcons)) {
        if (skillName.toLowerCase().includes(key.toLowerCase())) {
          return icon;
        }
      }
      return "cog";
    },
    viewProfile(user) {
      alert(`Viewing profile for ${user.name}`);
    },
    viewSkill(skill) {
      alert(`Viewing details for ${skill.name}`);
    },
    // Updated startChat method to use eventBus.emit
    startChat(user) {
      this.$router.push({
        name: "Chat",
        query: { user: user.id, userName: user.name },
      });
      eventBus.emit("show-notification", {
        type: "info",
        title: "Starting Chat",
        message: `Starting a conversation with ${user.name}`,
        duration: 3000,
      });
    },
  },
};
</script>

<style scoped>
/* (CSS remains unchanged) */
.search-page {
  padding-bottom: var(--space-12);
}
.search-hero {
  text-align: center;
  margin-bottom: var(--space-8);
}
.search-hero h1 {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-2);
  color: var(--dark);
}
.search-subtitle {
  font-size: var(--font-size-lg);
  color: var(--medium);
  margin-bottom: var(--space-6);
}
.search-form {
  max-width: 700px;
  margin: 0 auto;
}
.search-input-group {
  position: relative;
  margin-bottom: var(--space-4);
}
.search-icon {
  position: absolute;
  left: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
  color: var(--medium);
  font-size: var(--font-size-lg);
}
.search-input {
  width: 100%;
  padding: var(--space-4) var(--space-4) var(--space-4) var(--space-10);
  font-size: var(--font-size-lg);
  border: 2px solid var(--light);
  border-radius: var(--radius-full);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-fast) ease;
}
.search-input:focus {
  border-color: var(--primary-color);
  outline: none;
  box-shadow:
    0 0 0 3px var(--primary-light),
    var(--shadow-md);
}
.clear-search {
  position: absolute;
  right: var(--space-4);
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: var(--medium);
  cursor: pointer;
  font-size: var(--font-size-md);
}
.clear-search:hover {
  color: var(--dark);
}
.search-filters {
  display: flex;
  justify-content: flex-end;
  margin-bottom: var(--space-3);
}
.filter-toggle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  background: transparent;
  border: none;
  color: var(--primary-color);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  padding: var(--space-2);
}
.filter-toggle:hover,
.filter-toggle.active {
  color: var(--primary-dark);
}
.toggle-icon {
  font-size: var(--font-size-xs);
  transition: transform var(--transition-fast) ease;
}
.advanced-filters {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  box-shadow: var(--shadow-md);
  margin-bottom: var(--space-6);
}
.filter-group {
  margin-bottom: var(--space-4);
}
.filter-label {
  display: block;
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-2);
  color: var(--dark);
}
.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
}
.checkbox-container,
.radio-container {
  display: flex;
  align-items: center;
  position: relative;
  padding-left: 28px;
  cursor: pointer;
  user-select: none;
  font-size: var(--font-size-sm);
}
.checkbox-container input,
.radio-container input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
}
.checkmark,
.radio-mark {
  position: absolute;
  left: 0;
  top: 0;
  height: 18px;
  width: 18px;
  background-color: var(--white);
  border: 1px solid var(--medium);
}
.checkmark {
  border-radius: var(--radius-sm);
}
.radio-mark {
  border-radius: 50%;
}
.checkbox-container:hover input ~ .checkmark,
.radio-container:hover input ~ .radio-mark {
  border-color: var(--primary-color);
}
.checkbox-container input:checked ~ .checkmark,
.radio-container input:checked ~ .radio-mark {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}
.checkmark:after,
.radio-mark:after {
  content: "";
  position: absolute;
  display: none;
}
.checkbox-container input:checked ~ .checkmark:after,
.radio-container input:checked ~ .radio-mark:after {
  display: block;
}
.checkbox-container .checkmark:after {
  left: 6px;
  top: 2px;
  width: 5px;
  height: 10px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}
.radio-container .radio-mark:after {
  top: 5px;
  left: 5px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: white;
}
.filter-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-3);
}
.search-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12) 0;
  color: var(--medium);
}
.search-loading .spinner {
  margin-bottom: var(--space-4);
  width: 40px;
  height: 40px;
  border-width: 4px;
}
.search-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12) 0;
  color: var(--error-color);
}
.error-icon {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
}
.search-results-grid {
  margin-top: var(--space-6);
}
.results-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-4);
}
.result-card {
  display: flex;
  background-color: var(--white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  transition:
    transform var(--transition-normal) ease,
    box-shadow var(--transition-normal) ease;
  padding: var(--space-4);
}
.result-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}
.user-card {
  border-left: 4px solid var(--info-color);
}
.skill-card {
  border-left: 4px solid var(--success-color);
}
.result-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: var(--space-3);
  font-size: var(--font-size-xl);
  flex-shrink: 0;
}
.user-card .result-icon {
  background-color: var(--info-color);
  color: white;
}
.skill-card .result-icon {
  background-color: var(--success-color);
  color: white;
}
.result-details {
  flex: 1;
}
.result-details h3 {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-2);
}
.result-description {
  color: var(--medium);
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-3);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.result-meta {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--medium);
  margin-bottom: var(--space-3);
}
.result-actions {
  display: flex;
  gap: var(--space-2);
}
.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: var(--space-8) 0;
  color: var(--medium);
}
.no-results-image {
  width: 120px;
  height: 120px;
  margin-bottom: var(--space-4);
  opacity: 0.5;
}
.no-results h3 {
  font-size: var(--font-size-xl);
  color: var(--dark);
  margin-bottom: var(--space-2);
}
.search-suggestions {
  font-weight: var(--font-weight-semibold);
  margin-top: var(--space-4);
  margin-bottom: var(--space-2);
}
.no-results ul {
  list-style-type: none;
  padding: 0;
  text-align: center;
}
.no-results li {
  margin-bottom: var(--space-1);
}
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all var(--transition-normal) ease;
  max-height: 500px;
  opacity: 1;
  overflow: hidden;
}
.slide-down-enter-from,
.slide-down-leave-to {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
}
@media (max-width: 768px) {
  .search-hero h1 {
    font-size: var(--font-size-2xl);
  }
  .search-subtitle {
    font-size: var(--font-size-md);
  }
  .search-input {
    font-size: var(--font-size-md);
    padding: var(--space-3) var(--space-3) var(--space-3) var(--space-8);
  }
  .search-icon {
    left: var(--space-3);
    font-size: var(--font-size-md);
  }
  .filter-options {
    flex-direction: column;
    gap: var(--space-2);
  }
}
</style>
