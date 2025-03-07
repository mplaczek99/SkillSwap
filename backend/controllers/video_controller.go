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
func GetVideosList(c *gin.Context) {
	// Define the directory where uploads are stored.
	uploadDir := "./uploads"

	// Initialize an empty videos slice to ensure we always return a JSON array
	videos := []gin.H{}

	// Ensure directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		// Return empty array if directory doesn't exist
		c.JSON(http.StatusOK, videos)
		return
	}

	// Read the directory
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read videos directory"})
		return
	}

	// Collect video information
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
			// Skip this file if we can't get info
			continue
		}

		// Add to our results
		videos = append(videos, gin.H{
			"id":           filename, // Using filename as ID
			"name":         filename,
			"thumbnail":    thumbnailPath,
			"hasThumbnail": hasThumbnail,
			"size":         fileInfo.Size(),
			"uploadedAt":   fileInfo.ModTime(),
		})
	}

	c.JSON(http.StatusOK, videos)
}
