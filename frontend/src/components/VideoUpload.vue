<template>
  <div class="video-upload-container">
    <h2>Upload a Video</h2>
    <p class="description">Share your skills through video tutorials</p>

    <div class="upload-box">
      <div class="file-selection">
        <input
          type="file"
          @change="onFileSelected"
          accept="video/*"
          id="video-file"
          class="file-input"
        />
        <label for="video-file" class="file-label">
          <font-awesome-icon icon="upload" class="icon" />
          <span>Select video file</span>
        </label>
        <p class="file-info">Maximum size: 100MB</p>
      </div>

      <div v-if="selectedFile" class="selected-file-details">
        <h3>Selected File</h3>
        <p><strong>Name:</strong> {{ selectedFile.name }}</p>
        <p><strong>Size:</strong> {{ formatFileSize(selectedFile.size) }}</p>
        <p><strong>Type:</strong> {{ selectedFile.type }}</p>
      </div>
    </div>

    <div v-if="uploadProgress > 0" class="progress-area">
      <div class="progress-bar">
        <div
          class="progress-fill"
          :style="{ width: `${uploadProgress}%` }"
        ></div>
      </div>
      <div class="progress-text">{{ uploadProgress }}%</div>
    </div>

    <div class="button-area">
      <button
        @click="uploadVideo"
        :disabled="!selectedFile || uploadProgress > 0"
        class="upload-button"
      >
        <font-awesome-icon
          v-if="uploadProgress > 0 && uploadProgress < 100"
          icon="spinner"
          class="spin"
        />
        <span v-else>Upload Video</span>
      </button>
    </div>

    <div v-if="errorMessage" class="error-message">
      <font-awesome-icon icon="exclamation-circle" />
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="success-message">
      <font-awesome-icon icon="check-circle" />
      {{ successMessage }}
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "VideoUpload",
  data() {
    return {
      selectedFile: null,
      uploadProgress: 0,
      errorMessage: "",
      successMessage: "",
    };
  },
  methods: {
    onFileSelected(event) {
      const file = event.target.files[0];
      if (file) {
        if (file.size > 100 * 1024 * 1024) {
          // 100MB limit
          this.errorMessage = "File is too large. Maximum size is 100MB.";
          return;
        }

        this.selectedFile = file;
        this.errorMessage = "";
        this.successMessage = "";
        this.uploadProgress = 0;
      }
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

    async uploadVideo() {
      if (!this.selectedFile) return;

      this.errorMessage = "";
      this.successMessage = "";
      this.uploadProgress = 0;

      const formData = new FormData();
      formData.append("video", this.selectedFile);

      try {
        const response = await axios.post("/api/videos/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            // Authorization header will be added by the interceptor
          },
          onUploadProgress: (progressEvent) => {
            this.uploadProgress = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total,
            );
          },
        });

        this.successMessage =
          response.data.message || "Video uploaded successfully!";

        // Reset file selection after successful upload
        this.selectedFile = null;
        document.getElementById("video-file").value = "";

        // Reset progress after a delay
        setTimeout(() => {
          this.uploadProgress = 0;
        }, 3000);
      } catch (error) {
        console.error("Upload error:", error);

        // Handle authentication errors
        if (error.response && error.response.status === 401) {
          this.errorMessage = "Your session has expired. Please login again.";
        } else {
          this.errorMessage =
            error.response?.data?.error ||
            "Failed to upload video. Please try again.";
        }

        // Reset progress on error
        this.uploadProgress = 0;
      }
    },
  },
};
</script>

<style scoped>
.video-upload-container {
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
  background-color: var(--white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

h2 {
  color: var(--primary-color);
  text-align: center;
  margin-bottom: 0.5rem;
}

.description {
  text-align: center;
  color: var(--medium);
  margin-bottom: 2rem;
}

.upload-box {
  border: 2px dashed var(--medium);
  border-radius: var(--radius-md);
  padding: 2rem;
  margin-bottom: 1.5rem;
  background-color: var(--light);
}

.file-selection {
  text-align: center;
}

.file-input {
  width: 0.1px;
  height: 0.1px;
  opacity: 0;
  overflow: hidden;
  position: absolute;
  z-index: -1;
}

.file-label {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background-color: var(--primary-color);
  color: white;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-fast) ease;
}

.file-label:hover {
  background-color: var(--primary-dark);
  transform: translateY(-2px);
}

.icon {
  font-size: 1.25rem;
}

.file-info {
  margin-top: 1rem;
  font-size: var(--font-size-sm);
  color: var(--medium);
}

.selected-file-details {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--light);
}

.selected-file-details h3 {
  font-size: var(--font-size-lg);
  margin-bottom: 0.75rem;
  color: var(--dark);
}

.selected-file-details p {
  margin: 0.5rem 0;
  color: var(--medium);
}

.progress-area {
  margin: 1.5rem 0;
}

.progress-bar {
  height: 8px;
  background-color: var(--light);
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background-color: var(--primary-color);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: var(--font-size-sm);
  text-align: right;
  color: var(--medium);
  font-weight: var(--font-weight-medium);
}

.button-area {
  text-align: center;
  margin: 1.5rem 0;
}

.upload-button {
  padding: 0.75rem 2rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-fast) ease;
  min-width: 150px;
}

.upload-button:hover:not(:disabled) {
  background-color: var(--primary-dark);
}

.upload-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.error-message,
.success-message {
  margin-top: 1.5rem;
  padding: 1rem;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.error-message {
  background-color: var(--error-color);
  color: white;
}

.success-message {
  background-color: var(--success-color);
  color: white;
}
</style>
