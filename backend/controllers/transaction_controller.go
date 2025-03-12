package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/utils"
	"gorm.io/gorm"
)

type CreateTransactionRequest struct {
	RecipientEmail string `json:"recipientEmail" binding:"required"`
	Amount         int    `json:"amount" binding:"required,min=1"`
	Note           string `json:"note"`
}

// GetTransactions retrieves transactions for the current user
func GetTransactions(c *gin.Context) {
	// Get the user ID from the context (set by the auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	// Initialize repositories
	txRepo := repositories.NewTransactionRepository(db.(*gorm.DB))
	userRepo := repositories.NewUserRepository(db.(*gorm.DB))

	// Get transactions
	transactions, err := txRepo.GetTransactionsByUserID(userID.(uint))
	if err != nil {
		utils.Error("Failed to retrieve transactions: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	// Create a map of user IDs to names (fetch from database)
	userNames := make(map[uint]string)
	currentUserID := userID.(uint)

	// Enhanced response format for frontend
	type EnhancedTransaction struct {
		models.Transaction
		SenderName   string `json:"senderName"`
		ReceiverName string `json:"receiverName"`
	}

	enhancedTransactions := make([]EnhancedTransaction, len(transactions))
	for i, tx := range transactions {
		// Get sender name if not cached
		if _, exists := userNames[tx.SenderID]; !exists {
			if tx.SenderID == currentUserID {
				userNames[tx.SenderID] = "You" // Special case for current user
			} else {
				user, err := userRepo.GetUserByID(tx.SenderID)
				if err == nil {
					userNames[tx.SenderID] = user.Name
				} else {
					userNames[tx.SenderID] = "Unknown User"
				}
			}
		}

		// Get receiver name if not cached
		if _, exists := userNames[tx.ReceiverID]; !exists {
			if tx.ReceiverID == currentUserID {
				userNames[tx.ReceiverID] = "You" // Special case for current user
			} else {
				user, err := userRepo.GetUserByID(tx.ReceiverID)
				if err == nil {
					userNames[tx.ReceiverID] = user.Name
				} else {
					userNames[tx.ReceiverID] = "Unknown User"
				}
			}
		}

		// Create enhanced transaction
		enhancedTransactions[i] = EnhancedTransaction{
			Transaction:  tx,
			SenderName:   userNames[tx.SenderID],
			ReceiverName: userNames[tx.ReceiverID],
		}
	}

	c.JSON(http.StatusOK, enhancedTransactions)
}

// CreateTransaction handles the creation of a new transaction
func CreateTransaction(c *gin.Context) {
	// Get the sender ID from the context (set by the auth middleware)
	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get database from context
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.(*gorm.DB))
	txRepo := repositories.NewTransactionRepository(db.(*gorm.DB))

	// Find recipient by email
	recipient, err := userRepo.GetUserByEmail(req.RecipientEmail)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	// Check that sender is not recipient
	if senderID.(uint) == recipient.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send points to yourself"})
		return
	}

	// Create transaction
	transaction := models.Transaction{
		SenderID:   senderID.(uint),
		ReceiverID: recipient.ID,
		Amount:     req.Amount,
		Note:       req.Note,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save transaction and update balances
	if err := txRepo.CreateTransaction(&transaction); err != nil {
		utils.Error("Failed to create transaction: " + err.Error())

		// Return specific error for insufficient funds
		if err.Error() == "insufficient skill points" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have enough SkillPoints"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Return successful response with transaction details
	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction completed successfully",
		"transaction": gin.H{
			"id":           transaction.ID,
			"sender_id":    transaction.SenderID,
			"receiver_id":  transaction.ReceiverID,
			"amount":       transaction.Amount,
			"note":         transaction.Note,
			"created_at":   transaction.CreatedAt,
			"senderName":   "You",
			"receiverName": recipient.Name,
		},
	})
}
