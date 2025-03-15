<template>
  <div class="job-postings-page">
    <div class="container">
      <!-- Hero Section -->
      <section class="jobs-hero">
        <h1>Find Your Next Opportunity</h1>
        <p class="subtitle">Browse job postings from skills exchange users</p>

        <div class="search-container">
          <div class="search-input-group">
            <font-awesome-icon icon="search" class="search-icon" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search job titles, skills, or companies..."
              class="search-input"
              @input="filterJobs"
            />
            <button
              v-if="searchQuery"
              type="button"
              class="clear-search"
              @click="clearSearch"
            >
              <font-awesome-icon icon="times" />
            </button>
          </div>

          <router-link to="/post-job" class="btn btn-primary post-job-btn">
            <font-awesome-icon icon="plus" /> Post a Job
          </router-link>
        </div>

        <div class="filter-bar">
          <div class="filter-group">
            <label>Job Type</label>
            <select v-model="filters.jobType" @change="filterJobs">
              <option value="">All Types</option>
              <option value="Full-time">Full-time</option>
              <option value="Part-time">Part-time</option>
              <option value="Contract">Contract</option>
              <option value="Freelance">Freelance</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Experience</label>
            <select v-model="filters.experienceLevel" @change="filterJobs">
              <option value="">All Levels</option>
              <option value="Entry">Entry Level</option>
              <option value="Mid">Mid Level</option>
              <option value="Senior">Senior Level</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Location</label>
            <select v-model="filters.location" @change="filterJobs">
              <option value="">All Locations</option>
              <option
                v-for="location in uniqueLocations"
                :key="location"
                :value="location"
              >
                {{ location }}
              </option>
            </select>
          </div>

          <button
            class="btn btn-outline btn-sm filter-reset"
            @click="resetFilters"
          >
            <font-awesome-icon icon="undo" /> Reset Filters
          </button>
        </div>
      </section>

      <!-- Jobs Listing Section -->
      <section class="jobs-listing">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Loading job opportunities...</p>
        </div>

        <div v-else-if="error" class="error-message">
          <font-awesome-icon icon="exclamation-circle" />
          <p>{{ error }}</p>
          <button @click="fetchJobs" class="btn btn-outline btn-sm">
            Try Again
          </button>
        </div>

        <div v-else-if="filteredJobs.length === 0" class="no-jobs">
          <font-awesome-icon icon="briefcase" class="no-jobs-icon" />
          <h3>No Job Postings Found</h3>
          <p v-if="searchQuery || hasActiveFilters">
            We couldn't find any jobs matching your search criteria.
          </p>
          <p v-else>There are no job postings available at the moment.</p>
          <div class="no-jobs-actions">
            <button
              @click="resetFilters"
              class="btn btn-outline btn-sm"
              v-if="hasActiveFilters"
            >
              Clear Filters
            </button>
            <router-link to="/post-job" class="btn btn-primary btn-sm">
              Post a Job
            </router-link>
          </div>
        </div>

        <div v-else class="jobs-grid">
          <div v-for="job in paginatedJobs" :key="job.id" class="job-card">
            <div class="job-card-header">
              <h3 class="job-title">{{ job.title }}</h3>
              <span class="job-company">{{ job.company }}</span>
              <div class="job-meta">
                <span class="job-location">
                  <font-awesome-icon icon="map-marker-alt" />
                  {{ job.location }}
                </span>
                <span class="job-type" :class="getJobTypeClass(job.jobType)">
                  {{ job.jobType }}
                </span>
              </div>
            </div>

            <div class="job-card-body">
              <p class="job-description">
                {{ truncateDescription(job.description) }}
              </p>

              <div class="job-skills">
                <span
                  v-for="(skill, index) in job.skillsArray()"
                  :key="index"
                  class="skill-tag"
                >
                  {{ skill }}
                </span>
              </div>
            </div>

            <div class="job-card-footer">
              <div class="job-info">
                <span class="job-date"
                  >Posted {{ job.daysSincePosting() }} days ago</span
                >
                <span class="job-level">{{ job.experienceLevel }} Level</span>
              </div>

              <router-link
                :to="`/jobs/${job.id}`"
                class="btn btn-primary btn-sm view-job-btn"
              >
                View Details
              </router-link>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="filteredJobs.length > 0" class="pagination">
          <button
            class="pagination-btn"
            :disabled="currentPage === 1"
            @click="changePage(currentPage - 1)"
          >
            <font-awesome-icon icon="chevron-left" />
          </button>

          <span class="pagination-info">
            Page {{ currentPage }} of {{ totalPages }}
          </span>

          <button
            class="pagination-btn"
            :disabled="currentPage === totalPages"
            @click="changePage(currentPage + 1)"
          >
            <font-awesome-icon icon="chevron-right" />
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import JobPost from "@/models/JobPost";

export default {
  name: "JobPostings",
  data() {
    return {
      jobs: [],
      filteredJobs: [],
      loading: true,
      error: null,
      searchQuery: "",
      filters: {
        jobType: "",
        experienceLevel: "",
        location: "",
      },
      currentPage: 1,
      jobsPerPage: 9,
    };
  },
  computed: {
    uniqueLocations() {
      const locations = this.jobs.map((job) => job.location);
      return [...new Set(locations)].filter(Boolean);
    },
    hasActiveFilters() {
      return (
        Object.values(this.filters).some((value) => value !== "") ||
        this.searchQuery !== ""
      );
    },
    totalPages() {
      return Math.ceil(this.filteredJobs.length / this.jobsPerPage);
    },
    paginatedJobs() {
      const start = (this.currentPage - 1) * this.jobsPerPage;
      const end = start + this.jobsPerPage;
      return this.filteredJobs.slice(start, end);
    },
  },
  created() {
    this.fetchJobs();

    // Check if there's a query parameter in the URL
    if (this.$route.query.q) {
      this.searchQuery = this.$route.query.q;
    }

    // Check other filter params
    if (this.$route.query.type) {
      this.filters.jobType = this.$route.query.type;
    }
    if (this.$route.query.level) {
      this.filters.experienceLevel = this.$route.query.level;
    }
    if (this.$route.query.location) {
      this.filters.location = this.$route.query.location;
    }
  },
  methods: {
    async fetchJobs() {
      this.loading = true;
      this.error = null;

      try {
        // In a real implementation, this would call an API
        // For now, we'll use mock data
        setTimeout(() => {
          this.jobs = this.getMockJobs().map((job) => new JobPost(job));
          this.filterJobs();
          this.loading = false;
        }, 1000);

        // Real implementation would be:
        // const response = await axios.get('/api/jobs');
        // this.jobs = response.data.map(job => new JobPost(job));
      } catch (error) {
        console.error("Error fetching jobs:", error);
        this.error = "Failed to load job postings. Please try again.";
        this.loading = false;
      }
    },

    filterJobs() {
      // Update the URL with filter parameters
      this.$router.replace({
        query: {
          ...this.$route.query,
          q: this.searchQuery || undefined,
          type: this.filters.jobType || undefined,
          level: this.filters.experienceLevel || undefined,
          location: this.filters.location || undefined,
        },
      });

      let result = [...this.jobs];

      // Apply search filter
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        result = result.filter(
          (job) =>
            job.title.toLowerCase().includes(query) ||
            job.company.toLowerCase().includes(query) ||
            job.description.toLowerCase().includes(query) ||
            job
              .skillsArray()
              .some((skill) => skill.toLowerCase().includes(query)),
        );
      }

      // Apply dropdown filters
      if (this.filters.jobType) {
        result = result.filter((job) => job.jobType === this.filters.jobType);
      }

      if (this.filters.experienceLevel) {
        result = result.filter(
          (job) => job.experienceLevel === this.filters.experienceLevel,
        );
      }

      if (this.filters.location) {
        result = result.filter((job) => job.location === this.filters.location);
      }

      this.filteredJobs = result;
      this.currentPage = 1; // Reset to first page when filtering
    },

    resetFilters() {
      this.searchQuery = "";
      this.filters = {
        jobType: "",
        experienceLevel: "",
        location: "",
      };
      this.filterJobs();
    },

    clearSearch() {
      this.searchQuery = "";
      this.filterJobs();
    },

    getJobTypeClass(jobType) {
      const classes = {
        "Full-time": "full-time",
        "Part-time": "part-time",
        Contract: "contract",
        Freelance: "freelance",
      };
      return classes[jobType] || "";
    },

    truncateDescription(description, maxLength = 150) {
      if (description.length <= maxLength) return description;
      return description.substring(0, maxLength) + "...";
    },

    changePage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
      }
    },

    getMockJobs() {
      return [
        {
          id: 1,
          title: "Frontend Developer",
          company: "Tech Innovators",
          location: "San Francisco, CA",
          description:
            "We are looking for a skilled Frontend Developer to join our team. You will be responsible for building web applications using Vue.js and modern web technologies.",
          skillsRequired: ["Vue.js", "JavaScript", "CSS", "HTML"],
          experienceLevel: "Mid",
          jobType: "Full-time",
          salaryRange: "$80,000 - $110,000",
          contactEmail: "jobs@techinnovators.com",
          postedByUserID: 1,
          postedByName: "Alice Smith",
          createdAt: new Date(Date.now() - 5 * 86400000), // 5 days ago
        },
        {
          id: 2,
          title: "UX/UI Designer",
          company: "Creative Solutions",
          location: "Remote",
          description:
            "Join our design team and help create beautiful and functional user interfaces for our clients. You should have a strong portfolio and experience with design tools.",
          skillsRequired: ["Figma", "Adobe XD", "Prototyping", "User Research"],
          experienceLevel: "Senior",
          jobType: "Full-time",
          salaryRange: "$90,000 - $120,000",
          contactEmail: "careers@creativesolutions.com",
          postedByUserID: 2,
          postedByName: "Bob Johnson",
          createdAt: new Date(Date.now() - 3 * 86400000), // 3 days ago
        },
        {
          id: 3,
          title: "Content Writer",
          company: "Media Pulse",
          location: "New York, NY",
          description:
            "We need a talented Content Writer to create engaging content for our blog and social media channels. Must have excellent writing skills and SEO knowledge.",
          skillsRequired: ["Copywriting", "SEO", "Content Strategy", "Editing"],
          experienceLevel: "Entry",
          jobType: "Part-time",
          salaryRange: "$25 - $35 per hour",
          contactEmail: "hiring@mediapulse.com",
          postedByUserID: 3,
          postedByName: "Carol Williams",
          createdAt: new Date(Date.now() - 7 * 86400000), // 7 days ago
        },
        {
          id: 4,
          title: "Backend Developer",
          company: "Data Systems Inc.",
          location: "Boston, MA",
          description:
            "Looking for a Backend Developer with Go expertise to help build our next-generation API services. Must have experience with database design and RESTful APIs.",
          skillsRequired: ["Go", "SQL", "Docker", "RESTful APIs"],
          experienceLevel: "Senior",
          jobType: "Full-time",
          salaryRange: "$110,000 - $140,000",
          contactEmail: "tech-hiring@datasystems.com",
          postedByUserID: 1,
          postedByName: "Alice Smith",
          createdAt: new Date(Date.now() - 2 * 86400000), // 2 days ago
        },
        {
          id: 5,
          title: "Marketing Specialist",
          company: "Growth Hackers",
          location: "Chicago, IL",
          description:
            "Join our marketing team to develop and implement marketing strategies. You should have experience with digital marketing and analytics tools.",
          skillsRequired: [
            "Social Media Marketing",
            "Google Analytics",
            "SEO",
            "Content Creation",
          ],
          experienceLevel: "Mid",
          jobType: "Full-time",
          salaryRange: "$65,000 - $85,000",
          contactEmail: "jobs@growthhackers.com",
          postedByUserID: 4,
          postedByName: "David Brown",
          createdAt: new Date(Date.now() - 10 * 86400000), // 10 days ago
        },
        {
          id: 6,
          title: "Mobile App Developer",
          company: "App Wizards",
          location: "Seattle, WA",
          description:
            "We need a skilled mobile developer who can build native iOS applications. Knowledge of Swift and the Apple ecosystem is required.",
          skillsRequired: ["Swift", "iOS", "Xcode", "Mobile Design"],
          experienceLevel: "Mid",
          jobType: "Contract",
          salaryRange: "$70 - $90 per hour",
          contactEmail: "devjobs@appwizards.com",
          postedByUserID: 2,
          postedByName: "Bob Johnson",
          createdAt: new Date(Date.now() - 4 * 86400000), // 4 days ago
        },
        {
          id: 7,
          title: "Data Scientist",
          company: "Analytics Pro",
          location: "Remote",
          description:
            "Looking for a Data Scientist to join our team. You will analyze large datasets and build machine learning models to solve business problems.",
          skillsRequired: [
            "Python",
            "Machine Learning",
            "SQL",
            "Data Visualization",
          ],
          experienceLevel: "Senior",
          jobType: "Full-time",
          salaryRange: "$120,000 - $150,000",
          contactEmail: "talent@analyticspro.com",
          postedByUserID: 5,
          postedByName: "Eve Jones",
          createdAt: new Date(Date.now() - 1 * 86400000), // 1 day ago
        },
        {
          id: 8,
          title: "Product Manager",
          company: "Innovation Labs",
          location: "Austin, TX",
          description:
            "We are looking for a Product Manager to lead product development and work with cross-functional teams to deliver great user experiences.",
          skillsRequired: [
            "Product Strategy",
            "User Stories",
            "Agile",
            "Market Research",
          ],
          experienceLevel: "Senior",
          jobType: "Full-time",
          salaryRange: "$100,000 - $130,000",
          contactEmail: "pm-hiring@innovationlabs.com",
          postedByUserID: 3,
          postedByName: "Carol Williams",
          createdAt: new Date(Date.now() - 6 * 86400000), // 6 days ago
        },
        {
          id: 9,
          title: "DevOps Engineer",
          company: "Cloud Systems",
          location: "Denver, CO",
          description:
            "Join our team to build and maintain CI/CD pipelines and cloud infrastructure. Experience with AWS and containerization is required.",
          skillsRequired: ["AWS", "Docker", "Kubernetes", "CI/CD"],
          experienceLevel: "Mid",
          jobType: "Full-time",
          salaryRange: "$90,000 - $120,000",
          contactEmail: "careers@cloudsystems.com",
          postedByUserID: 4,
          postedByName: "David Brown",
          createdAt: new Date(Date.now() - 8 * 86400000), // 8 days ago
        },
        {
          id: 10,
          title: "Graphic Designer",
          company: "Creative Works",
          location: "Portland, OR",
          description:
            "We need a creative Graphic Designer to join our team. You will create visual concepts for web and print materials.",
          skillsRequired: [
            "Adobe Creative Suite",
            "Typography",
            "Branding",
            "Illustration",
          ],
          experienceLevel: "Entry",
          jobType: "Part-time",
          salaryRange: "$20 - $30 per hour",
          contactEmail: "design@creativeworks.com",
          postedByUserID: 5,
          postedByName: "Eve Jones",
          createdAt: new Date(Date.now() - 12 * 86400000), // 12 days ago
        },
      ];
    },
  },
};
</script>

<style scoped>
.job-postings-page {
  padding-bottom: var(--space-12);
}

/* Hero section */
.jobs-hero {
  text-align: center;
  margin-bottom: var(--space-8);
}

.jobs-hero h1 {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-2);
  color: var(--dark);
}

.subtitle {
  font-size: var(--font-size-lg);
  color: var(--medium);
  margin-bottom: var(--space-6);
}

.search-container {
  margin-bottom: var(--space-6);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--space-4);
}

.search-input-group {
  position: relative;
  flex: 1;
  max-width: 600px;
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
  padding: var(--space-3) var(--space-4) var(--space-3) var(--space-10);
  border: 2px solid var(--light);
  border-radius: var(--radius-full);
  font-size: var(--font-size-md);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-fast);
}

.search-input:focus {
  border-color: var(--primary-color);
  box-shadow:
    0 0 0 3px var(--primary-light),
    var(--shadow-md);
  outline: none;
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
}

.clear-search:hover {
  color: var(--dark);
}

.post-job-btn {
  white-space: nowrap;
}

/* Filter bar */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: var(--space-4);
  margin-bottom: var(--space-8);
  padding: var(--space-4);
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.filter-group {
  display: flex;
  flex-direction: column;
  min-width: 150px;
}

.filter-group label {
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-1);
  color: var(--dark);
  font-weight: var(--font-weight-medium);
}

.filter-group select {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  background-color: var(--white);
  font-size: var(--font-size-sm);
  color: var(--dark);
  cursor: pointer;
}

.filter-reset {
  align-self: flex-end;
  margin-top: var(--space-4);
}

/* Jobs Grid */
.jobs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--space-6);
  margin-bottom: var(--space-8);
}

.job-card {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  transition: all var(--transition-normal);
  display: flex;
  flex-direction: column;
  border-top: 4px solid var(--primary-color);
  height: 100%;
}

.job-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-lg);
}

.job-card-header {
  padding: var(--space-4);
  border-bottom: 1px solid var(--light);
}

.job-title {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-1);
  color: var(--dark);
}

.job-company {
  display: block;
  font-size: var(--font-size-md);
  color: var(--primary-color);
  margin-bottom: var(--space-2);
  font-weight: var(--font-weight-medium);
}

.job-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-sm);
}

.job-location {
  color: var(--medium);
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

.job-type {
  display: inline-block;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
}

.job-type.full-time {
  background-color: #e0f2fe;
  color: #0284c7;
}

.job-type.part-time {
  background-color: #fef3c7;
  color: #d97706;
}

.job-type.contract {
  background-color: #f3e8ff;
  color: #7c3aed;
}

.job-type.freelance {
  background-color: #dcfce7;
  color: #16a34a;
}

.job-card-body {
  padding: var(--space-4);
  flex: 1;
}

.job-description {
  color: var(--medium);
  margin-bottom: var(--space-3);
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.job-skills {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.skill-tag {
  background-color: var(--light);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  color: var(--dark);
}

.job-card-footer {
  padding: var(--space-4);
  border-top: 1px solid var(--light);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.job-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.job-date,
.job-level {
  font-size: var(--font-size-xs);
  color: var(--medium);
}

.view-job-btn {
  white-space: nowrap;
}

/* Loading and Error States */
.loading-state,
.error-message,
.no-jobs {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12);
  text-align: center;
}

.spinner {
  border: 4px solid var(--light);
  border-top: 4px solid var(--primary-color);
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: var(--space-4);
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.error-message {
  color: var(--error-color);
}

.error-message svg,
.no-jobs-icon {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
  opacity: 0.6;
}

.no-jobs-actions {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-4);
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--space-3);
  margin-top: var(--space-8);
}

.pagination-btn {
  background-color: var(--white);
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.pagination-btn:hover:not(:disabled) {
  background-color: var(--primary-light);
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-info {
  color: var(--medium);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .search-container {
    flex-direction: column;
    gap: var(--space-3);
  }

  .post-job-btn {
    width: 100%;
  }

  .filter-bar {
    flex-direction: column;
    gap: var(--space-3);
  }

  .filter-group {
    width: 100%;
  }

  .filter-reset {
    align-self: center;
  }

  .jobs-grid {
    grid-template-columns: 1fr;
  }
}
</style>
