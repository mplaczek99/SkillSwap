<template>
  <div class="auth-page">
    <div class="container">
      <div class="auth-container">
        <div class="auth-card">
          <div class="auth-header">
            <h1>Create Account</h1>
            <p>Join SkillSwap to start sharing and learning new skills!</p>
          </div>

          <form @submit.prevent="submitRegister" class="auth-form">
            <div class="form-group">
              <label for="name" class="form-label">Full Name</label>
              <div class="input-with-icon">
                <input
                  id="name"
                  type="text"
                  v-model="name"
                  placeholder="Your full name"
                  class="form-control"
                  :class="{ 'input-error': errors.name }"
                  required
                  autocomplete="name"
                />
                <font-awesome-icon icon="user" class="input-icon" />
              </div>
              <span v-if="errors.name" class="error-message">{{
                errors.name
              }}</span>
            </div>

            <div class="form-group">
              <label for="email" class="form-label">Email</label>
              <div class="input-with-icon">
                <input
                  id="email"
                  type="email"
                  v-model="email"
                  placeholder="Your email address"
                  class="form-control"
                  :class="{ 'input-error': errors.email }"
                  required
                  autocomplete="email"
                />
                <font-awesome-icon icon="envelope" class="input-icon" />
              </div>
              <span v-if="errors.email" class="error-message">{{
                errors.email
              }}</span>
            </div>

            <div class="form-group">
              <label for="password" class="form-label">Password</label>
              <div class="input-with-icon">
                <input
                  id="password"
                  :type="showPassword ? 'text' : 'password'"
                  v-model="password"
                  placeholder="Create a password"
                  class="form-control"
                  :class="{ 'input-error': errors.password }"
                  required
                  autocomplete="new-password"
                />
                <font-awesome-icon
                  :icon="showPassword ? 'eye-slash' : 'eye'"
                  class="input-icon clickable"
                  @click="togglePassword"
                />
              </div>
              <span v-if="errors.password" class="error-message">{{
                errors.password
              }}</span>
              <div class="password-strength" v-if="password">
                <div class="strength-meter">
                  <div
                    class="strength-meter-fill"
                    :style="{ width: `${passwordStrength.score * 25}%` }"
                    :class="passwordStrength.className"
                  ></div>
                </div>
                <span class="strength-text" :class="passwordStrength.className">
                  {{ passwordStrength.text }}
                </span>
              </div>
            </div>

            <div class="form-group terms-checkbox">
              <label class="checkbox-container">
                <input type="checkbox" v-model="agreeToTerms" required />
                <span class="checkmark"></span>
                I agree to the
                <a href="#" @click.prevent="showTerms">Terms of Service</a>
                and
                <a href="#" @click.prevent="showPrivacy">Privacy Policy</a>
              </label>
              <span v-if="errors.terms" class="error-message">{{
                errors.terms
              }}</span>
            </div>

            <div class="form-group">
              <button
                type="submit"
                class="btn btn-primary btn-full"
                :disabled="isLoading"
              >
                <span v-if="isLoading" class="spinner"></span>
                <span v-else>Sign Up</span>
              </button>
            </div>

            <div v-if="error" class="alert alert-danger">
              {{ error }}
            </div>
          </form>

          <div class="auth-footer">
            <p>
              Already have an account?
              <router-link to="/login">Sign in</router-link>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "RegisterForm",
  data() {
    return {
      name: "",
      email: "",
      password: "",
      agreeToTerms: false,
      showPassword: false,
      error: null,
      isLoading: false,
      errors: {
        name: null,
        email: null,
        password: null,
        terms: null,
      },
    };
  },
  computed: {
    passwordStrength() {
      if (!this.password) {
        return { score: 0, text: "", className: "" };
      }

      // Simple password strength calculation
      let score = 0;

      // Length check
      if (this.password.length >= 8) score++;
      if (this.password.length >= 12) score++;

      // Character variety checks
      if (/[A-Z]/.test(this.password)) score++;
      if (/[0-9]/.test(this.password)) score++;
      if (/[^A-Za-z0-9]/.test(this.password)) score++;

      let text = "";
      let className = "";

      switch (score) {
        case 0:
        case 1:
          text = "Weak";
          className = "weak";
          break;
        case 2:
          text = "Fair";
          className = "fair";
          break;
        case 3:
          text = "Good";
          className = "good";
          break;
        case 4:
        case 5:
          text = "Strong";
          className = "strong";
          break;
      }

      return { score, text, className };
    },
  },
  methods: {
    togglePassword() {
      this.showPassword = !this.showPassword;
    },
    validateForm() {
      let isValid = true;
      this.errors = {
        name: null,
        email: null,
        password: null,
        terms: null,
      };

      // Name validation
      if (!this.name.trim()) {
        this.errors.name = "Name is required";
        isValid = false;
      }

      // Email validation
      if (!this.email) {
        this.errors.email = "Email is required";
        isValid = false;
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.email)) {
        this.errors.email = "Please enter a valid email address";
        isValid = false;
      }

      // Password validation
      if (!this.password) {
        this.errors.password = "Password is required";
        isValid = false;
      } else if (this.password.length < 6) {
        this.errors.password = "Password must be at least 6 characters";
        isValid = false;
      }

      // Terms validation
      if (!this.agreeToTerms) {
        this.errors.terms =
          "You must agree to the Terms of Service and Privacy Policy";
        isValid = false;
      }

      return isValid;
    },
    async submitRegister() {
      if (!this.validateForm()) return;

      this.isLoading = true;
      this.error = null;

      try {
        await this.$store.dispatch("register", {
          name: this.name,
          email: this.email,
          password: this.password,
        });

        // Redirect to home page
        this.$router.push("/");
      } catch (err) {
        this.error =
          err.response?.data?.error || "Registration failed. Please try again.";
      } finally {
        this.isLoading = false;
      }
    },
    showTerms() {
      // In a real app, show terms of service modal or navigate to terms page
      alert("Terms of Service would be shown here");
    },
    showPrivacy() {
      // In a real app, show privacy policy modal or navigate to privacy page
      alert("Privacy Policy would be shown here");
    },
  },
};
</script>

<style scoped>
.auth-page {
  min-height: calc(100vh - 4rem - 108px);
  display: flex;
  align-items: center;
  padding: var(--space-8) 0;
}

.auth-container {
  max-width: 450px;
  margin: 0 auto;
  width: 100%;
}

.auth-card {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}

.auth-header {
  padding: var(--space-6);
  text-align: center;
  background-color: var(--primary-light);
}

.auth-header h1 {
  font-size: var(--font-size-2xl);
  color: var(--primary-color);
  margin-bottom: var(--space-2);
}

.auth-header p {
  color: var(--medium);
  margin-bottom: 0;
}

.auth-form {
  padding: var(--space-6);
}

.input-with-icon {
  position: relative;
}

.input-icon {
  position: absolute;
  right: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--medium);
}

.input-icon.clickable {
  cursor: pointer;
}

.input-error {
  border-color: var(--error-color) !important;
}

.error-message {
  color: var(--error-color);
  font-size: var(--font-size-sm);
  display: block;
  margin-top: var(--space-1);
}

/* Password strength meter */
.password-strength {
  margin-top: var(--space-2);
}

.strength-meter {
  height: 4px;
  background-color: #e0e0e0;
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-bottom: var(--space-1);
}

.strength-meter-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width var(--transition-normal) ease;
}

.strength-meter-fill.weak {
  background-color: var(--error-color);
}

.strength-meter-fill.fair {
  background-color: var(--warning-color);
}

.strength-meter-fill.good {
  background-color: var(--info-color);
}

.strength-meter-fill.strong {
  background-color: var(--success-color);
}

.strength-text {
  font-size: var(--font-size-xs);
  float: right;
}

.strength-text.weak {
  color: var(--error-color);
}

.strength-text.fair {
  color: var(--warning-color);
}

.strength-text.good {
  color: var(--info-color);
}

.strength-text.strong {
  color: var(--success-color);
}

/* Terms checkbox */
.terms-checkbox {
  margin-top: var(--space-4);
}

.checkbox-container {
  display: flex;
  align-items: center;
  position: relative;
  padding-left: 30px;
  cursor: pointer;
  user-select: none;
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.checkbox-container input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

.checkmark {
  position: absolute;
  left: 0;
  top: 0;
  height: 20px;
  width: 20px;
  background-color: var(--white);
  border: 1px solid var(--medium);
  border-radius: var(--radius-sm);
}

.checkbox-container:hover input ~ .checkmark {
  border-color: var(--primary-color);
}

.checkbox-container input:checked ~ .checkmark {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.checkmark:after {
  content: "";
  position: absolute;
  display: none;
}

.checkbox-container input:checked ~ .checkmark:after {
  display: block;
}

.checkbox-container .checkmark:after {
  left: 7px;
  top: 3px;
  width: 6px;
  height: 12px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.auth-footer {
  padding: var(--space-4);
  text-align: center;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  background-color: var(--light);
}

@media (max-width: 576px) {
  .auth-container {
    padding: 0 var(--space-4);
  }

  .auth-header,
  .auth-form {
    padding: var(--space-4);
  }
}
</style>
