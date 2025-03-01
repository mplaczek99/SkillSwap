package controllers

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// VideoUpload handles uploading and processing of video files.
// @Summary Upload a video
// @Description Upload a video file and generate a thumbnail
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param video formData file true "Video file to upload"
// @Success 200 {object} map[string]string "Successfully uploaded video"
// @Failure 400 {object} map[string]string "Video file is required"
// @Failure 500 {object} map[string]string "Failed to save video"
// @Security BearerAuth
// @Router /videos/upload [post]
func VideoUpload(c *gin.Context) {
	// Retrieve the file from the form data.
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file is required"})
		return
	}

	// Define a directory to store uploads.
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}
	filePath := filepath.Join(uploadDir, file.Filename)

	// Save the uploaded file.
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
		return
	}

	// Start processing asynchronously (e.g., generate a thumbnail).
	go processVideo(filePath)

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully", "file": file.Filename})
}

// processVideo demonstrates invoking FFmpeg to generate a thumbnail image.
func processVideo(filePath string) {
	// Example: generate a thumbnail 1 second into the video.
	thumbnailPath := filePath + ".jpg"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-ss", "00:00:01.000", "-vframes", "1", thumbnailPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Video processing failed for %s: %v", filePath, err)
		return
	}
	log.Printf("Thumbnail created: %s", thumbnailPath)
	// Additional processing such as transcoding can be added here.
}

// GetVideosList handles retrieving a list of all uploaded videos.
// @Summary Get videos list
// @Description Retrieve a list of all uploaded videos
// @Tags videos
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of videos"
// @Failure 500 {object} map[string]string "Failed to read videos directory"
// @Security BearerAuth
// @Router /videos [get]
func GetVideosList(c *gin.Context) {
	// Define the directory where uploads are stored.
	uploadDir := "./uploads"

	// Ensure directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Read the directory
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read videos directory"})
		return
	}

	// Collect video information
	var videos []gin.H

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		filename := file.Name()

		// Skip thumbnails (files ending with .jpg)
		if strings.HasSuffix(filename, ".jpg") {
			continue
		}

		// Check if a thumbnail exists for this video
		thumbnailPath := filename + ".jpg"
		thumbnailFullPath := filepath.Join(uploadDir, thumbnailPath)
		hasThumbnail := false
		if _, err := os.Stat(thumbnailFullPath); err == nil {
			hasThumbnail = true
		}

		// Get file info for additional metadata
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		videos = append(videos, gin.H{
			"id":           filename,
			"name":         filename,
			"size":         fileInfo.Size(),
			"uploadedAt":   fileInfo.ModTime(),
			"hasThumbnail": hasThumbnail,
		})
	}

	c.JSON(http.StatusOK, videos)
}
