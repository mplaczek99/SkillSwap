import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import eventBus from "./utils/eventBus"; // Import the eventBus

// Import the icon registration function
import { registerIcons } from "./utils/icons";

// Set up axios interceptors to add the auth token to all requests
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Handle 401 responses globally
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (
      error.response &&
      error.response.status === 401 &&
      router.currentRoute.value.path !== "/login"
    ) {
      console.log("Session expired, redirecting to login");
      // Token is expired or invalid
      store.dispatch("logout"); // Use dispatch instead of commit for actions
      router.push("/login");

      // Show notification to the user using eventBus
      eventBus.emit("show-notification", {
        type: "warning",
        title: "Session Expired",
        message: "Your session has expired. Please log in again.",
        duration: 5000,
      });
    }
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
