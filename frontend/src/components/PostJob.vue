<template>
  <div class="post-job-page">
    <div class="container">
      <div class="post-job-header">
        <h1>{{ isEditing ? "Edit Job Posting" : "Post a New Job" }}</h1>
        <p class="subtitle">
          {{
            isEditing
              ? "Update the details of your job posting"
              : "Share an opportunity with the SkillSwap community"
          }}
        </p>
      </div>

      <div class="form-container">
        <!-- Form stepper -->
        <div class="form-stepper">
          <div
            v-for="(step, index) in steps"
            :key="index"
            class="stepper-item"
            :class="{
              active: currentStep === index,
              completed: currentStep > index,
            }"
            @click="goToStep(index)"
          >
            <div class="step-counter">
              <span v-if="currentStep <= index">{{ index + 1 }}</span>
              <font-awesome-icon v-else icon="check" />
            </div>
            <div class="step-name">{{ step }}</div>
          </div>
        </div>

        <!-- Job posting form -->
        <form @submit.prevent="submitForm" class="job-form">
          <!-- Step 1: Basic Information -->
          <div v-show="currentStep === 0" class="form-step">
            <h2>Basic Information</h2>

            <div class="form-group">
              <label for="job-title">Job Title *</label>
              <input
                id="job-title"
                v-model="formData.title"
                type="text"
                placeholder="e.g. Frontend Developer"
                required
                :class="{ 'input-error': errors.title }"
              />
              <span v-if="errors.title" class="error-message">{{
                errors.title
              }}</span>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="job-company">Company Name *</label>
                <input
                  id="job-company"
                  v-model="formData.company"
                  type="text"
                  placeholder="e.g. Tech Innovators"
                  required
                  :class="{ 'input-error': errors.company }"
                />
                <span v-if="errors.company" class="error-message">{{
                  errors.company
                }}</span>
              </div>

              <div class="form-group">
                <label for="job-location">Location *</label>
                <input
                  id="job-location"
                  v-model="formData.location"
                  type="text"
                  placeholder="e.g. San Francisco, CA or Remote"
                  required
                  :class="{ 'input-error': errors.location }"
                />
                <span v-if="errors.location" class="error-message">{{
                  errors.location
                }}</span>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="job-type">Job Type *</label>
                <select
                  id="job-type"
                  v-model="formData.jobType"
                  required
                  :class="{ 'input-error': errors.jobType }"
                >
                  <option value="">Select Job Type</option>
                  <option value="Full-time">Full-time</option>
                  <option value="Part-time">Part-time</option>
                  <option value="Contract">Contract</option>
                  <option value="Freelance">Freelance</option>
                </select>
                <span v-if="errors.jobType" class="error-message">{{
                  errors.jobType
                }}</span>
              </div>

              <div class="form-group">
                <label for="job-experience">Experience Level *</label>
                <select
                  id="job-experience"
                  v-model="formData.experienceLevel"
                  required
                  :class="{ 'input-error': errors.experienceLevel }"
                >
                  <option value="">Select Experience Level</option>
                  <option value="Entry">Entry Level</option>
                  <option value="Mid">Mid Level</option>
                  <option value="Senior">Senior Level</option>
                </select>
                <span v-if="errors.experienceLevel" class="error-message">{{
                  errors.experienceLevel
                }}</span>
              </div>
            </div>

            <div class="form-group">
              <label for="job-salary">Salary Range</label>
              <input
                id="job-salary"
                v-model="formData.salaryRange"
                type="text"
                placeholder="e.g. $50,000 - $70,000 or $25 - $35 per hour"
              />
            </div>
          </div>

          <!-- Step 2: Job Description -->
          <div v-show="currentStep === 1" class="form-step">
            <h2>Job Description</h2>

            <div class="form-group">
              <label for="job-description">Detailed Description *</label>
              <textarea
                id="job-description"
                v-model="formData.description"
                rows="8"
                placeholder="Provide a comprehensive description of the job role, responsibilities, and requirements..."
                required
                :class="{ 'input-error': errors.description }"
              ></textarea>
              <div class="textarea-footer">
                <span v-if="errors.description" class="error-message">{{
                  errors.description
                }}</span>
                <span class="char-count"
                  >{{ formData.description.length }}/2000</span
                >
              </div>
            </div>

            <div class="form-group">
              <label for="job-skills">Required Skills *</label>
              <div
                class="tags-input-container"
                :class="{ 'input-error': errors.skillsRequired }"
              >
                <div class="tags-container">
                  <div
                    v-for="(skill, index) in formData.skillsRequired"
                    :key="index"
                    class="tag"
                  >
                    {{ skill }}
                    <button
                      type="button"
                      class="remove-tag"
                      @click="removeSkill(index)"
                    >
                      ×
                    </button>
                  </div>
                  <input
                    type="text"
                    v-model="newSkill"
                    @keydown.enter.prevent="addSkill"
                    @keydown="handleKeydown"
                    placeholder="Type skill and press Enter"
                    class="tags-input"
                  />
                </div>
              </div>
              <span v-if="errors.skillsRequired" class="error-message">{{
                errors.skillsRequired
              }}</span>
              <small class="helper-text"
                >Press Enter or comma after each skill</small
              >
            </div>
          </div>

          <!-- Step 3: Contact Information -->
          <div v-show="currentStep === 2" class="form-step">
            <h2>Contact Information</h2>

            <div class="form-group">
              <label for="job-email">Contact Email *</label>
              <input
                id="job-email"
                v-model="formData.contactEmail"
                type="email"
                placeholder="e.g. hiring@company.com"
                required
                :class="{ 'input-error': errors.contactEmail }"
              />
              <span v-if="errors.contactEmail" class="error-message">{{
                errors.contactEmail
              }}</span>
            </div>

            <div class="form-row">
              <div class="form-group checkbox-group">
                <input
                  type="checkbox"
                  id="job-notify"
                  v-model="notifyOnApplications"
                />
                <label for="job-notify">Notify me when someone applies</label>
              </div>
            </div>

            <div class="form-preview">
              <h3>Job Posting Preview</h3>
              <div class="preview-card">
                <div class="preview-header">
                  <h4>{{ formData.title || "Job Title" }}</h4>
                  <div class="preview-meta">
                    <span>{{ formData.company || "Company Name" }}</span>
                    <span>{{ formData.location || "Location" }}</span>
                  </div>
                </div>
                <div class="preview-body">
                  <p>
                    {{
                      truncateText(
                        formData.description || "No description provided.",
                        150,
                      )
                    }}
                  </p>
                  <div class="preview-skills">
                    <span
                      v-for="(skill, index) in formData.skillsRequired"
                      :key="index"
                      class="preview-skill"
                    >
                      {{ skill }}
                    </span>
                    <span
                      v-if="formData.skillsRequired.length === 0"
                      class="preview-skill empty"
                    >
                      No skills specified
                    </span>
                  </div>
                </div>
                <div class="preview-footer">
                  <span>{{ formData.jobType || "Job Type" }}</span>
                  <span
                    >{{
                      formData.experienceLevel || "Experience Level"
                    }}
                    Level</span
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- Form navigation buttons -->
          <div class="form-navigation">
            <button
              type="button"
              class="btn btn-outline btn-lg"
              @click="prevStep"
              v-if="currentStep > 0"
            >
              <font-awesome-icon icon="arrow-left" />
              Back
            </button>

            <button
              type="button"
              class="btn btn-primary btn-lg"
              @click="nextStep"
              v-if="currentStep < steps.length - 1"
            >
              Next
              <font-awesome-icon icon="arrow-right" />
            </button>

            <button
              type="submit"
              class="btn btn-primary btn-lg submit-btn"
              v-if="currentStep === steps.length - 1"
              :disabled="isSubmitting"
            >
              <span v-if="isSubmitting" class="spinner"></span>
              <span v-else>{{ isEditing ? "Update Job" : "Post Job" }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
// Import will be needed when connecting to real API
// import axios from 'axios';
import eventBus from "@/utils/eventBus";

export default {
  name: "PostJob",
  data() {
    return {
      isEditing: false,
      currentStep: 0,
      isSubmitting: false,
      steps: ["Basic Info", "Description", "Contact & Preview"],
      formData: {
        title: "",
        company: "",
        location: "",
        description: "",
        skillsRequired: [],
        experienceLevel: "",
        jobType: "",
        salaryRange: "",
        contactEmail: "",
      },
      newSkill: "",
      notifyOnApplications: true,
      errors: {},
    };
  },
  created() {
    // Check if we're editing an existing job
    const jobId = this.$route.params.id;
    if (jobId) {
      this.isEditing = true;
      this.fetchJobData(jobId);
    } else {
      // Set default contact email from user profile if available
      if (this.$store.state.user && this.$store.state.user.email) {
        this.formData.contactEmail = this.$store.state.user.email;
      }
    }
  },
  methods: {
    fetchJobData(jobId) {
      return new Promise((resolve, reject) => {
        try {
          // In a real implementation, this would call the API
          // For now, use mock data
          setTimeout(() => {
            const allJobs = this.getMockJobs();
            const job = allJobs.find((j) => j.id === parseInt(jobId));

            if (job) {
              this.formData = {
                title: job.title,
                company: job.company,
                location: job.location,
                description: job.description,
                skillsRequired: Array.isArray(job.skillsRequired)
                  ? job.skillsRequired
                  : job.skillsRequired.split(",").map((s) => s.trim()),
                experienceLevel: job.experienceLevel,
                jobType: job.jobType,
                salaryRange: job.salaryRange,
                contactEmail: job.contactEmail,
              };
              resolve();
            } else {
              this.$router.push("/jobs");
              reject(new Error("Job not found"));
            }
          }, 500);
        } catch (error) {
          reject(error);
        }
      });
    },

    validateStep() {
      this.errors = {};

      if (this.currentStep === 0) {
        // Validate basic info
        if (!this.formData.title.trim()) {
          this.errors.title = "Job title is required";
        }

        if (!this.formData.company.trim()) {
          this.errors.company = "Company name is required";
        }

        if (!this.formData.location.trim()) {
          this.errors.location = "Location is required";
        }

        if (!this.formData.jobType) {
          this.errors.jobType = "Job type is required";
        }

        if (!this.formData.experienceLevel) {
          this.errors.experienceLevel = "Experience level is required";
        }
      } else if (this.currentStep === 1) {
        // Validate description
        if (!this.formData.description.trim()) {
          this.errors.description = "Job description is required";
        } else if (this.formData.description.length < 50) {
          this.errors.description =
            "Description should be at least 50 characters";
        }

        if (this.formData.skillsRequired.length === 0) {
          this.errors.skillsRequired = "At least one skill is required";
        }
      } else if (this.currentStep === 2) {
        // Validate contact info
        if (!this.formData.contactEmail.trim()) {
          this.errors.contactEmail = "Contact email is required";
        } else if (!this.isValidEmail(this.formData.contactEmail)) {
          this.errors.contactEmail = "Please enter a valid email address";
        }
      }

      return Object.keys(this.errors).length === 0;
    },

    isValidEmail(email) {
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return re.test(email);
    },

    nextStep() {
      if (this.validateStep()) {
        this.currentStep += 1;
        window.scrollTo(0, 0);
      }
    },

    prevStep() {
      this.currentStep -= 1;
      window.scrollTo(0, 0);
    },

    goToStep(step) {
      if (step < this.currentStep) {
        this.currentStep = step;
        window.scrollTo(0, 0);
      } else if (step > this.currentStep) {
        // Only allow jumping ahead if current step is valid
        if (this.validateStep()) {
          this.currentStep = step;
          window.scrollTo(0, 0);
        }
      }
    },

    addSkill() {
      if (this.newSkill.trim()) {
        // Prevent duplicates
        if (!this.formData.skillsRequired.includes(this.newSkill.trim())) {
          this.formData.skillsRequired.push(this.newSkill.trim());
        }
        this.newSkill = "";
      }
    },

    handleKeydown(event) {
      // Handle comma key press for adding skills
      if (event.key === ",") {
        event.preventDefault();
        this.addSkill();
      }
    },

    removeSkill(index) {
      this.formData.skillsRequired.splice(index, 1);
    },

    truncateText(text, maxLength) {
      if (!text) return "";
      if (text.length <= maxLength) return text;

      // Get the substring exactly at max length
      const truncated = text.substring(0, maxLength);

      // Remove any trailing spaces and add ellipsis
      return truncated.trimEnd() + "...";
    },

    async submitForm() {
      if (!this.validateStep()) {
        return;
      }

      this.isSubmitting = true;

      try {
        // Prepare data
        // Note: We're not using this in our mock implementation, but will need it for API integration
        // const jobData = {
        //   ...this.formData,
        //   postedByUserID: this.$store.state.user ? this.$store.state.user.id : 1,
        //   postedByName: this.$store.state.user ? this.$store.state.user.name : 'Test User'
        // };

        // In a real implementation, this would call the API
        setTimeout(() => {
          this.isSubmitting = false;

          // Show success message
          eventBus.emit("show-notification", {
            type: "success",
            title: this.isEditing ? "Job Updated" : "Job Posted",
            message: this.isEditing
              ? "Your job posting has been successfully updated."
              : "Your job posting has been successfully published.",
            duration: 5000,
          });

          // Redirect to jobs page
          this.$router.push("/jobs");
        }, 1500);

        // Real implementation would be:
        // if (this.isEditing) {
        //   await axios.put(`/api/jobs/${this.$route.params.id}`, jobData);
        // } else {
        //   await axios.post('/api/jobs', jobData);
        // }
      } catch (error) {
        console.error("Error submitting job:", error);
        this.isSubmitting = false;

        // Show error message
        eventBus.emit("show-notification", {
          type: "error",
          title: "Error",
          message:
            "There was a problem submitting your job posting. Please try again.",
          duration: 5000,
        });
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
      ];
    },
  },
};
</script>

<style scoped>
.post-job-page {
  padding-bottom: var(--space-12);
}

.post-job-header {
  text-align: center;
  margin-bottom: var(--space-8);
}

.post-job-header h1 {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-2);
  color: var(--dark);
}

.subtitle {
  font-size: var(--font-size-lg);
  color: var(--medium);
}

.form-container {
  max-width: 800px;
  margin: 0 auto;
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}

/* Form stepper */
.form-stepper {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--space-6);
  padding: var(--space-4) var(--space-6);
  background-color: var(--light);
  position: relative;
}

.form-stepper::before {
  content: "";
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 2px;
  background-color: var(--medium);
  transform: translateY(-50%);
  z-index: 1;
}

.stepper-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  z-index: 2;
  cursor: pointer;
}

.step-counter {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: var(--white);
  border: 2px solid var(--medium);
  color: var(--medium);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-2);
  transition: all var(--transition-normal);
}

.step-name {
  font-size: var(--font-size-sm);
  color: var(--medium);
  transition: color var(--transition-normal);
}

.stepper-item.active .step-counter {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
}

.stepper-item.active .step-name {
  color: var(--primary-color);
  font-weight: var(--font-weight-semibold);
}

.stepper-item.completed .step-counter {
  background-color: var(--success-color);
  border-color: var(--success-color);
  color: white;
}

.stepper-item.completed .step-name {
  color: var(--success-color);
}

/* Job form */
.job-form {
  padding: var(--space-6);
}

.form-step {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.form-step h2 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-2);
  border-bottom: 1px solid var(--light);
  color: var(--dark);
}

.form-group {
  margin-bottom: var(--space-4);
}

.form-group label {
  display: block;
  margin-bottom: var(--space-2);
  font-weight: var(--font-weight-medium);
  color: var(--dark);
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  font-size: var(--font-size-md);
  transition: border-color var(--transition-fast);
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: var(--primary-color);
  outline: none;
}

.form-group input.input-error,
.form-group select.input-error,
.form-group textarea.input-error,
.tags-input-container.input-error {
  border-color: var(--error-color);
}

.error-message {
  color: var(--error-color);
  font-size: var(--font-size-sm);
  margin-top: var(--space-1);
}

.textarea-footer {
  display: flex;
  justify-content: space-between;
  margin-top: var(--space-1);
}

.char-count {
  font-size: var(--font-size-xs);
  color: var(--medium);
  text-align: right;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
}

/* Tags input for skills */
.tags-input-container {
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  padding: var(--space-2);
  min-height: 50px;
  transition: border-color var(--transition-fast);
}

.tags-input-container:focus-within {
  border-color: var(--primary-color);
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.tag {
  background-color: var(--primary-light);
  color: var(--primary-color);
  border-radius: var(--radius-full);
  padding: var(--space-1) var(--space-3);
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
}

.remove-tag {
  background: none;
  border: none;
  color: var(--primary-color);
  font-size: var(--font-size-lg);
  margin-left: var(--space-1);
  cursor: pointer;
  line-height: 1;
}

.tags-input {
  flex: 1;
  min-width: 100px;
  border: none;
  padding: var(--space-1);
  font-size: var(--font-size-sm);
  background: transparent;
}

.tags-input:focus {
  outline: none;
}

.helper-text {
  font-size: var(--font-size-xs);
  color: var(--medium);
  margin-top: var(--space-1);
}

/* Checkbox styling */
.checkbox-group {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.checkbox-group input[type="checkbox"] {
  width: 18px;
  height: 18px;
}

.checkbox-group label {
  margin-bottom: 0;
}

/* Preview card */
.form-preview {
  margin-top: var(--space-6);
  border-top: 1px solid var(--light);
  padding-top: var(--space-6);
}

.form-preview h3 {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-3);
  color: var(--dark);
}

.preview-card {
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.preview-header {
  background-color: var(--light);
  padding: var(--space-3);
  border-bottom: 1px solid var(--light);
}

.preview-header h4 {
  margin: 0 0 var(--space-1) 0;
  font-size: var(--font-size-md);
  color: var(--dark);
}

.preview-meta {
  display: flex;
  gap: var(--space-3);
  font-size: var(--font-size-sm);
  color: var(--medium);
}

.preview-body {
  padding: var(--space-3);
}

.preview-body p {
  margin: 0 0 var(--space-3) 0;
  color: var(--medium);
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.preview-skills {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.preview-skill {
  background-color: var(--light);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  color: var(--dark);
}

.preview-skill.empty {
  font-style: italic;
  color: var(--medium);
}

.preview-footer {
  display: flex;
  justify-content: space-between;
  padding: var(--space-3);
  background-color: var(--light);
  border-top: 1px solid var(--light);
  font-size: var(--font-size-sm);
  color: var(--medium);
}

/* Form navigation */
.form-navigation {
  display: flex;
  justify-content: space-between;
  margin-top: var(--space-6);
  padding-top: var(--space-6);
  border-top: 1px solid var(--light);
}

.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 120px;
}

.spinner {
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
    gap: var(--space-3);
  }

  .step-name {
    display: none;
  }
}

@media (max-width: 576px) {
  .form-navigation {
    flex-direction: column;
    gap: var(--space-3);
  }

  .form-navigation button {
    width: 100%;
  }
}
</style>
