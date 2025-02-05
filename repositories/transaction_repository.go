package repositories

import (
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/models"
)

// InsertTransaction creates a new transaction record.
func InsertTransaction(tx *models.Transaction) (*models.Transaction, error) {
	result := config.DB.Create(tx)
	return tx, result.Error
}

// GetTransactionByID retrieves a transaction by its ID.
func GetTransactionByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	result := config.DB.Where("id = ?", id).First(&transaction)
	return &transaction, result.Error
}

