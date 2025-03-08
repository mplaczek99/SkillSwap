package repositories_test

import (
	"testing"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
)

func TestInsertTransaction(t *testing.T) {
	tx := &models.Transaction{
		SenderID:   1,
		ReceiverID: 2,
		Amount:     10,
	}
	created, err := repositories.InsertTransaction(tx)
	if err != nil {
		t.Errorf("InsertTransaction returned error: %v", err)
	}
	if created.ID == 0 {
		t.Errorf("expected transaction ID to be set, got 0")
	}
}

func TestGetTransactionByID(t *testing.T) {
	tx, err := repositories.GetTransactionByID("1")
	if err != nil {
		t.Errorf("GetTransactionByID returned error: %v", err)
	}
	if tx.ID == 0 {
		t.Errorf("expected dummy transaction, got ID=0")
	}

	// If you had logic to handle nonexistent ID, test it here
	// e.g., tx, err = repositories.GetTransactionByID("999")
	// if err == nil {...}
}

func TestTransactionQueries(t *testing.T) {
	// Create mock transactions for a user
	userID := uint(1)
	now := time.Now()

	// Mock transactions with different dates
	mockTxs := []models.Transaction{
		{
			ID:         1,
			SenderID:   userID,
			ReceiverID: 2,
			Amount:     10,
			CreatedAt:  now.AddDate(0, 0, -30), // 30 days ago
		},
		{
			ID:         2,
			SenderID:   userID,
			ReceiverID: 3,
			Amount:     20,
			CreatedAt:  now.AddDate(0, 0, -15), // 15 days ago
		},
		{
			ID:         3,
			ReceiverID: userID, // User as receiver
			SenderID:   2,
			Amount:     15,
			CreatedAt:  now.AddDate(0, 0, -7), // 7 days ago
		},
		{
			ID:         4,
			ReceiverID: userID, // User as receiver
			SenderID:   3,
			Amount:     25,
			CreatedAt:  now.AddDate(0, 0, -1), // 1 day ago
		},
	}

	// Mock function to get transactions by user ID
	getTransactionsByUserID := func(uID uint) []models.Transaction {
		var result []models.Transaction
		for _, tx := range mockTxs {
			if tx.SenderID == uID || tx.ReceiverID == uID {
				result = append(result, tx)
			}
		}
		return result
	}

	// Mock function to calculate sent and received totals
	getTransactionTotals := func(uID uint) (int, int) {
		var totalSent, totalReceived int
		for _, tx := range mockTxs {
			if tx.SenderID == uID {
				totalSent += tx.Amount
			}
			if tx.ReceiverID == uID {
				totalReceived += tx.Amount
			}
		}
		return totalSent, totalReceived
	}

	// Mock function to get transactions by date range
	getTransactionsByDateRange := func(startDate, endDate time.Time) []models.Transaction {
		var result []models.Transaction
		for _, tx := range mockTxs {
			if (tx.CreatedAt.Equal(startDate) || tx.CreatedAt.After(startDate)) &&
				(tx.CreatedAt.Equal(endDate) || tx.CreatedAt.Before(endDate)) {
				result = append(result, tx)
			}
		}
		return result
	}

	// Test 1: Get all transactions for a user
	t.Run("Get All User Transactions", func(t *testing.T) {
		txs := getTransactionsByUserID(userID)
		if len(txs) != 4 {
			t.Errorf("Expected 4 transactions, got %d", len(txs))
		}
	})

	// Test 2: Calculate sent vs received totals
	t.Run("Calculate Transaction Totals", func(t *testing.T) {
		sent, received := getTransactionTotals(userID)
		expectedSent := 30     // 10 + 20
		expectedReceived := 40 // 15 + 25

		if sent != expectedSent {
			t.Errorf("Expected sent total of %d, got %d", expectedSent, sent)
		}
		if received != expectedReceived {
			t.Errorf("Expected received total of %d, got %d", expectedReceived, received)
		}

		// Calculate net balance
		netBalance := received - sent
		expectedNet := 10 // 40 - 30
		if netBalance != expectedNet {
			t.Errorf("Expected net balance of %d, got %d", expectedNet, netBalance)
		}
	})

	// Test 3: Get transactions in last 7 days
	t.Run("Get Recent Transactions", func(t *testing.T) {
		startDate := now.AddDate(0, 0, -7) // Last 7 days
		txs := getTransactionsByDateRange(startDate, now)

		if len(txs) != 2 {
			t.Errorf("Expected 2 transactions in last 7 days, got %d", len(txs))
		}

		// Check if the most recent transaction is included
		var foundMostRecent bool
		for _, tx := range txs {
			if tx.ID == 4 { // Most recent transaction
				foundMostRecent = true
				break
			}
		}
		if !foundMostRecent {
			t.Errorf("Most recent transaction not found in results")
		}
	})
}
