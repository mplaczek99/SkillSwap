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
      // Helper function to clear authentication state
      const clearAuth = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        localStorage.removeItem("rememberMe");
        state.user = null;
        state.token = null;
        state.rememberMe = false;
      };

      try {
        // Read from localStorage only once
        const storedToken = localStorage.getItem("token");

        // Early return if no token exists
        if (!storedToken) return;

        const storedRememberMe = localStorage.getItem("rememberMe") === "true";
        const sessionMarker = sessionStorage.getItem("sessionMarker");

        // If remember me is false and this is a new browser session, don't restore
        if (!storedRememberMe && !sessionMarker) {
          clearAuth();
          return;
        }

        // Validate token before restoring other state
        try {
          // Decode the JWT token
          const decoded = jwtDecode(storedToken);

          // Check if token has expired
          if (decoded.exp && decoded.exp < Date.now() / 1000) {
            console.log(
              "Stored token has expired, clearing authentication data",
            );
            clearAuth();
            return;
          }

          // Token is valid, mark this as an active session
          sessionStorage.setItem("sessionMarker", "active");

          // Restore state
          state.token = storedToken;
          state.rememberMe = storedRememberMe;

          // Only parse user JSON if we need it
          const storedUserJSON = localStorage.getItem("user");
          if (storedUserJSON) {
            state.user = JSON.parse(storedUserJSON);
          }
        } catch (tokenError) {
          console.error("Invalid token format:", tokenError);
          clearAuth();
        }
      } catch (e) {
        console.error("Error initializing store from localStorage:", e);
        clearAuth();
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
