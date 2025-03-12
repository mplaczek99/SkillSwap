package repositories

import (
	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/gorm"
)

// JobRepository handles database operations for jobs
type JobRepository struct {
	DB *gorm.DB
}

// NewJobRepository creates a new instance of JobRepository
func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{DB: db}
}

// GetAllJobs returns all job postings
func (r *JobRepository) GetAllJobs() ([]models.Job, error) {
	var jobs []models.Job
	err := r.DB.Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

// GetJobByID returns a job posting by ID
func (r *JobRepository) GetJobByID(id uint) (*models.Job, error) {
	var job models.Job
	err := r.DB.First(&job, id).Error
	return &job, err
}

// CreateJob creates a new job posting
func (r *JobRepository) CreateJob(job *models.Job) error {
	return r.DB.Create(job).Error
}

// UpdateJob updates an existing job posting
func (r *JobRepository) UpdateJob(job *models.Job) error {
	return r.DB.Save(job).Error
}

// DeleteJob deletes a job posting
func (r *JobRepository) DeleteJob(id uint) error {
	return r.DB.Delete(&models.Job{}, id).Error
}

// GetJobsByUser returns all job postings by a specific user
func (r *JobRepository) GetJobsByUser(userID uint) ([]models.Job, error) {
	var jobs []models.Job
	err := r.DB.Where("posted_by_user_id = ?", userID).Find(&jobs).Error
	return jobs, err
}
