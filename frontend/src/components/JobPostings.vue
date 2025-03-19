<template>
  <div class="job-postings-page" :class="{ 'dark-theme': isDark }">
    <div class="container">
      <!-- Hero Section -->
      <section class="jobs-hero">
        <div class="theme-toggle" @click="toggleDark">
          <font-awesome-icon :icon="isDark ? 'sun' : 'moon'" />
        </div>
        
        <h1 class="animated-title">
          <span class="title-word">Find</span>
          <span class="title-word">Your</span>
          <span class="title-word">Next</span>
          <span class="title-word">Opportunity</span>
        </h1>
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
            <font-awesome-icon icon="plus" class="btn-icon" /> Post a Job
          </router-link>
        </div>

        <div class="filter-bar">
          <div class="filter-group">
            <label><font-awesome-icon icon="briefcase" /> Job Type</label>
            <div class="select-wrapper">
              <select v-model="filters.jobType" @change="filterJobs">
                <option value="">All Types</option>
                <option value="Full-time">Full-time</option>
                <option value="Part-time">Part-time</option>
                <option value="Contract">Contract</option>
                <option value="Freelance">Freelance</option>
              </select>
              <font-awesome-icon icon="chevron-down" class="select-icon" />
            </div>
          </div>

          <div class="filter-group">
            <label><font-awesome-icon icon="user-graduate" /> Experience</label>
            <div class="select-wrapper">
              <select v-model="filters.experienceLevel" @change="filterJobs">
                <option value="">All Levels</option>
                <option value="Entry">Entry Level</option>
                <option value="Mid">Mid Level</option>
                <option value="Senior">Senior Level</option>
              </select>
              <font-awesome-icon icon="chevron-down" class="select-icon" />
            </div>
          </div>

          <div class="filter-group">
            <label><font-awesome-icon icon="map-marker-alt" /> Location</label>
            <div class="select-wrapper">
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
              <font-awesome-icon icon="chevron-down" class="select-icon" />
            </div>
          </div>

          <button
            class="btn btn-outline btn-sm filter-reset"
            @click="resetFilters"
          >
            <font-awesome-icon icon="undo" class="btn-icon" /> Reset Filters
          </button>
        </div>
      </section>

      <!-- Jobs Listing Section -->
      <section class="jobs-listing">
        <div v-if="loading" class="loading-state">
          <div class="spinner-container">
            <div class="spinner-circle"></div>
            <div class="spinner-circle-dot"></div>
          </div>
          <p>Loading job opportunities...</p>
        </div>

        <div v-else-if="error" class="error-message">
          <font-awesome-icon icon="exclamation-circle" class="error-icon pulse" />
          <p>{{ error }}</p>
          <button @click="fetchJobs" class="btn btn-outline btn-sm">
            <font-awesome-icon icon="sync" /> Try Again
          </button>
        </div>

        <div v-else-if="filteredJobs.length === 0" class="no-jobs">
          <div class="empty-illustration">
            <font-awesome-icon icon="briefcase" class="no-jobs-icon" />
            <div class="empty-dots">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
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
              <font-awesome-icon icon="filter-circle-xmark" /> Clear Filters
            </button>
            <router-link to="/post-job" class="btn btn-primary btn-sm">
              <font-awesome-icon icon="plus" /> Post a Job
            </router-link>
          </div>
        </div>

        <div v-else class="jobs-grid">
          <transition-group name="job-card" tag="div" class="jobs-grid-inner">
            <div 
              v-for="(job, index) in paginatedJobs" 
              :key="job.id" 
              class="job-card"
              :style="{ animationDelay: `${index * 0.1}s` }"
            >
              <div class="job-card-header">
                <div class="job-badge" :class="getJobTypeClass(job.jobType)">
                  <font-awesome-icon :icon="getJobTypeIcon(job.jobType)" />
                </div>
                <h3 class="job-title">{{ job.title }}</h3>
                <span class="job-company">
                  <font-awesome-icon icon="building" />
                  {{ job.company }}
                </span>
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
                  <span class="job-date">
                    <font-awesome-icon icon="calendar-alt" />
                    Posted {{ job.daysSincePosting() }} days ago
                  </span>
                  <span class="job-level">
                    <font-awesome-icon icon="signal" />
                    {{ job.experienceLevel }} Level
                  </span>
                </div>

                <router-link
                  :to="`/jobs/${job.id}`"
                  class="btn btn-primary btn-sm view-job-btn"
                >
                  View Details
                  <font-awesome-icon icon="arrow-right" class="btn-icon-right" />
                </router-link>
              </div>
              
              <div class="job-card-shine"></div>
            </div>
          </transition-group>
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

          <div class="pagination-numbers">
            <button 
              v-for="page in paginationRange" 
              :key="page" 
              @click="changePage(page)"
              class="page-number"
              :class="{ active: currentPage === page }"
            >
              {{ page }}
            </button>
          </div>

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
import { useDark, useToggle } from '@vueuse/core';

export default {
  name: "JobPostings",
  setup() {
    const isDark = useDark();
    const toggleDark = useToggle(isDark);
    
    return {
      isDark,
      toggleDark
    };
  },
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
    paginationRange() {
      const range = [];
      const maxVisiblePages = 5;
      
      if (this.totalPages <= maxVisiblePages) {
        // Show all pages if there are fewer than maxVisiblePages
        for (let i = 1; i <= this.totalPages; i++) {
          range.push(i);
        }
      } else {
        // Always include first page
        range.push(1);
        
        // Calculate start and end of visible range
        let start = Math.max(2, this.currentPage - 1);
        let end = Math.min(this.totalPages - 1, this.currentPage + 1);
        
        // Adjust if at the beginning or end
        if (this.currentPage <= 2) {
          end = Math.min(this.totalPages - 1, 4);
        } else if (this.currentPage >= this.totalPages - 1) {
          start = Math.max(2, this.totalPages - 3);
        }
        
        // Add visible range
        for (let i = start; i <= end; i++) {
          range.push(i);
        }
        
        // Always include last page
        if (this.totalPages > 1) {
          range.push(this.totalPages);
        }
      }
      
      return range;
    }
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
    
    getJobTypeIcon(jobType) {
      const icons = {
        "Full-time": "business-time",
        "Part-time": "clock",
        "Contract": "file-contract",
        "Freelance": "laptop-house"
      };
      return icons[jobType] || "briefcase";
    },

    truncateDescription(description, maxLength = 150) {
      if (description.length <= maxLength) return description;
      return description.substring(0, maxLength) + "...";
    },

    changePage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
        // Scroll to top of jobs listing with smooth animation
        const jobsListing = document.querySelector('.jobs-listing');
        if (jobsListing) {
          jobsListing.scrollIntoView({ behavior: 'smooth' });
        }
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
/* Base Styles */
.job-postings-page {
  padding: 2rem 0 4rem;
  min-height: 100vh;
  background-color: #f8fafc;
  transition: background-color 0.3s ease, color 0.3s ease;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
  position: relative;
}

/* Dark Theme */
.dark-theme {
  background-color: #0f172a;
  color: #e2e8f0;
}

.dark-theme .jobs-hero h1,
.dark-theme .animated-title .title-word {
  color: #e2e8f0;
  text-shadow: 0 2px 10px rgba(79, 70, 229, 0.3);
}

.dark-theme .subtitle {
  color: #94a3b8;
}

.dark-theme .search-input {
  background-color: #1e293b;
  border-color: #334155;
  color: #e2e8f0;
}

.dark-theme .search-icon,
.dark-theme .clear-search {
  color: #94a3b8;
}

.dark-theme .filter-bar {
  background-color: #1e293b;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2);
}

.dark-theme .filter-group label {
  color: #e2e8f0;
}

.dark-theme .select-wrapper select {
  background-color: #334155;
  border-color: #475569;
  color: #e2e8f0;
}

.dark-theme .select-icon {
  color: #94a3b8;
}

.dark-theme .job-card {
  background-color: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.dark-theme .job-card-header,
.dark-theme .job-card-footer {
  border-color: #334155;
}

.dark-theme .job-title {
  color: #e2e8f0;
}

.dark-theme .job-description {
  color: #94a3b8;
}

.dark-theme .skill-tag {
  background-color: #334155;
  color: #e2e8f0;
}

.dark-theme .pagination-btn,
.dark-theme .page-number {
  background-color: #1e293b;
  border-color: #334155;
  color: #e2e8f0;
}

.dark-theme .page-number.active {
  background-color: #4f46e5;
  color: white;
}

/* Theme Toggle */
.theme-toggle {
  position: absolute;
  top: 1rem;
  right: 1rem;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #3a0ca3);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  font-size: 1rem;
}

.theme-toggle:hover {
  transform: translateY(-3px) rotate(15deg);
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.15);
}

/* Hero section with animated title */
.jobs-hero {
  text-align: center;
  margin-bottom: 2.5rem;
  padding-top: 1rem;
}

.animated-title {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  font-size: 2.5rem;
  font-weight: 800;
  margin-bottom: 0.75rem;
  color: #1e293b;
  gap: 0.5rem;
}

.title-word {
  animation: color-animation 4s linear infinite;
  display: inline-block;
}

.title-word:nth-child(1) {
  animation-delay: 0s;
}

.title-word:nth-child(2) {
  animation-delay: 1s;
}

.title-word:nth-child(3) {
  animation-delay: 2s;
}

.title-word:nth-child(4) {
  animation-delay: 3s;
}

@keyframes color-animation {
  0% {
    color: #4f46e5;
  }
  32% {
    color: #4f46e5;
  }
  33% {
    color: #1e293b;
  }
  99% {
    color: #1e293b;
  }
  100% {
    color: #4f46e5;
  }
}

.subtitle {
  font-size: 1.125rem;
  color: #64748b;
  margin-bottom: 2rem;
  animation: fadeIn 1s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.search-container {
  margin-bottom: 1.5rem;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  animation: fadeIn 1s ease-in-out 0.2s both;
}

.search-input-group {
  position: relative;
  flex: 1;
  max-width: 600px;
  transition: all 0.3s ease;
}

.search-input-group:focus-within {
  transform: translateY(-2px);
}

.search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #64748b;
  font-size: 1rem;
  transition: color 0.3s ease;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 2px solid #e2e8f0;
  border-radius: 9999px;
  font-size: 1rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.search-input:focus {
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2), 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  outline: none;
}

.clear-search {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s ease;
}

.clear-search:hover {
  color: #1e293b;
  transform: translateY(-50%) scale(1.2);
}

.post-job-btn {
  white-space: nowrap;
  padding: 0.75rem 1.5rem;
  border-radius: 9999px;
  background: linear-gradient(135deg, #4f46e5, #3a0ca3);
  color: white;
  font-weight: 600;
  border: none;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.post-job-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.15);
}

.btn-icon {
  transition: transform 0.3s ease;
}

.post-job-btn:hover .btn-icon {
  transform: rotate(90deg);
}

.btn-icon-right {
  transition: transform 0.3s ease;
}

.view-job-btn:hover .btn-icon-right {
  transform: translateX(3px);
}

/* Filter bar */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
  margin-bottom: 2.5rem;
  padding: 1.25rem;
  background-color: white;
  border-radius: 1rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  animation: fadeIn 1s ease-in-out 0.4s both;
}

.filter-group {
  display: flex;
  flex-direction: column;
  min-width: 180px;
}

.filter-group label {
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  color: #1e293b;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.select-wrapper {
  position: relative;
}

.select-wrapper select {
  width: 100%;
  padding: 0.625rem 2rem 0.625rem 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  background-color: white;
  font-size: 0.875rem;
  color: #1e293b;
  cursor: pointer;
  appearance: none;
  transition: all 0.3s ease;
}

.select-wrapper select:focus {
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2);
  outline: none;
}

.select-icon {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #64748b;
  pointer-events: none;
  transition: transform 0.3s ease;
}

.select-wrapper:hover .select-icon {
  transform: translateY(-50%) translateY(2px);
}

.filter-reset {
  align-self: flex-end;
  margin-top: 1.5rem;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  background-color: transparent;
  color: #4f46e5;
  border: 1px solid #4f46e5;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-reset:hover {
  background-color: #eef2ff;
  transform: translateY(-2px);
}

/* Jobs Grid with animations */
.jobs-grid {
  margin-bottom: 2.5rem;
}

.jobs-grid-inner {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.job-card {
  background-color: white;
  border-radius: 1rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  border-top: 4px solid #4f46e5;
  height: 100%;
  position: relative;
  animation: cardAppear 0.5s ease-out both;
  transform-origin: center;
}

@keyframes cardAppear {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.job-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.job-card-header {
  padding: 1.25rem;
  border-bottom: 1px solid #f1f5f9;
  position: relative;
}

.job-badge {
  position: absolute;
  top: -4px;
  right: 1.25rem;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1rem;
  transform: translateY(-50%);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.job-badge.full-time {
  background: linear-gradient(135deg, #0284c7, #0369a1);
}

.job-badge.part-time {
  background: linear-gradient(135deg, #d97706, #b45309);
}

.job-badge.contract {
  background: linear-gradient(135deg, #7c3aed, #6d28d9);
}

.job-badge.freelance {
  background: linear-gradient(135deg, #16a34a, #15803d);
}

.job-title {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  color: #1e293b;
  font-weight: 700;
  transition: color 0.3s ease;
}

.job-company {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  color: #4f46e5;
  margin-bottom: 0.75rem;
  font-weight: 500;
}

.job-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.job-location {
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.job-type {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
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
  padding: 1.25rem;
  flex: 1;
}

.job-description {
  color: #64748b;
  margin-bottom: 1rem;
  font-size: 0.875rem;
  line-height: 1.5;
}

.job-skills {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.skill-tag {
  background-color: #f1f5f9;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  color: #334155;
  transition: all 0.3s ease;
}

.skill-tag:hover {
  background-color: #e0f2fe;
  transform: translateY(-2px);
}

.job-card-footer {
  padding: 1.25rem;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.job-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.job-date,
.job-level {
  font-size: 0.75rem;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.view-job-btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  background-color: #4f46e5;
  color: white;
  font-size: 0.875rem;
  font-weight: 500;
  border: none;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.view-job-btn:hover {
  background-color: #4338ca;
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

/* Shine effect */
.job-card-shine {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.4) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  z-index: 1;
  transform: translateX(-100%) rotate(45deg);
  pointer-events: none;
  transition: transform 0.5s;
}

.job-card:hover .job-card-shine {
  transform: translateX(100%) rotate(45deg);
  transition: transform 0.8s;
}

/* Loading and Error States */
.loading-state,
.error-message,
.no-jobs {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
  min-height: 300px;
}

.spinner-container {
  position: relative;
  width: 60px;
  height: 60px;
  margin-bottom: 1.5rem;
}

.spinner-circle {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 4px solid transparent;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1.2s linear infinite;
}

.spinner-circle-dot {
  position: absolute;
  top: 5px;
  left: 5px;
  right: 5px;
  bottom: 5px;
  border: 4px solid transparent;
  border-top-color: #818cf8;
  border-radius: 50%;
  animation: spin 0.8s linear infinite reverse;
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
  color: #ef4444;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  color: #ef4444;
}

.pulse {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.95);
    opacity: 0.8;
  }
  50% {
    transform: scale(1.05);
    opacity: 1;
  }
  100% {
    transform: scale(0.95);
    opacity: 0.8;
  }
}

.empty-illustration {
  position: relative;
  margin-bottom: 1.5rem;
}

.no-jobs-icon {
  font-size: 4rem;
  color: #94a3b8;
  opacity: 0.6;
}

.empty-dots {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 1rem;
}

.empty-dots span {
  width: 8px;
  height: 8px;
  background-color: #94a3b8;
  border-radius: 50%;
  opacity: 0.6;
}

.empty-dots span:nth-child(1) {
  animation: bounce 1s infinite 0.2s;
}

.empty-dots span:nth-child(2) {
  animation: bounce 1s infinite 0.4s;
}

.empty-dots span:nth-child(3) {
  animation: bounce 1s infinite 0.6s;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.no-jobs h3 {
  font-size: 1.5rem;
  color: #1e293b;
  margin-bottom: 0.75rem;
}

.no-jobs p {
  color: #64748b;
  margin-bottom: 1.5rem;
  max-width: 400px;
}

.no-jobs-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.75rem;
  margin-top: 2.5rem;
}

.pagination-btn {
  background-color: white;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #eef2ff;
  border-color: #4f46e5;
  color: #4f46e5;
  transform: translateY(-2px);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-numbers {
  display: flex;
  gap: 0.5rem;
}

.page-number {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
  background-color: white;
  font-size: 0.875rem;
  color: #1e293b;
  cursor: pointer;
  transition: all 0.3s ease;
}

.page-number:hover:not(.active) {
  background-color: #eef2ff;
  border-color: #4f46e5;
  color: #4f46e5;
  transform: translateY(-2px);
}

.page-number.active {
  background-color: #4f46e5;
  border-color: #4f46e5;
  color: white;
  font-weight: 600;
}

/* Transition group animations */
.job-card-enter-active, 
.job-card-leave-active {
  transition: all 0.5s ease;
}

.job-card-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.job-card-leave-to {
  opacity: 0;
  transform: translateY(-30px);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .animated-title {
    font-size: 2rem;
  }
  
  .search-container {
    flex-direction: column;
    gap: 1rem;
  }

  .post-job-btn {
    width: 100%;
  }

  .filter-bar {
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }

  .filter-group {
    width: 100%;
  }

  .filter-reset {
    align-self: center;
    width: 100%;
  }
  
  .pagination-numbers {
    display: none;
  }
  
  .pagination-info {
    display: block;
    color: #64748b;
    font-size: 0.875rem;
  }
}

@media (max-width: 480px) {
  .animated-title {
    font-size: 1.5rem;
  }
  
  .subtitle {
    font-size: 1rem;
  }
  
  .job-card-header {
    padding: 1rem;
  }
  
  .job-badge {
    width: 30px;
    height: 30px;
    font-size: 0.875rem;
  }
  
  .job-title {
    font-size: 1.125rem;
    padding-right: 1.5rem;
  }
  
  .job-card-body,
  .job-card-footer {
    padding: 1rem;
  }
  
  .view-job-btn {
    padding: 0.375rem 0.75rem;
    font-size: 0.75rem;
  }
}
</style>