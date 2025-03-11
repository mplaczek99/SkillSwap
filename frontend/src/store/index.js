import { createStore } from "vuex";
import axios from "axios";
import jwtDecode from "jwt-decode";

// Set up axios with the API URL from environment variables
const apiUrl = process.env.VUE_APP_API_URL || "http://localhost:8080";
axios.defaults.baseURL = apiUrl;

export default createStore({
  state: {
    user: null,
    token: null,
    rememberMe: false,
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
      localStorage.setItem("user", JSON.stringify(user));
    },
    setToken(state, token) {
      state.token = token;
      localStorage.setItem("token", token);
    },
    setRememberMe(state, value) {
      state.rememberMe = value;
      localStorage.setItem("rememberMe", value.toString());
    },
    updateUser(state, userUpdates) {
      state.user = { ...state.user, ...userUpdates };
      localStorage.setItem("user", JSON.stringify(state.user));
    },
    logout(state) {
      state.user = null;
      state.token = null;
      state.rememberMe = false;
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      localStorage.removeItem("rememberMe");
    },
    initializeStore(state) {
      // Initialize store from localStorage if available
      try {
        const storedUserJSON = localStorage.getItem("user");
        const storedToken = localStorage.getItem("token");
        const storedRememberMe = localStorage.getItem("rememberMe");

        // Check for remember me flag
        const isRemembered = storedRememberMe === "true";

        // If remember me is false and this is a new browser session, don't restore
        // We can detect a new session by checking a session variable
        const sessionMarker = sessionStorage.getItem("sessionMarker");

        if (!isRemembered && !sessionMarker && storedToken) {
          // This is a new session and rememberMe was false, clear stored data
          localStorage.removeItem("token");
          localStorage.removeItem("user");
          // Also clear the state data to ensure user is properly logged out
          state.user = null;
          state.token = null;
          state.rememberMe = false;
          return;
        }

        // Mark this as an active session
        sessionStorage.setItem("sessionMarker", "active");

        if (storedUserJSON) {
          state.user = JSON.parse(storedUserJSON);
        }

        if (storedToken) {
          state.token = storedToken;
        }

        if (storedRememberMe) {
          state.rememberMe = isRemembered;
        }
      } catch (e) {
        console.error("Error initializing store from localStorage:", e);
      }
    },
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/login", credentials);
        commit("setToken", response.data.token);
        commit("setRememberMe", !!credentials.rememberMe);

        try {
          const decoded = jwtDecode(response.data.token);
          commit("setUser", {
            id: decoded.user_id || 0,
            email: decoded.email || "",
            role: decoded.role || "User",
            name: decoded.name || "Test User",
            bio: decoded.bio || "",
          });
        } catch (decodeError) {
          console.error("Error decoding JWT token:", decodeError);
          throw new Error("Invalid authentication token received");
        }
      } catch (error) {
        throw error;
      }
    },
    async register({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/register", credentials);
        commit("setToken", response.data.token);
        commit("setRememberMe", true); // Default to remember for new registrations

        try {
          const decoded = jwtDecode(response.data.token);
          commit("setUser", {
            id: decoded.user_id || 0,
            email: decoded.email || "",
            role: decoded.role || "User",
            name: credentials.name || "New User",
            bio: "",
          });
        } catch (decodeError) {
          console.error("Error decoding JWT token:", decodeError);
          throw new Error("Invalid authentication token received");
        }
      } catch (error) {
        throw error;
      }
    },
    logout({ commit }) {
      commit("logout");
    },
    updateProfile({ commit }, profileData) {
      commit("updateUser", profileData);
    },
    initializeStore({ commit }) {
      commit("initializeStore");
    },
  },
  getters: {
    isAuthenticated: (state) => !!state.token,
    user: (state) => state.user,
  },
});
