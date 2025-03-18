<template>
  <div class="videos-list-container">
    <div class="container">
      <div class="page-header">
        <h2>My Video Library</h2>
        <p class="description">
          Create, manage and share your knowledge through video content
        </p>
        <router-link to="/upload-video" class="btn btn-primary upload-btn">
          <font-awesome-icon icon="upload" />
          Upload New Video
        </router-link>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="spinner-container">
          <font-awesome-icon icon="spinner" class="spin" />
        </div>
        <p>Loading your videos...</p>
      </div>

      <div v-else-if="error" class="error-message">
        <font-awesome-icon icon="exclamation-circle" />
        {{ error }}
      </div>

      <div v-else-if="videos.length === 0" class="empty-state">
        <div class="empty-state-icon">
          <font-awesome-icon icon="film" class="empty-icon" />
        </div>
        <h3>Your Video Library is Empty</h3>
        <p>Start sharing your skills by uploading your first video tutorial.</p>
        <router-link to="/upload-video" class="btn btn-primary">
          <font-awesome-icon icon="upload" />
          Upload Your First Video
        </router-link>
      </div>

      <div v-else class="videos-grid">
        <div v-for="video in videos" :key="video.id" class="video-card">
          <div class="video-thumbnail" @click="playVideo(video)">
            <div v-if="video.hasThumbnail" class="thumbnail-image">
              <img :src="getThumbnailUrl(video)" :alt="getDisplayName(video)" />
              <div class="play-overlay">
                <font-awesome-icon icon="play-circle" />
              </div>
            </div>
            <div v-else class="no-thumbnail">
              <font-awesome-icon icon="film" />
              <div class="play-text">Play Video</div>
            </div>
          </div>
          <div class="video-details">
            <h3 class="video-title" @click="playVideo(video)">
              {{ getDisplayName(video) }}
            </h3>
            <p class="video-meta">
              <span class="video-size"
                ><font-awesome-icon icon="file" />
                {{ formatFileSize(video.size) }}</span
              >
              <span class="video-date"
                ><font-awesome-icon icon="calendar-alt" />
                {{ formatDate(video.uploadedAt) }}</span
              >
            </p>
            <div class="video-actions">
              <button
                class="btn btn-primary btn-sm action-btn"
                @click="playVideo(video)"
              >
                <font-awesome-icon icon="play" />
                Play
              </button>
              <button
                class="btn btn-outline btn-sm action-btn"
                @click="downloadVideo(video)"
              >
                <font-awesome-icon icon="download" />
                Download
              </button>
              <button
                class="btn btn-outline btn-sm action-btn share-btn"
                @click="shareVideo(video)"
              >
                <font-awesome-icon icon="share-alt" />
                Share
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Video Modal -->
      <div v-if="showVideoModal" class="video-modal">
        <div class="video-modal-overlay" @click="closeVideoModal"></div>
        <div class="video-modal-content">
          <div class="modal-header">
            <h2 v-if="currentVideo">{{ getDisplayName(currentVideo) }}</h2>
            <button class="close-modal-btn" @click="closeVideoModal">
              <font-awesome-icon icon="times" />
            </button>
          </div>

          <div class="video-player-container">
            <video
              v-if="currentVideo && !videoError"
              controls
              autoplay
              class="video-player"
              @error="handleVideoError"
              ref="videoPlayer"
            >
              <source
                :src="getVideoUrl(currentVideo)"
                :type="getVideoType(currentVideo.name)"
              />
              Your browser does not support the video tag.
            </video>

            <div v-if="videoError" class="video-error">
              <div class="error-icon">
                <font-awesome-icon icon="exclamation-triangle" />
              </div>
              <h3>Playback Error</h3>
              <p>
                We couldn't play this video. Please try again or download it to
                watch locally.
              </p>

              <div class="modal-button-group">
                <button @click="retryVideo" class="btn btn-primary">
                  <font-awesome-icon icon="redo" />
                  Try Again
                </button>
                <button
                  @click="downloadVideo(currentVideo)"
                  class="btn btn-outline"
                >
                  <font-awesome-icon icon="download" />
                  Download
                </button>
                <button @click="closeVideoModal" class="btn btn-outline">
                  Close
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Share Modal -->
      <div v-if="showShareModal" class="share-modal">
        <div class="share-modal-overlay" @click="showShareModal = false"></div>
        <div class="share-modal-content">
          <div class="modal-header">
            <h3>Share Video</h3>
            <button class="close-modal-btn" @click="showShareModal = false">
              <font-awesome-icon icon="times" />
            </button>
          </div>
          <div class="share-options">
            <p>Share "{{ shareVideoName }}" with others:</p>
            <div class="share-link-container">
              <input
                type="text"
                readonly
                :value="shareVideoUrl"
                ref="shareUrlInput"
                class="share-url-input"
              />
              <button class="btn btn-primary copy-btn" @click="copyShareLink">
                <font-awesome-icon :icon="copiedIcon" />
                {{ copiedText }}
              </button>
            </div>
            <div class="share-platforms">
              <button class="platform-btn email-btn">
                <font-awesome-icon icon="envelope" />
                Email
              </button>
              <button class="platform-btn slack-btn">
                <font-awesome-icon icon="comment" />
                Slack
              </button>
              <button class="platform-btn teams-btn">
                <font-awesome-icon icon="users" />
                Teams
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
import eventBus from "@/utils/eventBus";

export default {
  name: "VideosList",
  data() {
    return {
      videos: [],
      loading: true,
      error: null,
      showVideoModal: false,
      currentVideo: null,
      videoError: false,
      showShareModal: false,
      shareVideoName: "",
      shareVideoUrl: "",
      copied: false,
    };
  },
  computed: {
    baseUrl() {
      return axios.defaults.baseURL || "";
    },
    copiedText() {
      return this.copied ? "Copied!" : "Copy Link";
    },
    copiedIcon() {
      return this.copied ? "check" : "copy";
    },
  },
  created() {
    this.fetchVideos();
  },
  methods: {
    async fetchVideos() {
      this.loading = true;
      this.error = null;

      try {
        const response = await axios.get("/api/videos", {
          timeout: 15000,
        });

        // Log the raw response in development for debugging
        if (process.env.NODE_ENV !== "production") {
          console.debug("API Response:", response.data);
        }

        this.videos = response.data || [];
      } catch (error) {
        console.error("Error fetching videos:", error);
        this.videos = [];
        this.error = this.getErrorMessage(error);

        eventBus.emit("show-notification", {
          type: "error",
          title: "Failed to Load Videos",
          message: this.error,
          duration: 5000,
        });
      } finally {
        this.loading = false;
      }
    },

    getThumbnailUrl(video) {
      return `${this.baseUrl}/uploads/${video.thumbnail}`;
    },

    getVideoUrl(video) {
      // Always use the internal name (not originalFilename) for the actual URL
      return `${this.baseUrl}/uploads/${encodeURIComponent(video.name)}`;
    },

    getDisplayName(video) {
      // First priority: Use originalFilename if available and not empty
      if (video.originalFilename && video.originalFilename.trim() !== "") {
        return video.originalFilename;
      }

      // Second priority: Generate a readable name from timestamp if it's a UUID/hash filename
      if (video.name && /^[a-f0-9]{8,}\.[\w]+$/i.test(video.name)) {
        const extension = video.name.split(".").pop();
        const date = new Date(video.uploadedAt);
        const formattedDate = date.toLocaleDateString(undefined, {
          year: "numeric",
          month: "short",
          day: "numeric",
        });
        return `Video (${formattedDate}).${extension}`;
      }

      // Third priority: Use the stored name
      if (video.name) {
        return video.name;
      }

      // Last resort
      return "Unnamed Video";
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
      this.currentVideo = video;
      this.showVideoModal = true;
      this.videoError = false;

      // Reset video position if it was played before
      this.$nextTick(() => {
        if (this.$refs.videoPlayer) {
          this.$refs.videoPlayer.currentTime = 0;
        }
      });
    },

    closeVideoModal() {
      this.showVideoModal = false;
      this.currentVideo = null;
      this.videoError = false;
    },

    handleVideoError() {
      console.error("Video playback error:", this.currentVideo);
      this.videoError = true;
    },

    retryVideo() {
      if (!this.currentVideo) return;

      this.videoError = false;

      // Force video reload by manipulating the DOM
      if (this.$refs.videoPlayer) {
        const video = this.$refs.videoPlayer;
        video.load();
      }
    },

    getVideoType(filename) {
      const ext = filename.split(".").pop().toLowerCase();
      const mimeTypes = {
        mp4: "video/mp4",
        webm: "video/webm",
        ogg: "video/ogg",
        mov: "video/quicktime",
        avi: "video/x-msvideo",
        wmv: "video/x-ms-wmv",
        mkv: "video/x-matroska",
      };
      return mimeTypes[ext] || "video/mp4";
    },

    downloadVideo(video) {
      if (!video) return;

      // Create a link to download the video with original filename
      const link = document.createElement("a");
      link.href = this.getVideoUrl(video);

      // Use the display name for the download
      link.download = this.getDisplayName(video);

      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // Show notification
      eventBus.emit("show-notification", {
        type: "success",
        title: "Download Started",
        message: `Downloading "${this.getDisplayName(video)}"`,
        duration: 3000,
      });
    },

    shareVideo(video) {
      if (!video) return;

      this.shareVideoName = this.getDisplayName(video);
      this.shareVideoUrl = `${window.location.origin}/view/${video.id}`;
      this.showShareModal = true;
      this.copied = false;

      // Focus on the input after the modal is shown
      this.$nextTick(() => {
        if (this.$refs.shareUrlInput) {
          this.$refs.shareUrlInput.focus();
          this.$refs.shareUrlInput.select();
        }
      });
    },

    copyShareLink() {
      if (!this.$refs.shareUrlInput) return;

      const textToCopy = this.$refs.shareUrlInput.value;

      // Use modern Clipboard API with fallback to older method
      if (navigator.clipboard && navigator.clipboard.writeText) {
        // Modern approach using Clipboard API
        navigator.clipboard
          .writeText(textToCopy)
          .then(() => {
            this.handleCopySuccess();
          })
          .catch((error) => {
            console.error("Failed to copy text: ", error);
            // Fallback to the old method if permission denied or other error
            this.fallbackCopyToClipboard();
          });
      } else {
        // Fallback for browsers that don't support clipboard API
        this.fallbackCopyToClipboard();
      }
    },

    // Fallback method using the older execCommand approach
    fallbackCopyToClipboard() {
      const input = this.$refs.shareUrlInput;
      input.select();
      try {
        const success = document.execCommand("copy");
        if (success) {
          this.handleCopySuccess();
        } else {
          console.error("Fallback clipboard copy failed");
          eventBus.emit("show-notification", {
            type: "error",
            title: "Copy Failed",
            message: "Could not copy to clipboard. Please try again.",
            duration: 3000,
          });
        }
      } catch (err) {
        console.error("Fallback clipboard copy error:", err);
      }
    },

    // Common success handler
    handleCopySuccess() {
      this.copied = true;

      // Reset copied state after a delay
      setTimeout(() => {
        this.copied = false;
      }, 2000);

      // Show notification
      eventBus.emit("show-notification", {
        type: "success",
        title: "Link Copied",
        message: "Video link copied to clipboard",
        duration: 2000,
      });
    },

    getErrorMessage(error) {
      let errorMessage = "Failed to load videos. Please try again later.";

      if (error.response) {
        const status = error.response.status;

        if (status === 401) {
          errorMessage = "Your session has expired. Please login again.";
          this.$router.push("/login");
        } else if (status === 403) {
          errorMessage = "You don't have permission to access these videos.";
        } else if (status >= 500) {
          errorMessage =
            "The video service is currently unavailable. Please try again later.";
        }

        if (error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error;
        }
      } else if (error.request) {
        if (error.code === "ECONNABORTED") {
          errorMessage = "Request timed out. Please try again later.";
        } else {
          errorMessage =
            "Network error. Please check your connection and try again.";
        }
      }

      return errorMessage;
    },
  },
};
</script>

<style scoped>
.videos-list-container {
  padding: 2rem 0 4rem;
  background: linear-gradient(to bottom, #f8fafc, #f1f5f9);
  min-height: 80vh;
}

.page-header {
  text-align: center;
  margin-bottom: 3rem;
  position: relative;
}

.page-header h2 {
  color: var(--primary-color);
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  font-weight: 700;
  position: relative;
  display: inline-block;
}

.page-header h2::after {
  content: "";
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 4px;
  background: var(--primary-color);
  border-radius: 2px;
}

.description {
  color: var(--medium);
  font-size: 1.2rem;
  max-width: 600px;
  margin: 1rem auto 2rem;
}

.upload-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 50px;
  font-weight: 600;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.upload-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  text-align: center;
}

.spinner-container {
  margin-bottom: 1.5rem;
}

.spin {
  font-size: 2.5rem;
  color: var(--primary-color);
  animation: spin 1.5s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);
  min-height: 400px;
}

.empty-state-icon {
  width: 120px;
  height: 120px;
  background: var(--primary-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 2rem;
}

.empty-icon {
  font-size: 3.5rem;
  color: var(--primary-color);
}

.empty-state h3 {
  font-size: 1.75rem;
  color: var(--dark);
  margin-bottom: 1rem;
}

.empty-state p {
  color: var(--medium);
  max-width: 400px;
  margin: 0 auto 2rem;
  font-size: 1.1rem;
}

/* Error State */
.error-message {
  background-color: #fee2e2;
  border-left: 4px solid var(--error-color);
  padding: 1.5rem;
  border-radius: 0.5rem;
  color: #b91c1c;
  display: flex;
  align-items: center;
  gap: 1rem;
  margin: 2rem 0;
  font-weight: 500;
}

.error-message svg {
  font-size: 1.5rem;
  flex-shrink: 0;
}

/* Videos Grid */
.videos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 2rem;
}

.video-card {
  background: white;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.video-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.1);
}

.video-thumbnail {
  aspect-ratio: 16 / 9;
  background-color: #1a1a1a;
  position: relative;
  overflow: hidden;
  cursor: pointer;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  position: relative;
}

.thumbnail-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.5s ease;
}

.video-card:hover .thumbnail-image img {
  transform: scale(1.05);
}

.play-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.play-overlay svg {
  font-size: 4rem;
  color: white;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

.video-card:hover .play-overlay {
  opacity: 1;
}

.no-thumbnail {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: linear-gradient(to bottom right, #2c3e50, #34495e);
  color: white;
}

.no-thumbnail svg {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.play-text {
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.video-details {
  padding: 1.5rem;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.video-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
  color: var(--dark);
  line-height: 1.4;
  cursor: pointer;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s ease;
}

.video-title:hover {
  color: var(--primary-color);
}

.video-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1.25rem;
  color: var(--medium);
  font-size: 0.9rem;
}

.video-size,
.video-date {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.video-actions {
  margin-top: auto;
  display: flex;
  gap: 0.75rem;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  white-space: nowrap;
}

.share-btn {
  background-color: #4f46e5;
  color: white;
  border-color: #4f46e5;
}

.share-btn:hover {
  background-color: #4338ca;
  border-color: #4338ca;
}

/* Modal Styles */
.video-modal,
.share-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-modal-overlay,
.share-modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(4px);
}

.video-modal-content {
  position: relative;
  background: white;
  border-radius: 1rem;
  width: 90%;
  max-width: 1000px;
  max-height: 90vh;
  overflow: hidden;
  z-index: 1001;
  box-shadow: 0 20px 25px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
}

.share-modal-content {
  position: relative;
  background: white;
  border-radius: 1rem;
  width: 90%;
  max-width: 500px;
  z-index: 1001;
  box-shadow: 0 20px 25px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h2,
.modal-header h3 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--dark);
}

.close-modal-btn {
  background: none;
  border: none;
  color: var(--medium);
  font-size: 1.25rem;
  cursor: pointer;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.close-modal-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--dark);
}

.video-player-container {
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #000;
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-player {
  width: 100%;
  height: 100%;
}

.video-error {
  padding: 2rem;
  text-align: center;
  background: #fee2e2;
  color: #b91c1c;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.video-error h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.video-error p {
  max-width: 500px;
  margin: 0 auto 2rem;
}

.modal-button-group {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

/* Share Modal Styles */
.share-options {
  padding: 1.5rem;
}

.share-options p {
  margin-top: 0;
  margin-bottom: 1rem;
  color: var(--dark);
}

.share-link-container {
  display: flex;
  margin-bottom: 1.5rem;
}

.share-url-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid #e2e8f0;
  border-right: none;
  border-radius: 0.5rem 0 0 0.5rem;
  font-size: 0.9rem;
}

.copy-btn {
  border-radius: 0 0.5rem 0.5rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 100px;
}

.share-platforms {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.platform-btn {
  width: 80px;
  height: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 1rem;
  border: none;
  background: #f8fafc;
  color: var(--dark);
  cursor: pointer;
  transition: all 0.2s ease;
  gap: 0.5rem;
}

.platform-btn svg {
  font-size: 1.5rem;
}

.platform-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.email-btn:hover {
  background: #f87171;
  color: white;
}

.slack-btn:hover {
  background: #4f46e5;
  color: white;
}

.teams-btn:hover {
  background: #4338ca;
  color: white;
}

/* Responsive styles */
@media (max-width: 768px) {
  .videos-grid {
    grid-template-columns: 1fr;
  }

  .modal-button-group {
    flex-direction: column;
  }

  .share-platforms {
    flex-wrap: wrap;
  }

  .platform-btn {
    width: 70px;
    height: 70px;
  }

  .page-header h2 {
    font-size: 2rem;
  }

  .description {
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .video-actions {
    flex-direction: column;
  }

  .share-link-container {
    flex-direction: column;
  }

  .share-url-input {
    border-radius: 0.5rem 0.5rem 0 0;
    border-right: 1px solid #e2e8f0;
    border-bottom: none;
  }

  .copy-btn {
    border-radius: 0 0 0.5rem 0.5rem;
  }
}
</style>
