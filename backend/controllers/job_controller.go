package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/utils"
	"gorm.io/gorm"
)

// GetJobs handles retrieving all job postings
func GetJobs(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	jobRepo := repositories.NewJobRepository(db.(*gorm.DB))
	jobs, err := jobRepo.GetAllJobs()
	if err != nil {
		utils.Error("Failed to retrieve jobs: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// GetJob handles retrieving a specific job posting by ID
func GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	jobRepo := repositories.NewJobRepository(db.(*gorm.DB))
	job, err := jobRepo.GetJobByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// CreateJob handles creating a new job posting
func CreateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job data"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	job.PostedByUserID = userID.(uint)

	// Get user name if available
	email, emailExists := c.Get("email")
	if emailExists {
		// Get user from database using email
		db, dbExists := c.Get("db")
		if dbExists {
			userRepo := repositories.NewUserRepository(db.(*gorm.DB))
			user, err := userRepo.GetUserByEmail(email.(string))
			if err == nil {
				job.PostedByName = user.Name
			}
		}
	}

	// Set timestamps
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	jobRepo := repositories.NewJobRepository(db.(*gorm.DB))
	if err := jobRepo.CreateJob(&job); err != nil {
		utils.Error("Failed to create job posting: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job posting"})
		return
	}

	c.JSON(http.StatusCreated, job)
}

// UpdateJob handles updating an existing job posting
func UpdateJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var jobUpdates models.Job
	if err := c.ShouldBindJSON(&jobUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job data"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	jobRepo := repositories.NewJobRepository(db.(*gorm.DB))

	// Get existing job
	existingJob, err := jobRepo.GetJobByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		return
	}

	// Check if user is the owner of the job posting
	if existingJob.PostedByUserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this job posting"})
		return
	}

	// Update fields
	jobUpdates.ID = uint(id)
	jobUpdates.PostedByUserID = userID.(uint)
	jobUpdates.PostedByName = existingJob.PostedByName
	jobUpdates.CreatedAt = existingJob.CreatedAt
	jobUpdates.UpdatedAt = time.Now()

	if err := jobRepo.UpdateJob(&jobUpdates); err != nil {
		utils.Error("Failed to update job posting: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job posting"})
		return
	}

	c.JSON(http.StatusOK, jobUpdates)
}

// DeleteJob handles deleting a job posting
func DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	jobRepo := repositories.NewJobRepository(db.(*gorm.DB))

	// Get existing job
	existingJob, err := jobRepo.GetJobByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		return
	}

	// Check if user is the owner of the job posting
	if existingJob.PostedByUserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this job posting"})
		return
	}

	if err := jobRepo.DeleteJob(uint(id)); err != nil {
		utils.Error("Failed to delete job posting: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job posting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job posting deleted successfully"})
}
