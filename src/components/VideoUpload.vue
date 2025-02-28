<template>
  <div class="video-upload">
    <h2>Upload a Video</h2>
    <input type="file" @change="onFileSelected" accept="video/*" />
    <button @click="uploadVideo" :disabled="!selectedFile">Upload</button>
    <div v-if="uploadProgress">Upload Progress: {{ uploadProgress }}%</div>
    <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
    <div v-if="successMessage" class="success">{{ successMessage }}</div>
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
      this.selectedFile = event.target.files[0];
    },
    async uploadVideo() {
      if (!this.selectedFile) return;
      const formData = new FormData();
      formData.append("video", this.selectedFile);
      try {
        const response = await axios.post("/api/videos/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
          onUploadProgress: (progressEvent) => {
            this.uploadProgress = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total,
            );
          },
        });
        this.successMessage = response.data.message;
      } catch (error) {
        this.errorMessage = error.response?.data?.error || "Upload failed";
      }
    },
  },
};
</script>

<style scoped>
.video-upload {
  max-width: 500px;
  margin: 2rem auto;
  padding: 1rem;
  border: 1px solid var(--medium);
  border-radius: var(--radius-md);
}
.error {
  color: var(--error-color);
  margin-top: var(--space-2);
}
.success {
  color: var(--success-color);
  margin-top: var(--space-2);
}
</style>
