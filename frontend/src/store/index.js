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
    updateUser(state, userUpdates) {
      state.user = { ...state.user, ...userUpdates };
      localStorage.setItem("user", JSON.stringify(state.user));
    },
    logout(state) {
      state.user = null;
      state.token = null;
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    },
    initializeStore(state) {
      // Initialize store from localStorage if available
      try {
        const storedUser = JSON.parse(localStorage.getItem("user"));
        const storedToken = localStorage.getItem("token");
        
        if (storedUser) {
          state.user = storedUser;
        }
        
        if (storedToken) {
          state.token = storedToken;
        }
      } catch (e) {
        console.error("Error initializing store from localStorage:", e);
      }
    }
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/login", credentials);
        commit("setToken", response.data.token);
        const decoded = jwtDecode(response.data.token);
        commit("setUser", {
          id: decoded.user_id,
          email: decoded.email,
          role: decoded.role,
          name: "Test User",
          bio: "",
        });
      } catch (error) {
        throw error;
      }
    },
    async register({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/register", credentials);
        commit("setToken", response.data.token);
        const decoded = jwtDecode(response.data.token);
        commit("setUser", {
          id: decoded.user_id,
          email: decoded.email,
          role: decoded.role,
          name: credentials.name,
          bio: "",
        });
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
    }
  },
  getters: {
    isAuthenticated: (state) => !!state.token,
    user: (state) => state.user,
  },
});
