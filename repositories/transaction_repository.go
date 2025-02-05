package repositories

import (
    "github.com/mplaczek99/SkillSwap/models"
)

// InsertTransaction returns the transaction with a dummy ID.
func InsertTransaction(tx *models.Transaction) (*models.Transaction, error) {
    tx.ID = 1
    return tx, nil
}

// GetTransactionByID returns a dummy transaction.
func GetTransactionByID(id string) (*models.Transaction, error) {
    dummyTx := models.Transaction{
        ID:         1,
        SenderID:   1,
        ReceiverID: 2,
        Amount:     10,
    }
    return &dummyTx, nil
}

