import { createStore } from "vuex";
import axios from "axios";
import jwtDecode from "jwt-decode";

// Set up axios with the API URL from environment variables
const apiUrl = process.env.VUE_APP_API_URL || "http://localhost:8080";
axios.defaults.baseURL = apiUrl;

// Improved LRU Token Cache with O(1) operations
class LRUCache {
  constructor(capacity) {
    this.capacity = capacity;
    this.cache = new Map(); // Map maintains insertion order
    this.hits = 0;
    this.misses = 0;
  }

  // Get with stats tracking
  get(key) {
    if (!this.cache.has(key)) {
      this.misses++;
      return null;
    }

    // Move the accessed item to the end (most recently used)
    const value = this.cache.get(key);
    this.cache.delete(key);
    this.cache.set(key, value);
    this.hits++;
    return value;
  }

  // Set with automatic eviction
  set(key, value) {
    // If key exists, delete it first to update access order
    if (this.cache.has(key)) {
      this.cache.delete(key);
    }
    // If at capacity, delete oldest item (first item in Map)
    else if (this.cache.size >= this.capacity) {
      const firstKey = this.cache.keys().next().value;
      this.cache.delete(firstKey);
    }

    // Add the new item (will be last in Map)
    this.cache.set(key, value);
  }

  // Check if key exists without updating access order
  has(key) {
    return this.cache.has(key);
  }

  // Get cache performance statistics
  getStats() {
    const total = this.hits + this.misses;
    const hitRate = total > 0 ? ((this.hits / total) * 100).toFixed(2) : 0;
    return { hits: this.hits, misses: this.misses, hitRate: `${hitRate}%` };
  }

  // Clear the cache
  clear() {
    this.cache.clear();
  }
}

// Initialize token cache with size limit
const tokenCache = new LRUCache(20); // Increased from 10 for better hit rate

// Storage access optimization with memoization
const storage = {
  memo: new Map(),

  // Get item with in-memory caching
  getItem(key) {
    if (!this.memo.has(key)) {
      try {
        const value = localStorage.getItem(key);
        this.memo.set(key, value);
      } catch (e) {
        console.warn(`Storage access error for ${key}:`, e);
        return null;
      }
    }
    return this.memo.get(key);
  },

  // Set item and update cache
  setItem(key, value) {
    try {
      localStorage.setItem(key, value);
      this.memo.set(key, value);
    } catch (e) {
      console.warn(`Storage write error for ${key}:`, e);
    }
  },

  // Remove item and update cache
  removeItem(key) {
    try {
      localStorage.removeItem(key);
      this.memo.delete(key);
    } catch (e) {
      console.warn(`Storage delete error for ${key}:`, e);
    }
  },

  // Batch operations for better performance
  batchOperation(operations) {
    operations.forEach((op) => {
      if (op.type === "set") {
        this.setItem(op.key, op.value);
      } else if (op.type === "remove") {
        this.removeItem(op.key);
      }
    });
  },
};

// Optimized token decoder with efficient error handling
const decodeToken = (token) => {
  if (!token) return null;

  // Return cached result if available - O(1) operation
  const cached = tokenCache.get(token);
  if (cached) return cached;

  try {
    // Decode and add to cache with timestamp
    const decoded = jwtDecode(token);
    tokenCache.set(token, decoded);
    return decoded;
  } catch (error) {
    // Only log detailed errors in development
    if (process.env.NODE_ENV !== "production") {
      console.error("Token decoding error:", error);
    }
    return null;
  }
};

// Fast token validation with minimal processing
const isTokenValid = (token) => {
  if (!token) return false;

  try {
    const decoded = decodeToken(token);
    // Simple numeric comparison is faster than Date operations
    return (
      decoded && decoded.exp && decoded.exp > Math.floor(Date.now() / 1000)
    );
  } catch {
    return false;
  }
};

// Helper function to efficiently clear authentication state
const clearAuthState = (state) => {
  // Update state
  state.user = null;
  state.token = null;
  state.rememberMe = false;

  // Batch localStorage operations
  storage.batchOperation([
    { type: "remove", key: "token" },
    { type: "remove", key: "user" },
    { type: "remove", key: "rememberMe" },
  ]);

  // Clear token cache to prevent memory leaks
  tokenCache.clear();

  // Clear session marker
  try {
    sessionStorage.removeItem("sessionMarker");
  } catch (e) {
    // Ignore sessionStorage errors
  }
};

export default createStore({
  state: {
    user: null,
    token: null,
    rememberMe: false,
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
      storage.setItem("user", JSON.stringify(user));
    },
    setToken(state, token) {
      state.token = token;
      storage.setItem("token", token);
    },
    setRememberMe(state, value) {
      state.rememberMe = value;
      storage.setItem("rememberMe", value.toString());
    },
    updateUser(state, userUpdates) {
      // Create new object only once instead of spreading twice
      const updatedUser = { ...state.user, ...userUpdates };
      state.user = updatedUser;
      storage.setItem("user", JSON.stringify(updatedUser));
    },
    logout(state) {
      clearAuthState(state);
    },
    initializeStore(state) {
      try {
        // Fast path: check for token first
        const storedToken = storage.getItem("token");
        if (!storedToken) return;

        // Check token validity before any other operations
        if (!isTokenValid(storedToken)) {
          clearAuthState(state);
          return;
        }

        const storedRememberMe = storage.getItem("rememberMe") === "true";
        let sessionMarker;

        try {
          sessionMarker = sessionStorage.getItem("sessionMarker");
        } catch {
          sessionMarker = null;
        }

        // Early exit for session-only logins on new session
        if (!storedRememberMe && !sessionMarker) {
          clearAuthState(state);
          return;
        }

        // Set session marker immediately to avoid redundant checks
        try {
          sessionStorage.setItem("sessionMarker", "active");
        } catch {
          // Fallback if sessionStorage is unavailable
          console.warn("Session storage unavailable");
        }

        // Restore state efficiently
        state.token = storedToken;
        state.rememberMe = storedRememberMe;

        // Parse user data only once
        const storedUserJSON = storage.getItem("user");
        if (storedUserJSON) {
          try {
            state.user = JSON.parse(storedUserJSON);
          } catch {
            state.user = null;
          }
        }
      } catch (error) {
        console.warn("Authentication state restoration failed:", error);
        clearAuthState(state);
      }
    },
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/login", credentials);
        const token = response.data.token;

        // Process with optimized flow to reduce work
        commit("setToken", token);
        commit("setRememberMe", !!credentials.rememberMe);

        const decoded = decodeToken(token);
        if (!decoded) {
          throw new Error("Invalid authentication token received");
        }

        // Construct user object once, avoid repeated spread operations
        const user = {
          id: decoded.user_id || 0,
          email: decoded.email || "",
          role: decoded.role || "User",
          name: decoded.name || "Test User",
          bio: decoded.bio || "",
        };

        commit("setUser", user);
        return user; // Return user for chaining
      } catch (error) {
        throw error;
      }
    },
    async register({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/register", credentials);
        const token = response.data.token;

        // Same optimized flow as login
        commit("setToken", token);
        commit("setRememberMe", true);

        const decoded = decodeToken(token);
        if (!decoded) {
          throw new Error("Invalid authentication token received");
        }

        const user = {
          id: decoded.user_id || 0,
          email: decoded.email || "",
          role: decoded.role || "User",
          name: credentials.name || "New User",
          bio: "",
        };

        commit("setUser", user);
        return user;
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
    // Memoize hasRole for better performance when called multiple times with same role
    hasRole: (state) => {
      const roleCache = new Map();
      return (role) => {
        if (!state.user) return false;
        if (!roleCache.has(role)) {
          roleCache.set(role, state.user.role === role);
        }
        return roleCache.get(role);
      };
    },
    // New getter to expose token cache stats in development
    tokenCacheStats: () => {
      return process.env.NODE_ENV !== "production"
        ? tokenCache.getStats()
        : null;
    },
  },
});
