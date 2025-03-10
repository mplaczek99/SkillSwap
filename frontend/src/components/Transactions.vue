<template>
  <div class="transactions-container">
    <div class="container">
      <h2>SkillPoints Transactions</h2>
      <p class="description">Track your SkillPoints earnings and spendings</p>

      <!-- Summary Card -->
      <div class="card summary-card">
        <div class="points-display">
          <div class="points-icon">
            <font-awesome-icon icon="star" />
          </div>
          <div class="points-info">
            <h3>Your Balance</h3>
            <p class="balance">{{ user.skillPoints || 0 }} SkillPoints</p>
          </div>
        </div>
        <div class="summary-stats">
          <div class="stat-item">
            <span class="stat-label">Earned</span>
            <span class="stat-value positive">+{{ totalEarned }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Spent</span>
            <span class="stat-value negative">-{{ totalSpent }}</span>
          </div>
        </div>
      </div>

      <!-- Transaction Actions -->
      <div class="action-buttons">
        <button class="btn btn-primary" @click="showSendModal = true">
          <font-awesome-icon icon="paper-plane" />
          Send SkillPoints
        </button>
        <button class="btn btn-outline" @click="fetchTransactions">
          <font-awesome-icon icon="sync" />
          Refresh
        </button>
      </div>

      <!-- Transactions List -->
      <div class="card transactions-card">
        <h3>Recent Transactions</h3>

        <div v-if="loading" class="loading-state">
          <font-awesome-icon icon="spinner" class="spin" />
          <p>Loading transactions...</p>
        </div>

        <div v-else-if="error" class="error-message">
          <font-awesome-icon icon="exclamation-circle" />
          {{ error }}
        </div>

        <div v-else-if="transactions.length === 0" class="empty-state">
          <font-awesome-icon icon="history" class="empty-icon" />
          <p>No transactions yet</p>
          <p class="empty-hint">
            Earn SkillPoints by teaching skills to others!
          </p>
        </div>

        <div v-else class="transactions-list">
          <div v-for="(transaction, index) in transactions" :key="index" class="transaction-item"
            :class="getTransactionClass(transaction)">
            <div class="transaction-icon">
              <font-awesome-icon :icon="getTransactionIcon(transaction)" />
            </div>
            <div class="transaction-details">
              <p class="transaction-title">
                {{ getTransactionTitle(transaction) }}
              </p>
              <p class="transaction-date">
                {{ formatDate(transaction.createdAt) }}
              </p>
            </div>
            <div class="transaction-amount" :class="getAmountClass(transaction)">
              {{ getAmountPrefix(transaction) }}{{ transaction.amount }} SP
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Send Points Modal -->
    <div v-if="showSendModal" class="modal-backdrop">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Send SkillPoints</h3>
          <button class="close-button" @click="showSendModal = false">
            <font-awesome-icon icon="times" />
          </button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="sendPoints">
            <div class="form-group">
              <label for="recipient">Recipient Email</label>
              <input id="recipient" v-model="sendForm.recipientEmail" type="email" placeholder="Enter recipient's email"
                required />
            </div>
            <div class="form-group">
              <label for="amount">Amount</label>
              <input id="amount" v-model.number="sendForm.amount" type="number" min="1" :max="user.skillPoints || 0"
                placeholder="Enter amount to send" required />
            </div>
            <div class="form-group">
              <label for="note">Note (Optional)</label>
              <textarea id="note" v-model="sendForm.note" placeholder="What are these points for?"></textarea>
            </div>
            <div class="form-actions">
              <button type="button" class="btn btn-outline" @click="showSendModal = false">
                Cancel
              </button>
              <button type="submit" class="btn btn-primary" :disabled="sendForm.amount <= 0 ||
                sendForm.amount > (user.skillPoints || 0)
                ">
                Send
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import axios from "axios";

export default {
  name: "Transactions",
  data() {
    return {
      transactions: [],
      loading: true,
      error: null,
      showSendModal: false,
      sendForm: {
        recipientEmail: "",
        amount: 1,
        note: "",
      },
      totalEarned: 0,
      totalSpent: 0,
    };
  },
  computed: {
    ...mapGetters(["user"]),
  },
  created() {
    this.fetchTransactions();
  },
  methods: {
    async fetchTransactions() {
      this.loading = true;
      this.error = null;

      try {
        // Call the real API endpoint
        const response = await axios.get("/api/transactions");
        this.transactions = response.data;
        this.calculateTotals();
      } catch (error) {
        console.error("Error fetching transactions:", error);
        this.error = "Failed to load transactions. Please try again.";
      } finally {
        this.loading = false;
      }
    },

    // Rest of the methods stay the same
    calculateTotals() {
      const userId = this.user ? this.user.id : 1;

      this.totalEarned = this.transactions
        .filter((t) => t.receiverId === userId)
        .reduce((sum, t) => sum + t.amount, 0);

      this.totalSpent = this.transactions
        .filter((t) => t.senderId === userId)
        .reduce((sum, t) => sum + t.amount, 0);
    },

    async sendPoints() {
      // In a real app, send to API
      alert(`Points would be sent to ${this.sendForm.recipientEmail}`);

      // Mock successful transaction
      const newTransaction = {
        id: this.transactions.length + 1,
        senderId: this.user ? this.user.id : 1,
        receiverId: 999, // Placeholder
        amount: this.sendForm.amount,
        createdAt: new Date(),
        senderName: this.user ? this.user.name : "You",
        receiverName: this.sendForm.recipientEmail,
        note: this.sendForm.note,
      };

      this.transactions.unshift(newTransaction);
      this.calculateTotals();

      // Reset form and close modal
      this.sendForm = {
        recipientEmail: "",
        amount: 1,
        note: "",
      };
      this.showSendModal = false;

      // In a real app:
      // try {
      //   await axios.post('/api/transactions', this.sendForm);
      //   this.fetchTransactions();
      //   this.sendForm = { recipientEmail: '', amount: 1, note: '' };
      //   this.showSendModal = false;
      // } catch (error) {
      //   console.error('Error sending points:', error);
      //   this.error = 'Failed to send points. Please try again.';
      // }
    },

    getTransactionClass(transaction) {
      const userId = this.user ? this.user.id : 1;
      return {
        received: transaction.receiverId === userId,
        sent: transaction.senderId === userId,
      };
    },

    getTransactionIcon(transaction) {
      const userId = this.user ? this.user.id : 1;
      return transaction.receiverId === userId ? "arrow-down" : "arrow-up";
    },

    getTransactionTitle(transaction) {
      const userId = this.user ? this.user.id : 1;

      if (transaction.receiverId === userId) {
        return `Received from ${transaction.senderName}${transaction.note ? ": " + transaction.note : ""}`;
      } else {
        return `Sent to ${transaction.receiverName}${transaction.note ? ": " + transaction.note : ""}`;
      }
    },

    getAmountClass(transaction) {
      const userId = this.user ? this.user.id : 1;
      return {
        positive: transaction.receiverId === userId,
        negative: transaction.senderId === userId,
      };
    },

    getAmountPrefix(transaction) {
      const userId = this.user ? this.user.id : 1;
      return transaction.receiverId === userId ? "+" : "-";
    },

    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    },
  },
};
</script>

<style scoped>
.transactions-container {
  padding-bottom: var(--space-12);
}

h2 {
  color: var(--primary-color);
  text-align: center;
  margin-bottom: var(--space-2);
  font-size: var(--font-size-3xl);
}

.description {
  text-align: center;
  color: var(--medium);
  margin-bottom: var(--space-8);
  font-size: var(--font-size-lg);
}

.summary-card {
  padding: var(--space-6);
  margin-bottom: var(--space-6);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: var(--space-4);
}

.points-display {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.points-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: var(--primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-2xl);
}

.points-info h3 {
  margin: 0;
  font-size: var(--font-size-lg);
  color: var(--medium);
}

.balance {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--dark);
  margin: var(--space-1) 0 0 0;
}

.summary-stats {
  display: flex;
  gap: var(--space-6);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--medium);
}

.stat-value {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
}

.action-buttons {
  display: flex;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.transactions-card {
  padding: var(--space-6);
}

.transactions-card h3 {
  margin-top: 0;
  margin-bottom: var(--space-4);
  font-size: var(--font-size-xl);
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-8) 0;
  color: var(--medium);
}

.spin {
  animation: spin 1s linear infinite;
  font-size: var(--font-size-2xl);
  margin-bottom: var(--space-2);
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.empty-icon {
  font-size: var(--font-size-3xl);
  opacity: 0.5;
  margin-bottom: var(--space-2);
}

.empty-hint {
  font-style: italic;
  font-size: var(--font-size-sm);
  margin-top: var(--space-2);
}

.error-message {
  background-color: var(--error-color);
  color: white;
  padding: var(--space-4);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.transactions-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.transaction-item {
  display: flex;
  align-items: center;
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background-color: var(--light);
  gap: var(--space-3);
}

.transaction-item.received {
  border-left: 4px solid var(--success-color);
}

.transaction-item.sent {
  border-left: 4px solid var(--info-color);
}

.transaction-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-sm);
  flex-shrink: 0;
}

.received .transaction-icon {
  background-color: var(--success-color);
  color: white;
}

.sent .transaction-icon {
  background-color: var(--info-color);
  color: white;
}

.transaction-details {
  flex: 1;
  min-width: 0;
}

.transaction-title {
  margin: 0;
  font-weight: var(--font-weight-medium);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.transaction-date {
  margin: var(--space-1) 0 0 0;
  font-size: var(--font-size-xs);
  color: var(--medium);
}

.transaction-amount {
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-lg);
}

.transaction-amount.positive {
  color: var(--success-color);
}

.transaction-amount.negative {
  color: var(--info-color);
}

/* Modal styles */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background-color: var(--white);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 500px;
  box-shadow: var(--shadow-xl);
}

.modal-header {
  padding: var(--space-4) var(--space-6);
  border-bottom: 1px solid var(--light);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h3 {
  margin: 0;
  font-size: var(--font-size-xl);
}

.close-button {
  background: none;
  border: none;
  font-size: var(--font-size-lg);
  color: var(--medium);
  cursor: pointer;
}

.modal-body {
  padding: var(--space-6);
}

.form-group {
  margin-bottom: var(--space-4);
}

.form-group label {
  display: block;
  margin-bottom: var(--space-2);
  font-weight: var(--font-weight-medium);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--light);
  border-radius: var(--radius-md);
  font-size: var(--font-size-md);
}

.form-group textarea {
  min-height: 100px;
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-4);
}

@media (max-width: 768px) {
  .summary-card {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .action-buttons {
    flex-direction: column;
  }

  .summary-stats {
    margin-top: var(--space-4);
  }
}
</style>
