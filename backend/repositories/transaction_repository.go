package repositories

import (
	"errors"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/gorm"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	DB *gorm.DB
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

// GetTransactionsByUserID returns all transactions for a given user
func (r *TransactionRepository) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

// CreateTransaction creates a new transaction and updates both users' balances
func (r *TransactionRepository) CreateTransaction(tx *models.Transaction) error {
	// Use a database transaction to ensure data consistency
	return r.DB.Transaction(func(dbTx *gorm.DB) error {
		// Find sender
		var sender models.User
		if err := dbTx.Where("id = ?", tx.SenderID).First(&sender).Error; err != nil {
			return errors.New("sender not found")
		}

		// Find receiver
		var receiver models.User
		if err := dbTx.Where("id = ?", tx.ReceiverID).First(&receiver).Error; err != nil {
			return errors.New("receiver not found")
		}

		// Check if sender has enough points
		if sender.SkillPoints < tx.Amount {
			return errors.New("insufficient skill points")
		}

		// Update sender's balance
		if err := dbTx.Model(&sender).Update("skill_points", sender.SkillPoints-tx.Amount).Error; err != nil {
			return err
		}

		// Update receiver's balance
		if err := dbTx.Model(&receiver).Update("skill_points", receiver.SkillPoints+tx.Amount).Error; err != nil {
			return err
		}

		// Set creation timestamp
		tx.CreatedAt = time.Now()
		tx.UpdatedAt = time.Now()

		// Create the transaction record
		if err := dbTx.Create(tx).Error; err != nil {
			return err
		}

		return nil
	})
}

// For backward compatibility with existing code
func InsertTransaction(tx *models.Transaction) (*models.Transaction, error) {
	// In a real implementation, you would use the repository pattern
	// For now, assign a dummy ID
	tx.ID = 1
	return tx, nil
}

// For backward compatibility with existing code
func GetTransactionByID(id string) (*models.Transaction, error) {
	// Dummy implementation
	dummyTx := models.Transaction{
		ID:         1,
		SenderID:   1,
		ReceiverID: 2,
		Amount:     10,
	}
	return &dummyTx, nil
}
