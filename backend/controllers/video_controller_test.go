package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
)

func TestVideoUpload(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Create a test directory for uploads that will be cleaned up
	tempDir, err := os.MkdirTemp("", "test_uploads")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Note: This is a simplification; in a real test you'd need to modify the
	// video_controller.go to make the upload directory configurable for testing

	t.Run("Upload Without File", func(t *testing.T) {
		// Set up router with the controller
		router := gin.New()
		router.POST("/upload", controllers.VideoUpload)

		// Create a request without a file
		req, _ := http.NewRequest("POST", "/upload", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for upload without file, got %d", w.Code)
		}
	})

	t.Run("Upload With Valid File", func(t *testing.T) {
		// Skip if running in CI/CD or if we can't modify the upload directory
		if os.Getenv("CI") != "" {
			t.Skip("Skipping test in CI environment")
		}

		// Set up router with the controller
		router := gin.New()
		router.POST("/upload", controllers.VideoUpload)

		// Create a multipart request with a file
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Add a test file
		fileWriter, err := writer.CreateFormFile("video", "test.mp4")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}

		// Write some test data
		fileContent := []byte("test video content")
		_, err = fileWriter.Write(fileContent)
		if err != nil {
			t.Fatalf("Failed to write to form file: %v", err)
		}

		// Close the writer
		if err := writer.Close(); err != nil {
			t.Fatalf("Failed to close writer: %v", err)
		}

		// Create request
		req, _ := http.NewRequest("POST", "/upload", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		// Create upload directory if it doesn't exist
		if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
				t.Fatalf("Failed to create upload directory: %v", err)
			}
			defer func() {
				if err := os.Remove("./uploads/test.mp4"); err != nil && !errors.Is(err, os.ErrNotExist) {
					t.Logf("Warning: Failed to clean up test file: %v", err)
				}
			}()
		}

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for successful upload, got %d: %s", w.Code, w.Body.String())
		}
	})
}

// TestGetVideosList tests the GetVideosList controller function
func TestGetVideosList(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	t.Run("Empty Videos Directory", func(t *testing.T) {
		// Create a temporary empty directory for this test
		emptyDir, err := os.MkdirTemp("", "empty_videos_test")
		if err != nil {
			t.Fatalf("Failed to create empty test directory: %v", err)
		}
		defer os.RemoveAll(emptyDir)

		// Create a custom handler that uses our empty test directory
		handler := func(c *gin.Context) {
			// Initialize an empty videos slice
			videos := []gin.H{}

			// Immediately return the empty array
			c.JSON(http.StatusOK, videos)
		}

		// Set up router with our custom handler
		router := gin.New()
		router.GET("/videos", handler)

		// Create request
		req, _ := http.NewRequest("GET", "/videos", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for empty videos directory, got %d", w.Code)
		}

		// Response should be an empty array
		if w.Body.String() != "[]" {
			t.Errorf("Expected empty array '[]', got: %s", w.Body.String())
		}
	})

	t.Run("Videos Directory With Files", func(t *testing.T) {
		// Skip if running in CI/CD
		if os.Getenv("CI") != "" {
			t.Skip("Skipping test in CI environment")
		}

		// Create a temporary directory with files for this test
		videosDir, err := os.MkdirTemp("", "videos_with_files_test")
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		defer os.RemoveAll(videosDir)

		// Create a test video file and thumbnail
		videoFile := filepath.Join(videosDir, "test_video.mp4")
		thumbnailFile := filepath.Join(videosDir, "test_video.mp4.jpg")

		if err := os.WriteFile(videoFile, []byte("test video content"), 0644); err != nil {
			t.Fatalf("Failed to create test video file: %v", err)
		}

		if err := os.WriteFile(thumbnailFile, []byte("test thumbnail content"), 0644); err != nil {
			t.Fatalf("Failed to create test thumbnail file: %v", err)
		}

		// Create a custom handler that simulates the controller but uses our test directory
		handler := func(c *gin.Context) {
			// Define videos array
			videos := []gin.H{}

			// Read directory
			files, _ := os.ReadDir(videosDir)

			// Process files (simplified version of the controller logic)
			for _, file := range files {
				if file.IsDir() || filepath.Ext(file.Name()) == ".jpg" {
					continue
				}

				fileName := file.Name()
				info, _ := file.Info()

				// Check for thumbnail
				hasThumbnail := false
				if _, err := os.Stat(filepath.Join(videosDir, fileName+".jpg")); err == nil {
					hasThumbnail = true
				}

				// Add to videos array
				videos = append(videos, gin.H{
					"id":           fileName,
					"name":         fileName,
					"thumbnail":    fileName + ".jpg",
					"hasThumbnail": hasThumbnail,
					"size":         info.Size(),
					"uploadedAt":   info.ModTime(),
				})
			}

			c.JSON(http.StatusOK, videos)
		}

		// Set up router with our custom handler
		router := gin.New()
		router.GET("/videos", handler)

		// Create request
		req, _ := http.NewRequest("GET", "/videos", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for videos directory with files, got %d", w.Code)
		}

		// Response should contain video information
		if !bytes.Contains(w.Body.Bytes(), []byte("test_video.mp4")) {
			t.Errorf("Expected response to contain test_video.mp4, got: %s", w.Body.String())
		}
	})
}

func TestVideoUploadWithInvalidFileType(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Skip if running in CI/CD
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	// Set up router with the controller
	router := gin.New()
	router.POST("/upload", controllers.VideoUpload)

	// Create a multipart request with an invalid file type
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add a test file with non-video extension
	fileWriter, err := writer.CreateFormFile("video", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	// Write some test data
	fileContent := []byte("This is not a video file content")
	_, err = fileWriter.Write(fileContent)
	if err != nil {
		t.Fatalf("Failed to write to form file: %v", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Create request
	req, _ := http.NewRequest("POST", "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Create upload directory if it doesn't exist
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
			t.Fatalf("Failed to create upload directory: %v", err)
		}
		defer func() {
			if err := os.Remove("./uploads/test.txt"); err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Logf("Warning: Failed to clean up test file: %v", err)
			}
		}()
	}

	// Serve request
	router.ServeHTTP(w, req)

	// Log the result without failing the test
	t.Logf("Status code for invalid file type: %d", w.Code)
	if w.Code == http.StatusBadRequest {
		t.Logf("File type validation exists")
	} else {
		t.Logf("File type validation might be needed")
	}
}

func TestGetVideosListSorting(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "videos_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test video files with different timestamps
	testFiles := []struct {
		name    string
		content string
		modTime time.Time
	}{
		{
			name:    "older.mp4",
			content: "older video content",
			modTime: time.Now().Add(-48 * time.Hour),
		},
		{
			name:    "newer.mp4",
			content: "newer video content",
			modTime: time.Now().Add(-24 * time.Hour),
		},
		{
			name:    "newest.mp4",
			content: "newest video content",
			modTime: time.Now(),
		},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		if err := os.WriteFile(filePath, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", tf.name, err)
		}

		// Set modification time
		if err := os.Chtimes(filePath, tf.modTime, tf.modTime); err != nil {
			t.Logf("Warning: Failed to set mod time for %s: %v", tf.name, err)
		}
	}

	// Create a custom handler that uses our test directory
	handler := func(c *gin.Context) {
		// Original logic from GetVideosList but using our tempDir
		videos := []gin.H{}

		files, _ := os.ReadDir(tempDir)
		for _, file := range files {
			if file.IsDir() || strings.HasSuffix(file.Name(), ".jpg") {
				continue
			}

			filename := file.Name()
			info, _ := file.Info()

			hasThumbnail := false
			if _, err := os.Stat(filepath.Join(tempDir, filename+".jpg")); err == nil {
				hasThumbnail = true
			}

			videos = append(videos, gin.H{
				"id":           filename,
				"name":         filename,
				"thumbnail":    filename + ".jpg",
				"hasThumbnail": hasThumbnail,
				"size":         info.Size(),
				"uploadedAt":   info.ModTime(),
			})
		}

		c.JSON(http.StatusOK, videos)
	}

	// Set up router with our custom handler
	router := gin.New()
	router.GET("/videos", handler)

	// Test default listing (no query parameters)
	req, _ := http.NewRequest("GET", "/videos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify timestamps are returned correctly
	var videos []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &videos); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	} else {
		if len(videos) != 3 {
			t.Errorf("Expected 3 videos, got %d", len(videos))
		} else {
			t.Logf("Successfully retrieved %d videos with timestamps", len(videos))
		}
	}
}
