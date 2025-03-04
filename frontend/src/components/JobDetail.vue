<template>
  <div class="job-detail-page">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading job details...</p>
      </div>

      <div v-else-if="error" class="error-message">
        <font-awesome-icon icon="exclamation-circle" />
        <p>{{ error }}</p>
        <div class="error-actions">
          <router-link to="/jobs" class="btn btn-outline btn-sm">
            Back to Jobs
          </router-link>
          <button @click="fetchJob" class="btn btn-primary btn-sm">
            Try Again
          </button>
        </div>
      </div>

      <template v-else-if="job">
        <!-- Back navigation -->
        <div class="back-nav">
          <router-link to="/jobs" class="back-link">
            <font-awesome-icon icon="arrow-left" />
            <span>Back to Job Listings</span>
          </router-link>
        </div>

        <!-- Job header section -->
        <section class="job-header">
          <div class="job-header-content">
            <div class="job-title-section">
              <h1>{{ job.title }}</h1>
              <div class="job-meta">
                <div class="job-company">
                  <font-awesome-icon icon="building" />
                  <span>{{ job.company }}</span>
                </div>
                <div class="job-location">
                  <font-awesome-icon icon="map-marker-alt" />
                  <span>{{ job.location }}</span>
                </div>
                <div class="job-type" :class="getJobTypeClass(job.jobType)">
                  <span>{{ job.jobType }}</span>
                </div>
              </div>
            </div>
            <div class="job-actions">
              <button class="btn btn-primary apply-btn">
                <font-awesome-icon icon="paper-plane" />
                Apply Now
              </button>
              <button class="btn btn-outline save-btn">
                <font-awesome-icon icon="bookmark" />
                Save Job
              </button>
            </div>
          </div>
        </section>

        <!-- Job details section -->
        <div class="job-content-container">
          <div class="job-details-grid">
            <!-- Main job details -->
            <div class="job-details-main">
              <section class="job-section">
                <h2>Job Description</h2>
                <div class="job-description">
                  <p>{{ job.description }}</p>
                </div>
              </section>

              <section class="job-section">
                <h2>Required Skills</h2>
                <div class="job-skills">
                  <div
                    v-for="(skill, index) in job.skillsArray()"
                    :key="index"
                    class="skill-tag"
                  >
                    {{ skill }}
                  </div>
                </div>
              </section>

              <section class="job-section">
                <h2>How to Apply</h2>
                <div class="apply-instructions">
                  <p>
                    To apply for this position, please click the "Apply Now"
                    button above and follow the instructions to submit your
                    application. For any questions, please contact the hiring
                    manager at
                    <a :href="`mailto:${job.contactEmail}`">{{
                      job.contactEmail
                    }}</a
                    >.
                  </p>
                </div>
              </section>
            </div>

            <!-- Job sidebar -->
            <div class="job-details-sidebar">
              <div class="job-sidebar-card">
                <h3>Job Overview</h3>
                <ul class="job-overview-list">
                  <li>
                    <div class="overview-icon">
                      <font-awesome-icon icon="calendar-alt" />
                    </div>
                    <div class="overview-content">
                      <span class="overview-label">Posted Date</span>
                      <span class="overview-value">{{
                        job.formattedDate()
                      }}</span>
                    </div>
                  </li>
                  <li>
                    <div class="overview-icon">
                      <font-awesome-icon icon="briefcase" />
                    </div>
                    <div class="overview-content">
                      <span class="overview-label">Job Type</span>
                      <span class="overview-value">{{ job.jobType }}</span>
                    </div>
                  </li>
                  <li>
                    <div class="overview-icon">
                      <font-awesome-icon icon="graduation-cap" />
                    </div>
                    <div class="overview-content">
                      <span class="overview-label">Experience Level</span>
                      <span class="overview-value">{{
                        job.experienceLevel
                      }}</span>
                    </div>
                  </li>
                  <li v-if="job.salaryRange">
                    <div class="overview-icon">
                      <font-awesome-icon icon="money-bill-alt" />
                    </div>
                    <div class="overview-content">
                      <span class="overview-label">Salary Range</span>
                      <span class="overview-value">{{ job.salaryRange }}</span>
                    </div>
                  </li>
                </ul>
              </div>

              <div class="job-sidebar-card">
                <h3>About the Company</h3>
                <p class="company-description">
                  {{ job.company }} is a forward-thinking organization that
                  values innovation and collaboration. Join our team to work on
                  exciting projects and grow your career.
                </p>
                <button class="btn btn-outline btn-sm btn-full company-btn">
                  <font-awesome-icon icon="external-link-alt" />
                  Visit Company Profile
                </button>
              </div>

              <div class="job-sidebar-card">
                <h3>Job Posted By</h3>
                <div class="posted-by-user">
                  <div class="user-avatar">
                    <font-awesome-icon icon="user" />
                  </div>
                  <div class="user-info">
                    <span class="user-name">{{ job.postedByName }}</span>
                    <button class="btn btn-outline btn-sm contact-btn">
                      <font-awesome-icon icon="comment" />
                      Contact
                    </button>
                  </div>
                </div>
              </div>

              <div class="share-job">
                <h3>Share This Job</h3>
                <div class="share-buttons">
                  <button class="share-btn" style="background-color: #3b5998">
                    Share
                  </button>
                  <button class="share-btn" style="background-color: #1da1f2">
                    Tweet
                  </button>
                  <button class="share-btn" style="background-color: #0077b5">
                    Post
                  </button>
                  <button class="share-btn" style="background-color: #ea4335">
                    <font-awesome-icon icon="envelope" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Similar jobs section -->
        <section class="similar-jobs-section">
          <h2>Similar Jobs</h2>
          <div class="similar-jobs-grid">
            <div
              v-for="similarJob in similarJobs"
              :key="similarJob.id"
              class="job-card"
            >
              <div class="job-card-header">
                <h3 class="job-title">{{ similarJob.title }}</h3>
                <span class="job-company">{{ similarJob.company }}</span>
                <div class="job-meta">
                  <span class="job-location">
                    <font-awesome-icon icon="map-marker-alt" />
                    {{ similarJob.location }}
                  </span>
                  <span
                    class="job-type"
                    :class="getJobTypeClass(similarJob.jobType)"
                  >
                    {{ similarJob.jobType }}
                  </span>
                </div>
              </div>

              <div class="job-card-body">
                <p class="job-description">
                  {{ truncateDescription(similarJob.description) }}
                </p>

                <div class="job-skills">
                  <span
                    v-for="(skill, index) in similarJob.skillsArray()"
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
                    >Posted {{ similarJob.daysSincePosting() }} days ago</span
                  >
                  <span class="job-level"
                    >{{ similarJob.experienceLevel }} Level</span
                  >
                </div>

                <router-link
                  :to="`/jobs/${similarJob.id}`"
                  class="btn btn-primary btn-sm view-job-btn"
                >
                  View Details
                </router-link>
              </div>
            </div>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script>
import JobPost from "@/models/JobPost";

export default {
  name: "JobDetail",
  data() {
    return {
      job: null,
      loading: true,
      error: null,
      similarJobs: [],
    };
  },
  created() {
    this.fetchJob();
  },
  watch: {
    $route(to, from) {
      // Refetch job data if route param changes (user clicked on a similar job)
      if (to.params.id !== from.params.id) {
        this.fetchJob();
      }
    },
  },
  methods: {
    async fetchJob() {
      this.loading = true;
      this.error = null;
      const jobId = parseInt(this.$route.params.id);

      try {
        // In a real implementation, this would be an API call
        // For now, use mock data
        setTimeout(() => {
          const allJobs = this.getMockJobs().map((job) => new JobPost(job));
          this.job = allJobs.find((j) => j.id === jobId);

          if (!this.job) {
            this.error = `Job with ID ${jobId} not found`;
          } else {
            // Find similar jobs based on skills or job type
            this.findSimilarJobs(allJobs);
          }

          this.loading = false;
        }, 800);

        // Real implementation would be:
        // const response = await axios.get(`/api/jobs/${jobId}`);
        // this.job = new JobPost(response.data);
      } catch (error) {
        console.error("Error fetching job details:", error);
        this.error = "Failed to load job details. Please try again.";
        this.loading = false;
      }
    },

    findSimilarJobs(allJobs) {
      // Filter out the current job
      const otherJobs = allJobs.filter((j) => j.id !== this.job.id);

      // Score each job based on similarity
      const scoredJobs = otherJobs.map((j) => {
        let score = 0;

        // Same job type gets 2 points
        if (j.jobType === this.job.jobType) score += 2;

        // Same experience level gets 1 point
        if (j.experienceLevel === this.job.experienceLevel) score += 1;

        // Each matching skill gets 2 points
        const currentSkills = this.job.skillsArray();
        const otherSkills = j.skillsArray();

        for (const skill of otherSkills) {
          if (currentSkills.includes(skill)) score += 2;
        }

        return { job: j, score };
      });

      // Sort by score and take top 3
      this.similarJobs = scoredJobs
        .sort((a, b) => b.score - a.score)
        .slice(0, 3)
        .map((item) => item.job);
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

    truncateDescription(description, maxLength = 120) {
      if (description.length <= maxLength) return description;
      return description.substring(0, maxLength) + "...";
    },

    getMockJobs() {
      return [
        {
          id: 1,
          title: "Frontend Developer",
          company: "Tech Innovators",
          location: "San Francisco, CA",
          description:
            "We are looking for a skilled Frontend Developer to join our team. You will be responsible for building web applications using Vue.js and modern web technologies. The ideal candidate should have experience with JavaScript frameworks, HTML5, CSS3, and responsive design principles.\n\nResponsibilities include implementing user interface components, collaborating with UX designers and backend developers, optimizing applications for maximum speed and scalability, and ensuring cross-browser compatibility. You should be comfortable working in an agile environment with frequent iterations.",
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
            "Join our design team and help create beautiful and functional user interfaces for our clients. You should have a strong portfolio and experience with design tools. The ideal candidate will have a keen eye for aesthetics, understanding of user experience principles, and the ability to transform complex requirements into intuitive interfaces.\n\nResponsibilities include creating wireframes, prototypes, user flows, and visual designs. You will work closely with product managers and developers to ensure designs are implemented correctly.",
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
            "We need a talented Content Writer to create engaging content for our blog and social media channels. Must have excellent writing skills and SEO knowledge. The ideal candidate will have a way with words, the ability to research complex topics, and strong editing skills.\n\nResponsibilities include creating compelling blog posts, social media content, newsletter copy, and website text. You will work with our marketing team to develop effective content strategies that drive engagement and conversions.",
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
            "Looking for a Backend Developer with Go expertise to help build our next-generation API services. Must have experience with database design and RESTful APIs. The ideal candidate will have strong problem-solving skills, experience with cloud infrastructure, and a passion for writing clean, maintainable code.\n\nResponsibilities include designing and implementing APIs, optimizing database performance, integrating with third-party services, and ensuring security best practices are followed throughout the codebase.",
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
            "Join our marketing team to develop and implement marketing strategies. You should have experience with digital marketing and analytics tools. The ideal candidate will be data-driven, creative, and able to work across multiple marketing channels.\n\nResponsibilities include running digital marketing campaigns, analyzing performance metrics, managing social media accounts, and collaborating with content creators to develop marketing materials.",
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
            "We need a skilled mobile developer who can build native iOS applications. Knowledge of Swift and the Apple ecosystem is required. The ideal candidate will have a portfolio of published apps, experience with the App Store submission process, and a deep understanding of iOS design patterns.\n\nResponsibilities include developing new features, fixing bugs, optimizing performance, and working with designers to implement user interfaces.",
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
            "Looking for a Data Scientist to join our team. You will analyze large datasets and build machine learning models to solve business problems. The ideal candidate will have strong statistical knowledge, programming skills, and the ability to communicate complex findings to non-technical stakeholders.\n\nResponsibilities include exploratory data analysis, feature engineering, model development, and creating data visualizations. You will work with cross-functional teams to implement machine learning solutions that drive business value.",
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
            "We are looking for a Product Manager to lead product development and work with cross-functional teams to deliver great user experiences. The ideal candidate will have a blend of business acumen, technical understanding, and user empathy.\n\nResponsibilities include defining product vision and strategy, managing the product roadmap, gathering and prioritizing requirements, and working closely with engineering, design, and marketing teams throughout the product lifecycle.",
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
            "Join our team to build and maintain CI/CD pipelines and cloud infrastructure. Experience with AWS and containerization is required. The ideal candidate will have a strong understanding of infrastructure as code, automation tools, and security best practices.\n\nResponsibilities include designing and implementing cloud infrastructure, automating deployment processes, monitoring system performance, and troubleshooting issues in production environments.",
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
            "We need a creative Graphic Designer to join our team. You will create visual concepts for web and print materials. The ideal candidate will have a strong portfolio showing their design skills, creativity, and attention to detail.\n\nResponsibilities include creating branding materials, marketing collateral, social media graphics, and website design elements. You will need to maintain brand consistency while producing fresh and engaging designs.",
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
.job-detail-page {
  padding-bottom: var(--space-12);
}

/* Loading and Error States */
.loading-state,
.error-message {
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

.error-message svg {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
}

.error-actions {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-4);
}

/* Back navigation */
.back-nav {
  margin-bottom: var(--space-6);
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--primary-color);
  font-weight: var(--font-weight-medium);
  transition: color var(--transition-fast);
}

.back-link:hover {
  color: var(--primary-dark);
  text-decoration: none;
}

/* Job header */
.job-header {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-6);
  margin-bottom: var(--space-6);
}

.job-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: var(--space-4);
}

.job-title-section {
  flex: 1;
  min-width: 250px;
}

.job-title-section h1 {
  font-size: var(--font-size-2xl);
  margin-bottom: var(--space-2);
  color: var(--dark);
}

.job-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-4);
  margin-top: var(--space-3);
}

.job-company,
.job-location {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--medium);
}

.job-type {
  display: inline-block;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--font-size-sm);
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

.job-actions {
  display: flex;
  gap: var(--space-3);
}

.apply-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.save-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

/* Job content layout */
.job-content-container {
  margin-bottom: var(--space-8);
}

.job-details-grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
  gap: var(--space-6);
}

/* Main job details */
.job-details-main {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.job-section {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-6);
}

.job-section h2 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-4);
  color: var(--dark);
  padding-bottom: var(--space-2);
  border-bottom: 1px solid var(--light);
}

.job-description p {
  color: var(--medium);
  line-height: 1.6;
  white-space: pre-line;
}

.job-skills {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.skill-tag {
  background-color: var(--light);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--font-size-sm);
  color: var(--dark);
}

.apply-instructions p {
  color: var(--medium);
  line-height: 1.6;
}

.apply-instructions a {
  color: var(--primary-color);
  font-weight: var(--font-weight-medium);
}

/* Job sidebar */
.job-details-sidebar {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.job-sidebar-card {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-4);
}

.job-sidebar-card h3 {
  font-size: var(--font-size-md);
  margin-bottom: var(--space-3);
  padding-bottom: var(--space-2);
  border-bottom: 1px solid var(--light);
}

.job-overview-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.job-overview-list li {
  display: flex;
  align-items: center;
  padding: var(--space-2) 0;
  border-bottom: 1px solid var(--light);
}

.job-overview-list li:last-child {
  border-bottom: none;
}

.overview-icon {
  width: 40px;
  height: 40px;
  background-color: var(--primary-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary-color);
  margin-right: var(--space-3);
}

.overview-content {
  display: flex;
  flex-direction: column;
}

.overview-label {
  font-size: var(--font-size-xs);
  color: var(--medium);
}

.overview-value {
  font-weight: var(--font-weight-medium);
  color: var(--dark);
}

.company-description {
  color: var(--medium);
  margin-bottom: var(--space-3);
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.btn-full {
  width: 100%;
}

.company-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
}

.posted-by-user {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: var(--primary-light);
  color: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-lg);
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.user-name {
  font-weight: var(--font-weight-medium);
  color: var(--dark);
}

.contact-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-xs);
  padding: var(--space-1) var(--space-2);
}

.share-job h3 {
  margin-bottom: var(--space-3);
}

.share-buttons {
  display: flex;
  gap: var(--space-2);
}

.share-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  border: none;
  cursor: pointer;
  transition: opacity var(--transition-fast);
  margin-right: var(--space-2);
  font-size: var(--font-size-xs);
}

.share-btn:hover {
  opacity: 0.9;
}

/* Similar jobs section */
.similar-jobs-section {
  margin-top: var(--space-8);
}

.similar-jobs-section h2 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-4);
  color: var(--dark);
}

.similar-jobs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-4);
}

/* Job card styles */
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

/* Responsive adjustments */
@media (max-width: 992px) {
  .job-details-grid {
    grid-template-columns: 1fr;
  }

  .job-actions {
    margin-top: var(--space-3);
    width: 100%;
  }

  .apply-btn,
  .save-btn {
    flex: 1;
  }

  .similar-jobs-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 576px) {
  .job-meta {
    flex-direction: column;
    gap: var(--space-2);
  }

  .job-actions {
    flex-direction: column;
    gap: var(--space-2);
  }
}
</style>
