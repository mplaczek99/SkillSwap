import { mount } from "@vue/test-utils";
import NotificationComponent from "@/components/NotificationComponent.vue";
import eventBus from "@/utils/eventBus";

// Mock the event bus
jest.mock("@/utils/eventBus", () => ({
  on: jest.fn(),
  off: jest.fn(),
  emit: jest.fn(),
}));

describe("NotificationComponent.vue", () => {
  let wrapper;

  beforeEach(() => {
    // Clear mock calls before each test
    jest.clearAllMocks();

    // Mount the component
    wrapper = mount(NotificationComponent);
  });

  it("sets up and cleans up event listeners", () => {
    // Verify that the event bus registers event listeners
    expect(eventBus.on).toHaveBeenCalledWith(
      "chat:incoming-message",
      expect.any(Function),
    );
    expect(eventBus.on).toHaveBeenCalledWith(
      "chat:read-messages",
      expect.any(Function),
    );

    // Unmount the component to trigger beforeUnmount
    wrapper.unmount();

    // Verify that the event listeners are removed
    expect(eventBus.off).toHaveBeenCalledWith(
      "chat:incoming-message",
      expect.any(Function),
    );
    expect(eventBus.off).toHaveBeenCalledWith(
      "chat:read-messages",
      expect.any(Function),
    );
  });

  it("shows notification when a message is received", async () => {
    // Get the handleIncomingMessage function that was registered
    const handleIncomingMessage = eventBus.on.mock.calls.find(
      (call) => call[0] === "chat:incoming-message",
    )[1];

    // Create a mock message
    const mockMessage = {
      conversationId: 1,
      senderName: "Test User",
      text: "Hello, this is a test message",
      timestamp: new Date(),
    };

    // Verify initial state
    expect(wrapper.vm.activeMessage).toBe(null);
    expect(wrapper.vm.messageQueue).toHaveLength(0);

    // Call the handler with the mock message
    handleIncomingMessage(mockMessage);

    // Check if the message is added to the queue and shown
    expect(wrapper.vm.activeMessage).toEqual(mockMessage);
  });

  it("clears messages when read", async () => {
    // Get the handleReadMessages function that was registered
    const handleReadMessages = eventBus.on.mock.calls.find(
      (call) => call[0] === "chat:read-messages",
    )[1];

    // Set some initial state
    wrapper.vm.activeMessage = { id: 1, text: "Test message" };
    wrapper.vm.messageQueue = [{ id: 2, text: "Another message" }];

    // Call the handler
    handleReadMessages();

    // Check if the messages are cleared
    expect(wrapper.vm.activeMessage).toBe(null);
    expect(wrapper.vm.messageQueue).toHaveLength(0);
  });
});
