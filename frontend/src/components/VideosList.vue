<template>
  <div class="videos-list-container">
    <div class="container">
      <h2>My Videos</h2>
      <p class="description">Browse your uploaded video content</p>

      <div v-if="loading" class="loading-state">
        <font-awesome-icon icon="spinner" class="spin" />
        <p>Loading videos...</p>
      </div>

      <div v-else-if="error" class="error-message">
        <font-awesome-icon icon="exclamation-circle" />
        {{ error }}
      </div>

      <div v-else-if="videos.length === 0" class="empty-state">
        <font-awesome-icon icon="film" class="empty-icon" />
        <h3>No Videos Found</h3>
        <p>You haven't uploaded any videos yet.</p>
        <router-link to="/upload-video" class="btn btn-primary">
          Upload Your First Video
        </router-link>
      </div>

      <div v-else class="videos-grid">
        <div v-for="video in videos" :key="video.id" class="video-card">
          <div class="video-thumbnail">
            <div v-if="video.hasThumbnail" class="thumbnail-image">
              <img :src="getThumbnailUrl(video)" :alt="video.name" />
              <div class="play-overlay">
                <font-awesome-icon icon="play-circle" />
              </div>
            </div>
            <div v-else class="no-thumbnail">
              <font-awesome-icon icon="film" />
            </div>
          </div>
          <div class="video-details">
            <h3 class="video-title">{{ formatFileName(video.name) }}</h3>
            <p class="video-meta">
              <span>{{ formatFileSize(video.size) }}</span>
              <span>· {{ formatDate(video.uploadedAt) }}</span>
            </p>
            <div class="video-actions">
              <button class="btn btn-primary btn-sm" @click="playVideo(video)">
                <font-awesome-icon icon="play" />
                Play
              </button>
              <button
                class="btn btn-outline btn-sm"
                @click="downloadVideo(video)"
              >
                <font-awesome-icon icon="download" />
                Download
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "VideosList",
  data() {
    return {
      videos: [],
      loading: true,
      error: null,
    };
  },
  created() {
    this.fetchVideos();
  },
  methods: {
    async fetchVideos() {
      this.loading = true;
      this.error = null;

      try {
        const response = await axios.get("/api/videos");
        this.videos = response.data;
      } catch (error) {
        console.error("Error fetching videos:", error);

        if (error.response && error.response.status === 401) {
          this.error = "Your session has expired. Please login again.";
        } else {
          this.error = "Failed to load videos. Please try again later.";
        }
      } finally {
        this.loading = false;
      }
    },

    getThumbnailUrl(video) {
      // Use the API URL from the store configuration
      const baseUrl = axios.defaults.baseURL || "";
      return `${baseUrl}/uploads/${video.thumbnail}`;
    },

    formatFileName(fileName) {
      // Remove file extension and replace underscores/hyphens with spaces
      return fileName
        .replace(/\.[^/.]+$/, "") // Remove extension
        .replace(/_|-/g, " "); // Replace underscores and hyphens with spaces
    },

    formatFileSize(bytes) {
      if (bytes < 1024) {
        return bytes + " bytes";
      } else if (bytes < 1024 * 1024) {
        return (bytes / 1024).toFixed(1) + " KB";
      } else {
        return (bytes / (1024 * 1024)).toFixed(1) + " MB";
      }
    },

    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    },

    playVideo(video) {
      // Get the API base URL from the store
      const baseUrl = axios.defaults.baseURL || "";
      const videoUrl = `${baseUrl}/uploads/${video.name}`;

      // Create a modal or overlay to play the video
      const win = window.open("", "_blank");
      win.document.write(`
        <html>
          <head>
            <title>${this.formatFileName(video.name)}</title>
            <style>
              body {
                margin: 0;
                padding: 0;
                background: #000;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
              }
              video {
                max-width: 100%;
                max-height: 100vh;
              }
            </style>
          </head>
          <body>
            <video controls autoplay>
              <source src="${videoUrl}" type="video/mp4">
              Your browser does not support the video tag.
            </video>
          </body>
        </html>
      `);
    },

    downloadVideo(video) {
      // Get the API base URL from the store
      const baseUrl = axios.defaults.baseURL || "";
      const videoUrl = `${baseUrl}/uploads/${video.name}`;

      // Create a link to download the video
      const link = document.createElement("a");
      link.href = videoUrl;
      link.download = video.name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
  },
};
</script>

<style scoped>
.videos-list-container {
  padding-bottom: var(--space-12);
}

h2 {
  color: var(--primary-color);
  text-align: center;
  margin-bottom: var(--space-2);
  font-size: var(--font-size-3xl);
}

.description {
  text-align: center;
  color: var(--medium);
  margin-bottom: var(--space-8);
  font-size: var(--font-size-lg);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--space-12) 0;
}

.spin {
  animation: spin 1s linear infinite;
  font-size: var(--font-size-3xl);
  margin-bottom: var(--space-4);
  color: var(--primary-color);
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.error-message {
  background-color: var(--error-color);
  color: white;
  padding: var(--space-4);
  border-radius: var(--radius-md);
  margin: var(--space-8) 0;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.empty-state {
  text-align: center;
  padding: var(--space-12) 0;
}

.empty-icon {
  font-size: var(--font-size-4xl);
  color: var(--medium);
  opacity: 0.5;
  margin-bottom: var(--space-4);
}

.videos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-6);
}

.video-card {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  transition:
    transform var(--transition-normal) ease,
    box-shadow var(--transition-normal) ease;
}

.video-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.video-thumbnail {
  height: 180px;
  background-color: var(--dark);
  position: relative;
}

.thumbnail-image {
  position: relative;
  width: 100%;
  height: 100%;
}

.thumbnail-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.play-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.3);
  opacity: 0;
  transition: opacity var(--transition-fast) ease;
}

.play-overlay svg {
  font-size: var(--font-size-4xl);
  color: white;
}

.thumbnail-image:hover .play-overlay {
  opacity: 1;
}

.no-thumbnail {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  font-size: var(--font-size-3xl);
  color: var(--medium);
  opacity: 0.5;
}

.video-details {
  padding: var(--space-4);
}

.video-title {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.video-meta {
  display: flex;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--medium);
  margin-bottom: var(--space-4);
}

.video-actions {
  display: flex;
  gap: var(--space-2);
}

@media (max-width: 768px) {
  .videos-grid {
    grid-template-columns: 1fr;
  }
}
</style>
