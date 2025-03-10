package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
)

// GetTransactions retrieves transactions for the current user
func GetTransactions(c *gin.Context) {
	// Get the user ID from the context (set by the auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// In a production app, you would query the database for transactions
	// where the user is either sender or receiver
	// For this implementation, we'll return mock data similar to what the frontend expects

	// Create mock transactions
	transactions := []models.Transaction{
		{
			ID:         1,
			SenderID:   2,
			ReceiverID: userID.(uint),
			Amount:     15,
			Note:       "For JavaScript tutoring",
			CreatedAt:  time.Now().Add(-24 * time.Hour), // 1 day ago
		},
		{
			ID:         2,
			SenderID:   userID.(uint),
			ReceiverID: 3,
			Amount:     5,
			Note:       "For cooking lessons",
			CreatedAt:  time.Now().Add(-48 * time.Hour), // 2 days ago
		},
		{
			ID:         3,
			SenderID:   4,
			ReceiverID: userID.(uint),
			Amount:     10,
			Note:       "For guitar lessons",
			CreatedAt:  time.Now().Add(-72 * time.Hour), // 3 days ago
		},
	}

	// Enhance transactions with sender and receiver names for frontend display
	type EnhancedTransaction struct {
		models.Transaction
		SenderName   string `json:"senderName"`
		ReceiverName string `json:"receiverName"`
	}

	// Create a map of user IDs to names (in a real app, you'd fetch this from the DB)
	userNames := map[uint]string{
		1: "Test User",
		2: "Alice Smith",
		3: "Bob Johnson",
		4: "Carol Williams",
	}

	// Convert to enhanced transactions with names
	enhancedTransactions := make([]EnhancedTransaction, len(transactions))
	for i, tx := range transactions {
		// Get sender name
		senderName, exists := userNames[tx.SenderID]
		if !exists {
			senderName = "Unknown User"
		}

		// Get receiver name
		receiverName, exists := userNames[tx.ReceiverID]
		if !exists {
			receiverName = "Unknown User"
		}

		// Create enhanced transaction
		enhancedTransactions[i] = EnhancedTransaction{
			Transaction:  tx,
			SenderName:   senderName,
			ReceiverName: receiverName,
		}
	}

	c.JSON(http.StatusOK, enhancedTransactions)
}
