import { createStore } from 'vuex';
import axios from 'axios';
import jwtDecode from 'jwt-decode';

axios.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8080';

export default createStore({
  state: {
    user: JSON.parse(localStorage.getItem('user')) || null,
    token: localStorage.getItem('token') || null,
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
      localStorage.setItem('user', JSON.stringify(user));
    },
    setToken(state, token) {
      state.token = token;
      localStorage.setItem('token', token);
    },
    updateUser(state, userUpdates) {
      state.user = { ...state.user, ...userUpdates };
      localStorage.setItem('user', JSON.stringify(state.user));
    },
    logout(state) {
      state.user = null;
      state.token = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    },
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post('/api/auth/login', credentials);
        commit('setToken', response.data.token);
        const decoded = jwtDecode(response.data.token);
        // Using a dummy name as token may not include it
        commit('setUser', {
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
        const response = await axios.post('/api/auth/register', credentials);
        commit('setToken', response.data.token);
        const decoded = jwtDecode(response.data.token);
        commit('setUser', {
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
      commit('logout');
    },
    updateProfile({ commit }, profileData) {
      // Simulate a profile update; replace with an API call if available.
      commit('updateUser', profileData);
    },
  },
  getters: {
    isAuthenticated: (state) => !!state.token,
    user: (state) => state.user,
  },
});

