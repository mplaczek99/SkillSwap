import { createStore } from 'vuex';
import axios from 'axios';

// Optionally, set a default base URL (adjust if your backend is hosted elsewhere)
axios.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8080';

export default createStore({
  state: {
    user: null,
    token: localStorage.getItem('token') || null,
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
    },
    setToken(state, token) {
      state.token = token;
      localStorage.setItem('token', token);
    },
    logout(state) {
      state.user = null;
      state.token = null;
      localStorage.removeItem('token');
    },
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post('/api/auth/login', credentials);
        commit('setToken', response.data.token);
        // Optionally, decode the token here to set user details
      } catch (error) {
        throw error;
      }
    },
    async register({ commit }, credentials) {
      try {
        const response = await axios.post('/api/auth/register', credentials);
        commit('setToken', response.data.token);
      } catch (error) {
        throw error;
      }
    },
    logout({ commit }) {
      commit('logout');
    },
  },
  getters: {
    isAuthenticated: (state) => !!state.token,
  },
});

