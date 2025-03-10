import { mount } from "@vue/test-utils";
import Transactions from "@/components/Transactions.vue";
import { createStore } from "vuex";
import axios from "axios";

// Mock axios
jest.mock("axios");

describe("Transactions.vue", () => {
  let wrapper;
  let store;

  beforeEach(() => {
    // Mock axios response for /api/transactions
    axios.get.mockResolvedValue({
      data: [
        {
          id: 1,
          senderId: 2,
          receiverId: 1,
          amount: 15,
          createdAt: new Date(Date.now() - 86400000), // 1 day ago
          senderName: "Alice Smith",
          receiverName: "Test User",
          note: "For JavaScript tutoring",
        },
        {
          id: 2,
          senderId: 1,
          receiverId: 3,
          amount: 5,
          createdAt: new Date(Date.now() - 172800000), // 2 days ago
          senderName: "Test User",
          receiverName: "Bob Johnson",
          note: "For cooking lessons",
        },
        {
          id: 3,
          senderId: 4,
          receiverId: 1,
          amount: 10,
          createdAt: new Date(Date.now() - 259200000), // 3 days ago
          senderName: "Carol Williams",
          receiverName: "Test User",
          note: "For guitar lessons",
        },
      ],
    });

    // Create a mock store
    store = createStore({
      state: {
        user: {
          id: 1,
          name: "Test User",
          email: "test@example.com",
          skillPoints: 100,
        },
      },
      getters: {
        user: (state) => state.user,
      },
    });

    // Mount the component with the store
    wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });
  });

  it("renders transactions list", async () => {
    // Wait for async operations to complete
    await wrapper.vm.$nextTick();
    await wrapper.vm.$nextTick(); // Sometimes we need multiple ticks for the DOM to update

    // Should display transaction items
    const transactions = wrapper.findAll(".transaction-item");
    expect(transactions.length).toBeGreaterThan(0);
  });

  it("shows send points modal", async () => {
    // Test showing and hiding the modal
    expect(wrapper.find(".modal-backdrop").exists()).toBe(false);

    // Find and click the button to show the modal
    const sendButton = wrapper.find(".btn-primary");
    await sendButton.trigger("click");

    // Modal should now be visible
    expect(wrapper.find(".modal-backdrop").exists()).toBe(true);

    // Find and click the cancel button
    const cancelButton = wrapper.find(".modal-content .btn-outline");
    await cancelButton.trigger("click");

    // Modal should be hidden again
    expect(wrapper.find(".modal-backdrop").exists()).toBe(false);
  });

  it("calculates totals correctly", async () => {
    // Wait for async operations to complete
    await wrapper.vm.$nextTick();

    // Call calculateTotals directly to test its logic
    wrapper.vm.calculateTotals();

    // Verify the expected totals
    // Received: transactions with ID 1 and 3 (total: 25)
    // Sent: transaction with ID 2 (total: 5)
    expect(wrapper.vm.totalEarned).toBe(25);
    expect(wrapper.vm.totalSpent).toBe(5);
  });
});
