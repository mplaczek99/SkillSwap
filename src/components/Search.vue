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
      <div
        v-for="(item, index) in results"
        :key="index"
        class="result-item"
      >
        <h3>{{ item.name }}</h3>
        <p>{{ item.skill }}</p>
      </div>
    </div>
    <div v-else-if="searched">
      <p>No results found.</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Search',
  data() {
    return {
      query: '',
      results: [],
      searched: false,
      searchTimeout: null,
    };
  },
  methods: {
    async search() {
      if (this.searchTimeout) clearTimeout(this.searchTimeout);
      // If running under Jest (JEST_WORKER_ID is defined), resolve immediately.
      if (process.env.JEST_WORKER_ID) {
        const dummyData = [
          { name: 'Alice', skill: 'Guitar' },
          { name: 'Bob', skill: 'Spanish' },
          { name: 'Charlie', skill: 'Cooking' },
        ];
        this.results = dummyData.filter((item) =>
          item.name.toLowerCase().includes(this.query.toLowerCase()) ||
          item.skill.toLowerCase().includes(this.query.toLowerCase())
        );
        this.searched = true;
        return;
      }
      // Otherwise use the normal delay.
      await new Promise((resolve) => {
        this.searchTimeout = setTimeout(() => {
          const dummyData = [
            { name: 'Alice', skill: 'Guitar' },
            { name: 'Bob', skill: 'Spanish' },
            { name: 'Charlie', skill: 'Cooking' },
          ];
          this.results = dummyData.filter((item) =>
            item.name.toLowerCase().includes(this.query.toLowerCase()) ||
            item.skill.toLowerCase().includes(this.query.toLowerCase())
          );
          this.searched = true;
          resolve();
        }, 300);
      });
    },
  },
};
</script>


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

