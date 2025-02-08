package repositories_test

import (
    "testing"

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

