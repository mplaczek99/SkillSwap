import { mount, flushPromises } from "@vue/test-utils";
import Transactions from "@/components/Transactions.vue";
import { createStore } from "vuex";

describe("Transactions.vue", () => {
  let store;

  beforeEach(() => {
    // Create a test store with a user state
    store = createStore({
      state: {
        user: {
          id: 1,
          name: "Test User",
          email: "test@example.com",
          skillPoints: 50,
        },
      },
      getters: {
        user: (state) => state.user,
      },
    });

    // Mock setTimeout to make tests run faster
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("renders user balance and transaction stats", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Advance timers to trigger the mocked API response
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Should display user balance
    expect(wrapper.find(".balance").text()).toContain("50");

    // Should display statistics
    expect(wrapper.find(".stat-item").exists()).toBe(true);
    expect(wrapper.text()).toContain("Earned");
    expect(wrapper.text()).toContain("Spent");
  });

  it("renders transactions list", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Advance timers to trigger the mocked API response
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Should display transaction items
    const transactions = wrapper.findAll(".transaction-item");
    expect(transactions.length).toBeGreaterThan(0);
  });

  it("shows send points modal", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Initially modal should be hidden
    expect(wrapper.find(".modal-backdrop").exists()).toBe(false);

    // Click send button
    await wrapper.find("button.btn-primary").trigger("click");

    // Modal should be visible
    expect(wrapper.find(".modal-backdrop").exists()).toBe(true);
    expect(wrapper.find(".modal-header").text()).toContain("Send SkillPoints");
  });

  it("validates send points form", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Open modal
    await wrapper.find("button.btn-primary").trigger("click");

    // Submit button should be disabled with invalid amount
    await wrapper.find("#amount").setValue(-10);
    expect(
      wrapper.find(".form-actions .btn-primary").attributes("disabled"),
    ).toBeDefined();

    // Submit button should be disabled with amount greater than balance
    await wrapper.find("#amount").setValue(100); // Balance is 50
    expect(
      wrapper.find(".form-actions .btn-primary").attributes("disabled"),
    ).toBeDefined();

    // Submit button should be enabled with valid amount
    await wrapper.find("#amount").setValue(10);
    expect(
      wrapper.find(".form-actions .btn-primary").attributes("disabled"),
    ).toBeUndefined();
  });

  it("submits send points form", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Open modal
    await wrapper.find("button.btn-primary").trigger("click");

    // Fill form
    await wrapper.find("#recipient").setValue("alice@example.com");
    await wrapper.find("#amount").setValue(10);
    await wrapper.find("#note").setValue("For JavaScript lessons");

    // Mock alert
    const originalAlert = window.alert;
    window.alert = jest.fn();

    // Submit form
    await wrapper.find(".modal-body form").trigger("submit.prevent");

    // Alert should be called
    expect(window.alert).toHaveBeenCalled();

    // Modal should be closed
    expect(wrapper.find(".modal-backdrop").exists()).toBe(false);

    // Restore original alert
    window.alert = originalAlert;
  });

  it("refreshes transactions", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Create a replacement for fetchTransactions
    const originalFetchTransactions = wrapper.vm.fetchTransactions;
    const mockFetchTransactions = jest.fn();
    wrapper.vm.fetchTransactions = mockFetchTransactions;

    // Click refresh button
    await wrapper.find(".action-buttons .btn-outline").trigger("click");

    // Method should be called
    expect(mockFetchTransactions).toHaveBeenCalled();

    // Restore original method after test
    wrapper.vm.fetchTransactions = originalFetchTransactions;
  });

  it("formats transaction items correctly", async () => {
    const wrapper = mount(Transactions, {
      global: {
        plugins: [store],
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Advance timers to trigger the mocked API response
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Test transaction formatting helpers
    expect(
      wrapper.vm.getTransactionClass({
        receiverId: 1, // Same as user ID
        senderId: 2,
      }),
    ).toEqual({ received: true, sent: false });

    expect(
      wrapper.vm.getTransactionClass({
        receiverId: 2,
        senderId: 1, // Same as user ID
      }),
    ).toEqual({ received: false, sent: true });

    expect(
      wrapper.vm.getAmountPrefix({
        receiverId: 1, // Same as user ID
      }),
    ).toBe("+");

    expect(
      wrapper.vm.getAmountPrefix({
        receiverId: 2,
        senderId: 1, // Same as user ID
      }),
    ).toBe("-");
  });
});
