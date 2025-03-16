import { createStore } from "vuex";
import axios from "axios";
import jwtDecode from "jwt-decode";

// Perfect token cache with optimized cleanup and size limiting
class TokenCache {
  constructor(options = {}) {
    this.cache = new Map();
    this.maxSize = options.maxSize || 100; // Limit number of cached tokens
    this.lastCleanup = Date.now();
    this.cleanupInterval = options.cleanupInterval || 60000; // Minimum ms between full cleanups
  }

  // Throttled expired token cleanup
  cleanExpired() {
    const now = Date.now();

    // Only perform full cleanup if enough time has passed since last cleanup
    // or if the cache is getting full
    if (now - this.lastCleanup > this.cleanupInterval || this.cache.size > this.maxSize * 0.8) {
      this.lastCleanup = now;

      let expiredCount = 0;
      for (const [token, item] of this.cache.entries()) {
        if (now > item.expiresAt) {
          this.cache.delete(token);
          expiredCount++;
        }
      }

      // If we're still over max size after removing expired items,
      // remove oldest items based on last accessed time
      if (this.cache.size > this.maxSize) {
        const entries = Array.from(this.cache.entries())
          .sort((a, b) => a[1].lastAccessed - b[1].lastAccessed);

        const toRemove = entries.slice(0, entries.length - this.maxSize);
        for (const [token] of toRemove) {
          this.cache.delete(token);
        }
      }

      return expiredCount; // Return number of expired items removed
    }

    return 0; // No cleanup performed
  }

  get(token) {
    if (!token) return null;

    const item = this.cache.get(token);
    if (!item) return null;

    const now = Date.now();

    // Check if this specific token is expired
    if (now > item.expiresAt) {
      this.cache.delete(token);
      return null;
    }

    // Update last accessed time
    item.lastAccessed = now;

    // Run cleanup occasionally, but not on every access
    if (this.cache.size > 10 && Math.random() < 0.1) {
      this.cleanExpired();
    }

    return item.value;
  }

  set(token, value, ttl) {
    if (!token) return;

    const now = Date.now();
    const defaultTTL = 30 * 60 * 1000;

    this.cache.set(token, {
      value,
      expiresAt: now + (ttl || defaultTTL),
      lastAccessed: now,
      created: now
    });

    // Run cleanup if we've exceeded max size
    if (this.cache.size > this.maxSize) {
      this.cleanExpired();
    }
  }

  has(token) {
    if (!token) return false;

    const item = this.cache.get(token);
    if (!item) return false;

    if (Date.now() > item.expiresAt) {
      this.cache.delete(token);
      return false;
    }

    return true;
  }

  // Returns cache stats - useful for debugging
  getStats() {
    const now = Date.now();
    let expired = 0;

    for (const item of this.cache.values()) {
      if (now > item.expiresAt) expired++;
    }

    return {
      size: this.cache.size,
      expired,
      active: this.cache.size - expired
    };
  }

  clear() {
    this.cache.clear();
    this.lastCleanup = Date.now();
  }
}

// Initialize token cache with options
const tokenCache = new TokenCache({
  maxSize: 50,              // Maximum number of tokens to store
  cleanupInterval: 30000    // Minimum 30 seconds between full cleanups
});

// Simple storage wrapper for localStorage
const storage = {
  getItem(key) {
    try {
      return localStorage.getItem(key);
    } catch (e) {
      console.warn(`Storage access error for ${key}:`, e);
      return null;
    }
  },

  setItem(key, value) {
    try {
      localStorage.setItem(key, value);
    } catch (e) {
      console.warn(`Storage write error for ${key}:`, e);
    }
  },

  removeItem(key) {
    try {
      localStorage.removeItem(key);
    } catch (e) {
      console.warn(`Storage delete error for ${key}:`, e);
    }
  },
};

// Token decoding function with better error handling
const decodeToken = (token) => {
  if (!token) return null;

  // Try to get from cache first
  const cached = tokenCache.get(token);
  if (cached) return cached;

  try {
    const decoded = jwtDecode(token);

    // Validate required fields are present
    if (!decoded || typeof decoded !== 'object') {
      return null;
    }

    // Calculate TTL based on token expiration
    let ttl;
    if (decoded.exp) {
      // Convert exp to milliseconds and subtract current time
      ttl = Math.max(0, decoded.exp * 1000 - Date.now() - 5 * 60 * 1000);

      // If token is already expired, don't cache it
      if (ttl <= 0) return decoded;
    }

    // Cache the token
    tokenCache.set(token, decoded, ttl);
    return decoded;
  } catch (error) {
    if (process.env.NODE_ENV !== "production") {
      console.error("Token decoding error:", error);
    }
    return null;
  }
};

// Token validation function with improved safety
const isTokenValid = (token) => {
  if (!token) return false;

  try {
    const decoded = decodeToken(token);

    // More thorough validation
    if (!decoded) return false;
    if (!decoded.exp) return false;
    if (!decoded.user_id && !decoded.sub) return false;

    // Check expiration
    const nowInSeconds = Math.floor(Date.now() / 1000);
    return decoded.exp > nowInSeconds;
  } catch (error) {
    return false;
  }
};

// Helper to clear auth state
const clearAuthState = (state) => {
  state.user = null;
  state.token = null;
  state.rememberMe = false;

  storage.removeItem("token");
  storage.removeItem("user");
  storage.removeItem("rememberMe");
  tokenCache.clear();

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
      const updatedUser = { ...state.user, ...userUpdates };
      state.user = updatedUser;
      storage.setItem("user", JSON.stringify(updatedUser));
    },
    logout(state) {
      clearAuthState(state);
    },
    initializeStore(state) {
      try {
        const storedToken = storage.getItem("token");
        if (!storedToken) return;

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

        if (!storedRememberMe && !sessionMarker) {
          clearAuthState(state);
          return;
        }

        try {
          sessionStorage.setItem("sessionMarker", "active");
        } catch {
          // Fallback if sessionStorage is unavailable
          console.warn("Session storage unavailable");
        }

        state.token = storedToken;
        state.rememberMe = storedRememberMe;

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

        commit("setToken", token);
        commit("setRememberMe", !!credentials.rememberMe);

        const decoded = decodeToken(token);
        if (!decoded) {
          throw new Error("Invalid authentication token received");
        }

        const user = {
          id: decoded.user_id || 0,
          email: decoded.email || "",
          role: decoded.role || "User",
          name: decoded.name || "Test User",
          bio: decoded.bio || "",
        };

        commit("setUser", user);
        return user;
      } catch (error) {
        throw error;
      }
    },
    async register({ commit }, credentials) {
      try {
        const response = await axios.post("/api/auth/register", credentials);
        const token = response.data.token;

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
    hasRole: (state) => {
      return (role) => {
        return state.user && state.user.role === role;
      };
    },
  },
});
