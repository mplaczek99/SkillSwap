<template>
  <div class="auth-page">
    <div class="container">
      <div class="auth-container">
        <div class="auth-card">
          <div class="auth-header">
            <h1>Login</h1>
            <p>Welcome back! Please enter your credentials.</p>
          </div>

          <form @submit.prevent="submitLogin" class="auth-form">
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
              <div class="label-with-link">
                <label for="password" class="form-label">Password</label>
                <a href="#" class="forgot-password">Forgot password?</a>
              </div>
              <div class="input-with-icon">
                <input
                  id="password"
                  :type="showPassword ? 'text' : 'password'"
                  v-model="password"
                  placeholder="Your password"
                  class="form-control"
                  :class="{ 'input-error': errors.password }"
                  required
                  autocomplete="current-password"
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
            </div>

            <div class="form-group remember-me">
              <label class="checkbox-container">
                <input type="checkbox" v-model="rememberMe" />
                <span class="checkmark"></span>
                Remember me
              </label>
            </div>

            <div class="form-group">
              <button
                type="submit"
                class="btn btn-primary btn-full"
                :disabled="isLoading"
              >
                <span v-if="isLoading" class="spinner"></span>
                <span v-else>Sign In</span>
              </button>
            </div>

            <div v-if="error" class="alert alert-danger">
              {{ error }}
            </div>
          </form>

          <div class="auth-footer">
            <p>
              Don't have an account?
              <router-link to="/register">Sign up</router-link>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "LoginForm",
  data() {
    return {
      email: "",
      password: "",
      rememberMe: false,
      showPassword: false,
      error: null,
      isLoading: false,
      errors: {
        email: null,
        password: null,
      },
    };
  },
  methods: {
    togglePassword() {
      this.showPassword = !this.showPassword;
    },
    validateForm() {
      let isValid = true;
      this.errors = {
        email: null,
        password: null,
      };

      // Simple email validation
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
      }

      return isValid;
    },
    async submitLogin() {
      if (!this.validateForm()) return;

      this.isLoading = true;
      this.error = null;

      try {
        await this.$store.dispatch("login", {
          email: this.email,
          password: this.password,
          rememberMe: this.rememberMe, // Pass this value to the login action
        });

        // Redirect to home page or intended destination
        this.$router.push("/");
      } catch (err) {
        this.error =
          err.response?.data?.error ||
          "Login failed. Please check your credentials.";
      } finally {
        this.isLoading = false;
      }
    },
  },
};
</script>

<style scoped>
.auth-page {
  min-height: calc(100vh - 4rem - 108px);
  /* Adjust for navbar and footer */
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

.label-with-link {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-2);
}

.forgot-password {
  font-size: var(--font-size-sm);
}

.remember-me {
  display: flex;
  align-items: center;
}

.checkbox-container {
  display: flex;
  align-items: center;
  position: relative;
  padding-left: 30px;
  cursor: pointer;
  user-select: none;
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

.input-error {
  border-color: var(--error-color) !important;
}

.error-message {
  color: var(--error-color);
  font-size: var(--font-size-sm);
  display: block;
  margin-top: var(--space-1);
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
