import { createStore } from "vuex";
import axios from "axios";
import jwtDecode from "jwt-decode";

// Optimized token cache with more efficient cleanup and size limiting
class TokenCache {
  constructor(options = {}) {
    this.cache = new Map();
    this.maxSize = options.maxSize || 100; // Limit number of cached tokens
    this.lastCleanup = Date.now();
    this.cleanupInterval = options.cleanupInterval || 60000; // Minimum ms between full cleanups
    this.accessMap = new Map(); // Map to track access times for O(1) operations
    this.disabled = false; // Safety switch in case of errors
  }

  // More efficient expired token cleanup
  cleanExpired() {
    try {
      const now = Date.now();

      // Only perform full cleanup if enough time has passed since last cleanup
      // or if the cache is getting full
      if (
        now - this.lastCleanup > this.cleanupInterval ||
        this.cache.size > this.maxSize * 0.8
      ) {
        this.lastCleanup = now;

        // First pass: remove expired items
        let expiredCount = 0;
        for (const [token, item] of this.cache.entries()) {
          if (now > item.expiresAt) {
            this.cache.delete(token);
            this.accessMap.delete(token);
            expiredCount++;
          }
        }

        // If we're still over max size after removing expired items,
        // remove least recently used tokens
        if (this.cache.size > this.maxSize) {
          // If we somehow lost track of access times, reset them
          if (this.accessMap.size < this.cache.size) {
            this._resyncAccessMap();
          }

          // Get sorted tokens by access time (oldest first)
          const sortedTokens = [...this.accessMap.entries()]
            .sort((a, b) => a[1] - b[1])
            .map((entry) => entry[0]);

          // Calculate how many tokens to remove
          const excessTokens = this.cache.size - this.maxSize;
          const tokensToRemove = Math.min(excessTokens, sortedTokens.length);

          // Remove the oldest tokens
          for (let i = 0; i < tokensToRemove; i++) {
            const token = sortedTokens[i];
            this.cache.delete(token);
            this.accessMap.delete(token);
          }
        }

        return expiredCount; // Return number of expired items removed
      }

      return 0; // No cleanup performed
    } catch (error) {
      console.error("Token cache cleanup error:", error);
      this.disabled = true; // Disable cache on critical errors
      return 0;
    }
  }

  // Helper method to resync access map if it gets out of sync
  _resyncAccessMap() {
    try {
      const now = Date.now();
      this.accessMap.clear();

      // Rebuild access map from cache
      for (const [token, item] of this.cache.entries()) {
        this.accessMap.set(token, item.lastAccessed || now);
      }
    } catch (error) {
      console.error("Access map resync error:", error);
    }
  }

  get(token) {
    if (!token || this.disabled) return null;

    try {
      const item = this.cache.get(token);
      if (!item) return null;

      const now = Date.now();

      // Check if this specific token is expired
      if (now > item.expiresAt) {
        this.cache.delete(token);
        this.accessMap.delete(token);
        return null;
      }

      // Update last accessed time
      item.lastAccessed = now;
      this.accessMap.set(token, now);

      // Run cleanup occasionally based on time and cache size
      if (
        this.cache.size > 10 &&
        now - this.lastCleanup > this.cleanupInterval / 10
      ) {
        this.cleanExpired();
      }

      return item.value;
    } catch (error) {
      console.warn("Token cache get error:", error);
      return null;
    }
  }

  set(token, value, ttl) {
    if (!token || this.disabled) return;

    try {
      const now = Date.now();
      const defaultTTL = 30 * 60 * 1000;

      this.cache.set(token, {
        value,
        expiresAt: now + (ttl || defaultTTL),
        lastAccessed: now,
        created: now,
      });

      // Update the access map
      this.accessMap.set(token, now);

      // Run cleanup if we've exceeded max size
      if (this.cache.size > this.maxSize) {
        this.cleanExpired();
      }
    } catch (error) {
      console.warn("Token cache set error:", error);
    }
  }

  has(token) {
    if (!token || this.disabled) return false;

    try {
      const item = this.cache.get(token);
      if (!item) return false;

      if (Date.now() > item.expiresAt) {
        this.cache.delete(token);
        this.accessMap.delete(token);
        return false;
      }

      return true;
    } catch (error) {
      console.warn("Token cache has error:", error);
      return false;
    }
  }

  // Returns cache stats - useful for debugging
  getStats() {
    try {
      const now = Date.now();
      let expired = 0;

      for (const item of this.cache.values()) {
        if (now > item.expiresAt) expired++;
      }

      return {
        size: this.cache.size,
        expired,
        active: this.cache.size - expired,
        accessMapSize: this.accessMap.size,
        disabled: this.disabled,
      };
    } catch (error) {
      return { error: "Failed to get stats" };
    }
  }

  clear() {
    try {
      this.cache.clear();
      this.accessMap.clear();
      this.lastCleanup = Date.now();
      this.disabled = false; // Reset disabled state on clear
    } catch (error) {
      console.warn("Token cache clear error:", error);
    }
  }
}

// Initialize token cache with options
const tokenCache = new TokenCache({
  maxSize: 50, // Maximum number of tokens to store
  cleanupInterval: 30000, // Minimum 30 seconds between full cleanups
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
    if (!decoded || typeof decoded !== "object") {
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
