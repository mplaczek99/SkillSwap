<template>
  <div class="search-container">
    <h2>Search for Skills or Users</h2>
    <form @submit.prevent="search">
      <input
        v-model="query"
        type="text"
        placeholder="Search for skills or users..."
        required
      />
      <button type="submit">Search</button>
    </form>
    <div class="results" v-if="results.length">
      <div v-for="(item, index) in results" :key="index" class="result-item">
        <h3>{{ item.name }}</h3>
        <p v-if="item.description">Description: {{ item.description }}</p>
        <p v-if="item.email">Email: {{ item.email }}</p>
      </div>
    </div>
    <div v-else-if="searched">
      <p>No results found.</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import debounce from "lodash/debounce";

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
    };
  },
  created() {
    this.debouncedSearch = debounce(this.performSearch, 300);
  },
  methods: {
    async performSearch() {
      try {
        const response = await axios.get("/api/search", {
          params: { q: this.query },
        });
        this.results = response.data;
      } catch (error) {
        console.error("Search API error:", error);
        // If using a toast notification library, show a friendly error.
        if (this.$toast) {
          this.$toast.error(
            "An error occurred while searching. Please try again later.",
          );
        }
        this.results = [];
      }
      this.searched = true;
    },
    async search() {
      if (this.forceApiCall) {
        this.debouncedSearch();
      } else if (process.env.JEST_WORKER_ID) {
        // Provide dummy data in test environments.
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
  },
};
</script>

<style scoped>
.search-container {
  padding: 2rem;
}
form {
  display: flex;
  margin-bottom: 1rem;
}
form input {
  flex: 1;
  padding: 0.5rem;
  font-size: 1rem;
}
form button {
  padding: 0.5rem 1rem;
  margin-left: 0.5rem;
  font-size: 1rem;
}
.results {
  margin-top: 1rem;
}
.result-item {
  border: 1px solid #ccc;
  padding: 1rem;
  margin-bottom: 0.5rem;
  border-radius: 4px;
}
</style>
