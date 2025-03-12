import { mount } from "@vue/test-utils";
import NotificationComponent from "@/components/NotificationComponent.vue";
import eventBus from "@/utils/eventBus";

// Mock the router
const mockRouter = { push: jest.fn() };
const mockRoute = { name: "Dashboard", query: {} };

// Mock Vue Router's useRouter function
jest.mock("vue-router", () => ({
  useRouter: () => mockRouter,
}));

describe("NotificationComponent", () => {
  let wrapper;

  beforeEach(() => {
    // Reset mock between tests
    jest.clearAllMocks();

    // Create the component wrapper
    wrapper = mount(NotificationComponent, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: mockRoute,
        },
      },
    });
  });

  it("should render when a message is received", async () => {
    // Initially no active message
    expect(wrapper.find(".message-preview-container").exists()).toBe(false);

    // Emit a chat message event
    const message = {
      conversationId: 1,
      senderName: "Test User",
      text: "Hello, this is a test message",
      timestamp: new Date(),
    };

    // Simulate receiving a message
    eventBus.emit("chat:incoming-message", message);

    // Wait for DOM updates
    await wrapper.vm.$nextTick();

    // Component should now be visible
    expect(wrapper.find(".message-preview-container").exists()).toBe(true);
    expect(wrapper.find(".message-sender").text()).toBe("Test User");
    expect(wrapper.find(".message-text").text()).toBe(
      "Hello, this is a test message",
    );
  });

  it("should render when a notification is shown", async () => {
    // Initially no active message
    expect(wrapper.find(".message-preview-container").exists()).toBe(false);

    // Emit a general notification event
    const notification = {
      type: "success",
      title: "Success",
      message: "Operation completed successfully",
      duration: 3000,
    };

    // Simulate showing a notification
    eventBus.emit("show-notification", notification);

    // Wait for DOM updates
    await wrapper.vm.$nextTick();

    // Component should now be visible with correct notification
    expect(wrapper.find(".message-preview-container").exists()).toBe(true);
    expect(wrapper.find(".message-sender").text()).toBe("Success");
    expect(wrapper.find(".message-text").text()).toBe(
      "Operation completed successfully",
    );
    expect(wrapper.find(".message-preview").classes()).toContain(
      "notification-success",
    );
  });

  it("should close when the close button is clicked", async () => {
    // Show a notification first
    eventBus.emit("show-notification", {
      title: "Test",
      message: "Test message",
    });

    await wrapper.vm.$nextTick();

    // Component should be visible
    expect(wrapper.find(".message-preview-container").exists()).toBe(true);

    // Click the close button
    await wrapper.find(".message-close").trigger("click");

    // Component should be hidden
    await wrapper.vm.$nextTick();
    expect(wrapper.find(".message-preview-container").exists()).toBe(false);
  });

  it("should navigate to conversation when a chat message is clicked", async () => {
    // Show a chat message
    const message = {
      conversationId: 123,
      senderName: "Test User",
      text: "Hello, this is a test message",
      timestamp: new Date(),
    };

    eventBus.emit("chat:incoming-message", message);
    await wrapper.vm.$nextTick();

    // Click on the message
    await wrapper.find(".message-preview").trigger("click");

    // Router should have been called with correct route
    expect(mockRouter.push).toHaveBeenCalledWith({
      name: "Chat",
      query: { conversation: 123 },
    });
  });

  it("should clear messages when chat:read-messages is emitted", async () => {
    // Show a message first
    eventBus.emit("chat:incoming-message", {
      conversationId: 1,
      senderName: "Test User",
      text: "Test message",
      timestamp: new Date(),
    });

    await wrapper.vm.$nextTick();
    expect(wrapper.find(".message-preview-container").exists()).toBe(true);

    // Emit read messages event
    eventBus.emit("chat:read-messages");

    // Component should be hidden
    await wrapper.vm.$nextTick();
    expect(wrapper.find(".message-preview-container").exists()).toBe(false);
  });
});
