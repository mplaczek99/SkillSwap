import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import eventBus from "./utils/eventBus"; // Import the eventBus

// Import the icon registration function
import { registerIcons } from "./utils/icons";

// Add default timeout to all axios requests
axios.defaults.timeout = 15000; // 15 seconds

// Set up axios interceptors to add the auth token to all requests, with CSRF protection
axios.interceptors.request.use(
  (config) => {
    const token = store.state.token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      // Add CSRF token if available
      const csrfToken = document
        .querySelector('meta[name="csrf-token"]')
        ?.getAttribute("content");
      if (csrfToken) {
        config.headers["X-CSRF-Token"] = csrfToken;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Handle API errors globally with enhanced error handling
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    // Default error message
    let title = "Error";
    let message = "An unexpected error occurred. Please try again.";
    let notify = true;
    let duration = 5000;

    // Handle different error scenarios
    if (error.response) {
      // The server responded with an error status code
      const status = error.response.status;

      switch (status) {
        case 401:
          // Unauthorized - Handle token expiration
          if (router.currentRoute.value.path !== "/login") {
            console.log("Session expired, redirecting to login");
            store.dispatch("logout");
            router.push("/login");

            title = "Session Expired";
            message = "Your session has expired. Please log in again.";
          } else {
            // Don't show notification if already on login page
            notify = false;
          }
          break;

        case 403:
          // Forbidden
          title = "Access Denied";
          message = "You don't have permission to access this resource.";
          break;

        case 404:
          // Not Found
          title = "Not Found";
          message = "The requested resource was not found.";
          break;

        case 422:
          // Validation Error
          title = "Validation Error";
          // Use server-provided message if available
          message =
            error.response.data?.error ||
            "Please check your input and try again.";
          break;

        case 429:
          // Too Many Requests
          title = "Rate Limited";
          message = "Too many requests. Please try again later.";
          break;

        case 500:
        case 502:
        case 503:
        case 504:
          // Server Errors
          title = "Server Error";
          message =
            "We're experiencing technical difficulties. Please try again later.";
          break;

        default:
          // Use server-provided error message if available
          if (error.response.data?.error) {
            message = error.response.data.error;
          }
      }
    } else if (error.request) {
      // The request was made but no response was received
      if (error.code === "ECONNABORTED") {
        title = "Request Timeout";
        message = "The request timed out. Please try again.";
      } else {
        title = "Network Error";
        message =
          "Unable to connect to the server. Please check your internet connection.";
      }
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error("Error during request setup:", error.message);
    }

    // Show notification to the user using eventBus if needed
    if (notify) {
      eventBus.emit("show-notification", {
        type: "error",
        title,
        message,
        duration,
      });
    }

    // Pass the error along to the component that made the request
    return Promise.reject(error);
  },
);

// Create the app
const app = createApp(App);

// Register Font Awesome icons
registerIcons(app);

// Initialize the store from localStorage before mounting the app
store.dispatch("initializeStore");

// Set base URL from env
const apiUrl = process.env.VUE_APP_API_URL || "http://localhost:8080";
axios.defaults.baseURL = apiUrl;

console.log("API URL:", apiUrl);

// Mount the app
app.use(store).use(router).mount("#app");
