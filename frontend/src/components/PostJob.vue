<template>
  <div class="post-job-page" :class="{ 'dark-theme': isDark }">
    <div class="container">
      <div class="post-job-header">
        <div class="header-content">
          <h1>{{ isEditing ? "Edit Job Posting" : "Post a New Job" }}</h1>
          <p class="subtitle">
            {{
              isEditing
                ? "Update the details of your job posting"
                : "Share an opportunity with the SkillSwap community"
            }}
          </p>
        </div>
        <div class="theme-toggle" @click="toggleDark">
          <font-awesome-icon :icon="isDark ? 'sun' : 'moon'" />
        </div>
      </div>

      <div class="form-container">
        <!-- Form stepper -->
        <div class="form-stepper">
          <div class="stepper-progress" :style="{ width: progressWidth }"></div>
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
          <transition name="fade" mode="out-in">
            <div v-if="currentStep === 0" class="form-step">
              <h2>Basic Information</h2>
              <p class="step-description">
                Let's start with the essential details about this job opportunity.
              </p>

              <div class="form-group">
                <label for="job-title">Job Title <span class="required">*</span></label>
                <div class="input-wrapper">
                  <font-awesome-icon icon="briefcase" class="input-icon" />
                  <input
                    id="job-title"
                    v-model="formData.title"
                    type="text"
                    placeholder="e.g. Frontend Developer"
                    required
                    :class="{ 'input-error': errors.title }"
                  />
                </div>
                <span v-if="errors.title" class="error-message">
                  <font-awesome-icon icon="exclamation-circle" />
                  {{ errors.title }}
                </span>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label for="job-company">Company Name <span class="required">*</span></label>
                  <div class="input-wrapper">
                    <font-awesome-icon icon="building" class="input-icon" />
                    <input
                      id="job-company"
                      v-model="formData.company"
                      type="text"
                      placeholder="e.g. Tech Innovators"
                      required
                      :class="{ 'input-error': errors.company }"
                    />
                  </div>
                  <span v-if="errors.company" class="error-message">
                    <font-awesome-icon icon="exclamation-circle" />
                    {{ errors.company }}
                  </span>
                </div>

                <div class="form-group">
                  <label for="job-location">Location <span class="required">*</span></label>
                  <div class="input-wrapper">
                    <font-awesome-icon icon="map-marker-alt" class="input-icon" />
                    <input
                      id="job-location"
                      v-model="formData.location"
                      type="text"
                      placeholder="e.g. San Francisco, CA or Remote"
                      required
                      :class="{ 'input-error': errors.location }"
                    />
                  </div>
                  <span v-if="errors.location" class="error-message">
                    <font-awesome-icon icon="exclamation-circle" />
                    {{ errors.location }}
                  </span>
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label for="job-type">Job Type <span class="required">*</span></label>
                  <div class="select-wrapper">
                    <font-awesome-icon icon="clock" class="input-icon" />
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
                    <font-awesome-icon icon="chevron-down" class="select-arrow" />
                  </div>
                  <span v-if="errors.jobType" class="error-message">
                    <font-awesome-icon icon="exclamation-circle" />
                    {{ errors.jobType }}
                  </span>
                </div>

                <div class="form-group">
                  <label for="job-experience">Experience Level <span class="required">*</span></label>
                  <div class="select-wrapper">
                    <font-awesome-icon icon="user-graduate" class="input-icon" />
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
                    <font-awesome-icon icon="chevron-down" class="select-arrow" />
                  </div>
                  <span v-if="errors.experienceLevel" class="error-message">
                    <font-awesome-icon icon="exclamation-circle" />
                    {{ errors.experienceLevel }}
                  </span>
                </div>
              </div>

              <div class="form-group">
                <label for="job-salary">Salary Range</label>
                <div class="input-wrapper">
                  <font-awesome-icon icon="dollar-sign" class="input-icon" />
                  <input
                    id="job-salary"
                    v-model="formData.salaryRange"
                    type="text"
                    placeholder="e.g. $50,000 - $70,000 or $25 - $35 per hour"
                  />
                </div>
              </div>
            </div>
          </transition>

          <!-- Step 2: Job Description -->
          <transition name="fade" mode="out-in">
            <div v-if="currentStep === 1" class="form-step">
              <h2>Job Description</h2>
              <p class="step-description">
                Provide details about the role, responsibilities, and required skills.
              </p>

              <div class="form-group">
                <label for="job-description">Detailed Description <span class="required">*</span></label>
                <div class="textarea-wrapper">
                  <textarea
                    id="job-description"
                    v-model="formData.description"
                    rows="8"
                    placeholder="Provide a comprehensive description of the job role, responsibilities, and requirements..."
                    required
                    :class="{ 'input-error': errors.description }"
                  ></textarea>
                </div>
                <div class="textarea-footer">
                  <span v-if="errors.description" class="error-message">
                    <font-awesome-icon icon="exclamation-circle" />
                    {{ errors.description }}
                  </span>
                  <span class="char-count" :class="{ 'warning': formData.description.length > 1800 }">
                    {{ formData.description.length }}/2000
                  </span>
                </div>
              </div>

              <div class="form-group">
                <label for="job-skills">Required Skills <span class="required">*</span></label>
                <div
                  class="tags-input-container"
                  :class="{ 'input-error': errors.skillsRequired }"
                >
                  <font-awesome-icon icon="tags" class="tags-icon" />
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
                <span v-if="errors.skillsRequired" class="error-message">
                  <font-awesome-icon icon="exclamation-circle" />
                  {{ errors.skillsRequired }}
                </span>
                <small class="helper-text">
                  <font-awesome-icon icon="info-circle" />
                  Press Enter or comma after each skill
                </small>
              </div>
              
              <div class="skill-suggestions" v-if="newSkill.length > 0">
                <p class="suggestion-title">Suggestions:</p>
                <div class="suggestion-tags">
                  <button 
                    v-for="(suggestion, index) in filteredSuggestions" 
                    :key="index"
                    type="button" 
                    class="suggestion-tag"
                    @click="selectSuggestion(suggestion)"
                  >
                    {{ suggestion }}
                  </button>
                </div>
              </div>
            </div>
          </transition>

          <!-- Step 3: Contact Information -->
          <transition name="fade" mode="out-in">
            <div v-if="currentStep === 2" class="form-step">
              <h2>Contact Information</h2>
              <p class="step-description">
                Add contact details and review your job posting before submitting.
              </p>

              <div class="form-group">
                <label for="job-email">Contact Email <span class="required">*</span></label>
                <div class="input-wrapper">
                  <font-awesome-icon icon="envelope" class="input-icon" />
                  <input
                    id="job-email"
                    v-model="formData.contactEmail"
                    type="email"
                    placeholder="e.g. hiring@company.com"
                    required
                    :class="{ 'input-error': errors.contactEmail }"
                  />
                </div>
                <span v-if="errors.contactEmail" class="error-message">
                  <font-awesome-icon icon="exclamation-circle" />
                  {{ errors.contactEmail }}
                </span>
              </div>

              <div class="form-row">
                <div class="form-group checkbox-group">
                  <input
                    type="checkbox"
                    id="job-notify"
                    v-model="notifyOnApplications"
                  />
                  <label for="job-notify">
                    <span class="checkbox-custom"></span>
                    Notify me when someone applies
                  </label>
                </div>
              </div>

              <div class="form-preview">
                <h3>
                  <font-awesome-icon icon="eye" />
                  Job Posting Preview
                </h3>
                <div class="preview-card">
                  <div class="preview-header">
                    <h4>{{ formData.title || "Job Title" }}</h4>
                    <div class="preview-meta">
                      <span>
                        <font-awesome-icon icon="building" />
                        {{ formData.company || "Company Name" }}
                      </span>
                      <span>
                        <font-awesome-icon icon="map-marker-alt" />
                        {{ formData.location || "Location" }}
                      </span>
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
                    <span>
                      <font-awesome-icon icon="clock" />
                      {{ formData.jobType || "Job Type" }}
                    </span>
                    <span>
                      <font-awesome-icon icon="user-graduate" />
                      {{
                        formData.experienceLevel || "Experience Level"
                      }}
                      Level
                    </span>
                    <span v-if="formData.salaryRange">
                      <font-awesome-icon icon="dollar-sign" />
                      {{ formData.salaryRange }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </transition>

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
            <div v-else class="btn-placeholder"></div>

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
              <span v-else>
                <font-awesome-icon :icon="isEditing ? 'edit' : 'paper-plane'" />
                {{ isEditing ? "Update Job" : "Post Job" }}
              </span>
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
import { useDark, useToggle } from '@vueuse/core';

export default {
  name: "PostJob",
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
      skillSuggestions: [
        "JavaScript", "TypeScript", "React", "Vue.js", "Angular", "Node.js", 
        "Python", "Java", "C#", "PHP", "Ruby", "Go", "Swift", "Kotlin",
        "HTML", "CSS", "SASS", "LESS", "Tailwind CSS", "Bootstrap",
        "SQL", "MongoDB", "PostgreSQL", "MySQL", "Firebase", "AWS", "Azure",
        "Docker", "Kubernetes", "Git", "CI/CD", "DevOps", "Agile", "Scrum",
        "UI/UX Design", "Figma", "Adobe XD", "Photoshop", "Illustrator",
        "Project Management", "Communication", "Leadership", "Problem Solving"
      ]
    };
  },
  computed: {
    progressWidth() {
      const percentage = ((this.currentStep + 1) / this.steps.length) * 100;
      return `${percentage}%`;
    },
    filteredSuggestions() {
      if (!this.newSkill) return [];
      
      const input = this.newSkill.toLowerCase();
      const existingSkills = this.formData.skillsRequired.map(s => s.toLowerCase());
      
      return this.skillSuggestions
        .filter(skill => 
          skill.toLowerCase().includes(input) && 
          !existingSkills.includes(skill.toLowerCase())
        )
        .slice(0, 5); // Limit to 5 suggestions
    }
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
        window.scrollTo({ top: 0, behavior: 'smooth' });
      } else {
        // Shake the form to indicate errors
        const formEl = document.querySelector('.form-step');
        if (formEl) {
          formEl.classList.add('shake');
          setTimeout(() => {
            formEl.classList.remove('shake');
          }, 500);
        }
      }
    },

    prevStep() {
      this.currentStep -= 1;
      window.scrollTo({ top: 0, behavior: 'smooth' });
    },

    goToStep(step) {
      if (step < this.currentStep) {
        this.currentStep = step;
        window.scrollTo({ top: 0, behavior: 'smooth' });
      } else if (step > this.currentStep) {
        // Only allow jumping ahead if current step is valid
        if (this.validateStep()) {
          this.currentStep = step;
          window.scrollTo({ top: 0, behavior: 'smooth' });
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

    selectSuggestion(skill) {
      if (!this.formData.skillsRequired.includes(skill)) {
        this.formData.skillsRequired.push(skill);
      }
      this.newSkill = "";
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
        // Use a promise-based approach to handle errors in the timeout
        await new Promise((resolve, reject) => {
          setTimeout(() => {
            try {
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
              resolve();
            } catch (err) {
              reject(err);
            }
          }, 1500);
        });
      } catch (error) {
        console.error("Error submitting job:", error);

        // Show error message
        eventBus.emit("show-notification", {
          type: "error",
          title: "Error",
          message:
            "There was a problem submitting your job posting. Please try again.",
          duration: 5000,
        });
      } finally {
        this.isSubmitting = false;
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
/* Base Styles */
.post-job-page {
  padding: 2rem 0 4rem;
  min-height: 100vh;
  background-color: #f8fafc;
  transition: background-color 0.3s ease, color 0.3s ease;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.post-job-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2.5rem;
}

.header-content {
  text-align: center;
  flex: 1;
}

.post-job-header h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  color: #1e293b;
  font-weight: 800;
  background: linear-gradient(135deg, #4f46e5, #3a0ca3);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  transition: all 0.3s ease;
}

.subtitle {
  font-size: 1.125rem;
  color: #64748b;
  max-width: 600px;
  margin: 0 auto;
  transition: color 0.3s ease;
}

/* Theme Toggle */
.theme-toggle {
  position: relative;
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

/* Dark Theme */
.dark-theme {
  background-color: #0f172a;
  color: #e2e8f0;
}

.dark-theme .post-job-header h1 {
  background: linear-gradient(135deg, #818cf8, #6366f1);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.dark-theme .subtitle {
  color: #94a3b8;
}

.dark-theme .form-container {
  background-color: #1e293b;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  border-color: #334155;
}

.dark-theme .form-stepper {
  background-color: #0f172a;
}

.dark-theme .form-stepper::before {
  background-color: #334155;
}

.dark-theme .step-counter {
  background-color: #1e293b;
  border-color: #475569;
  color: #94a3b8;
}

.dark-theme .step-name {
  color: #94a3b8;
}

.dark-theme .form-step h2 {
  color: #e2e8f0;
  border-bottom-color: #334155;
}

.dark-theme .step-description {
  color: #94a3b8;
}

.dark-theme .form-group label {
  color: #e2e8f0;
}

.dark-theme .input-wrapper,
.dark-theme .select-wrapper,
.dark-theme .textarea-wrapper,
.dark-theme .tags-input-container {
  background-color: #334155;
  border-color: #475569;
}

.dark-theme .input-wrapper input,
.dark-theme .select-wrapper select,
.dark-theme .textarea-wrapper textarea,
.dark-theme .tags-input {
  color: #e2e8f0;
  background-color: transparent;
}

.dark-theme .input-icon,
.dark-theme .select-arrow,
.dark-theme .tags-icon {
  color: #94a3b8;
}

.dark-theme .tag {
  background-color: #3e4c6b;
  color: #a5b4fc;
}

.dark-theme .preview-card {
  background-color: #1e293b;
  border-color: #334155;
}

.dark-theme .preview-header,
.dark-theme .preview-footer {
  background-color: #0f172a;
  border-color: #334155;
}

.dark-theme .preview-skill {
  background-color: #334155;
  color: #e2e8f0;
}

.dark-theme .checkbox-custom {
  border-color: #475569;
}

.dark-theme .suggestion-tag {
  background-color: #334155;
  color: #e2e8f0;
}

.dark-theme .suggestion-tag:hover {
  background-color: #4f46e5;
}

/* Form Container */
.form-container {
  max-width: 800px;
  margin: 0 auto;
  background-color: white;
  border-radius: 1rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -5px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

/* Form stepper */
.form-stepper {
  display: flex;
  justify-content: space-between;
  padding: 1.5rem 2rem;
  background-color: #f1f5f9;
  position: relative;
  transition: background-color 0.3s ease;
}

.form-stepper::before {
  content: "";
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 2px;
  background-color: #cbd5e1;
  transform: translateY(-50%);
  z-index: 1;
}

.stepper-progress {
  position: absolute;
  top: 50%;
  left: 0;
  height: 2px;
  background: linear-gradient(90deg, #4f46e5, #818cf8);
  transform: translateY(-50%);
  z-index: 2;
  transition: width 0.5s ease;
}

.stepper-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  z-index: 3;
  cursor: pointer;
  transition: all 0.3s ease;
}

.step-counter {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: white;
  border: 2px solid #94a3b8;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  margin-bottom: 0.5rem;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.step-counter::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #4f46e5, #818cf8);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
}

.stepper-item.active .step-counter {
  background-color: #4f46e5;
  border-color: #4f46e5;
  color: white;
  transform: scale(1.1);
  box-shadow: 0 0 0 5px rgba(79, 70, 229, 0.2);
}

.stepper-item.completed .step-counter {
  background-color: #10b981;
  border-color: #10b981;
  color: white;
}

.step-name {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.stepper-item.active .step-name {
  color: #4f46e5;
  font-weight: 600;
  transform: scale(1.05);
}

.stepper-item.completed .step-name {
  color: #10b981;
}

/* Job form */
.job-form {
  padding: 2rem;
}

.form-step {
  animation: fadeIn 0.5s ease;
}

.step-description {
  color: #64748b;
  margin-bottom: 1.5rem;
  font-size: 1rem;
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

.form-step h2 {
  font-size: 1.5rem;
  margin-bottom: 0.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
  color: #1e293b;
  font-weight: 700;
  transition: all 0.3s ease;
}

.form-group {
  margin-bottom: 1.5rem;
  position: relative;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #334155;
  transition: color 0.3s ease;
}

.required {
  color: #ef4444;
  margin-left: 2px;
}

.input-wrapper,
.select-wrapper,
.textarea-wrapper {
  position: relative;
  transition: all 0.3s ease;
}

.input-icon,
.select-arrow,
.tags-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #64748b;
  transition: color 0.3s ease;
}

.select-arrow {
  left: auto;
  right: 1rem;
  pointer-events: none;
}

.input-wrapper input,
.select-wrapper select,
.textarea-wrapper textarea {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.3s ease;
  background-color: white;
}

.textarea-wrapper textarea {
  padding-left: 1rem;
  resize: vertical;
  min-height: 120px;
}

.input-wrapper input:focus,
.select-wrapper select:focus,
.textarea-wrapper textarea:focus {
  border-color: #4f46e5;
  outline: none;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2);
}

.input-wrapper input.input-error,
.select-wrapper select.input-error,
.textarea-wrapper textarea.input-error,
.tags-input-container.input-error {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.error-message {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  animation: shake 0.5s ease;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
  20%, 40%, 60%, 80% { transform: translateX(5px); }
}

.textarea-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 0.5rem;
}

.char-count {
  font-size: 0.75rem;
  color: #64748b;
  text-align: right;
}

.char-count.warning {
  color: #f59e0b;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

/* Tags input for skills */
.tags-input-container {
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  padding: 0.5rem 0.75rem 0.5rem 2.5rem;
  min-height: 50px;
  transition: all 0.3s ease;
  position: relative;
  background-color: white;
}

.tags-icon {
  top: 1rem;
}

.tags-input-container:focus-within {
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2);
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  background-color: #eef2ff;
  color: #4f46e5;
  border-radius: 9999px;
  padding: 0.25rem 0.75rem;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
  animation: tagAppear 0.3s ease;
}

@keyframes tagAppear {
  from {
    opacity: 0;
    transform: scale(0.8);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.tag:hover {
  background-color: #e0e7ff;
  transform: translateY(-2px);
}

.remove-tag {
  background: none;
  border: none;
  color: #4f46e5;
  font-size: 1.25rem;
  margin-left: 0.25rem;
  cursor: pointer;
  line-height: 1;
  transition: all 0.2s ease;
}

.remove-tag:hover {
  color: #ef4444;
  transform: scale(1.2);
}

.tags-input {
  flex: 1;
  min-width: 100px;
  border: none;
  padding: 0.25rem;
  font-size: 0.875rem;
  background: transparent;
}

.tags-input:focus {
  outline: none;
}

.helper-text {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* Skill suggestions */
.skill-suggestions {
  margin-top: 0.5rem;
  animation: fadeIn 0.3s ease;
}

.suggestion-title {
  font-size: 0.75rem;
  color: #64748b;
  margin-bottom: 0.5rem;
}

.suggestion-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.suggestion-tag {
  background-color: #f1f5f9;
  color: #334155;
  border: none;
  border-radius: 9999px;
  padding: 0.25rem 0.75rem;
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.suggestion-tag:hover {
  background-color: #4f46e5;
  color: white;
  transform: translateY(-2px);
}

/* Checkbox styling */
.checkbox-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  position: relative;
}

.checkbox-group input[type="checkbox"] {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  margin-bottom: 0;
}

.checkbox-custom {
  width: 20px;
  height: 20px;
  border: 2px solid #cbd5e1;
  border-radius: 4px;
  display: inline-block;
  position: relative;
  transition: all 0.2s ease;
}

.checkbox-group input[type="checkbox"]:checked + label .checkbox-custom {
  background-color: #4f46e5;
  border-color: #4f46e5;
}

.checkbox-group input[type="checkbox"]:checked + label .checkbox-custom::after {
  content: '';
  position: absolute;
  left: 6px;
  top: 2px;
  width: 5px;
  height: 10px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.checkbox-group input[type="checkbox"]:focus + label .checkbox-custom {
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.2);
}

/* Preview card */
.form-preview {
  margin-top: 2rem;
  border-top: 1px solid #e2e8f0;
  padding-top: 2rem;
}

.form-preview h3 {
  font-size: 1.25rem;
  margin-bottom: 1rem;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.preview-card {
  border: 1px solid #e2e8f0;
  border-radius: 0.75rem;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.preview-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

.preview-header {
  background-color: #f8fafc;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  transition: background-color 0.3s ease;
}

.preview-header h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
  color: #1e293b;
  font-weight: 600;
}

.preview-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  font-size: 0.875rem;
  color: #64748b;
}

.preview-meta span {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.preview-body {
  padding: 1rem;
}

.preview-body p {
  margin: 0 0 1rem 0;
  color: #475569;
  font-size: 0.875rem;
  line-height: 1.5;
}

.preview-skills {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.preview-skill {
  background-color: #f1f5f9;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  color: #334155;
}

.preview-skill.empty {
  font-style: italic;
  color: #94a3b8;
}

.preview-footer {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  padding: 1rem;
  background-color: #f8fafc;
  border-top: 1px solid #e2e8f0;
  font-size: 0.875rem;
  color: #64748b;
  gap: 0.5rem;
}

.preview-footer span {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* Form navigation */
.form-navigation {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e2e8f0;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  transition: all 0.3s ease;
  cursor: pointer;
  text-decoration: none;
  padding: 0.75rem 1.5rem;
  font-size: 0.875rem;
}

.btn-lg {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
}

.btn-primary {
  background-color: #4f46e5;
  color: white;
  border: 1px solid #4f46e5;
}

.btn-primary:hover {
  background-color: #4338ca;
  border-color: #4338ca;
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.btn-outline {
  background-color: transparent;
  color: #4f46e5;
  border: 1px solid #4f46e5;
}

.btn-outline:hover {
  background-color: #eef2ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.btn-placeholder {
  width: 120px;
}

.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 150px;
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.spinner {
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  width: 20px;
  height: 20px;
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

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

/* Shake animation for validation errors */
.shake {
  animation: shake 0.5s ease;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .post-job-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .theme-toggle {
    position: absolute;
    top: 1rem;
    right: 1rem;
  }
  
  .form-row {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .step-name {
    display: none;
  }
  
  .preview-footer {
    flex-direction: column;
    align-items: flex-start;
  }
}

@media (max-width: 576px) {
  .post-job-header h1 {
    font-size: 1.75rem;
  }
  
  .subtitle {
    font-size: 1rem;
  }
  
  .form-navigation {
    flex-direction: column;
    gap: 1rem;
  }

  .form-navigation button {
    width: 100%;
  }
  
  .btn-placeholder {
    display: none;
  }
  
  .job-form {
    padding: 1.5rem 1rem;
  }
}
</style>