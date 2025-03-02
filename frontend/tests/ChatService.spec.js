// tests/ChatService.spec.js
import store from "@/store";
import ChatService from "@/services/ChatService";

jest.mock("jwt-decode", () =>
  jest.fn(() => ({
    user_id: 1,
    email: "test@example.com",
    role: "User",
  }))
);

describe("ChatService", () => {
  beforeEach(() => {
    // Set the store state for testing
    store.replaceState({ user: { id: 1, name: "Test User" } });
    // Speed up tests by mocking network delay.
    ChatService.simulateNetworkDelay = jest.fn(() => Promise.resolve());
  });

  describe("getConversations", () => {
    it("returns conversations for the current user", async () => {
      const conversations = await ChatService.getConversations();
      expect(Array.isArray(conversations)).toBe(true);
      expect(conversations.length).toBeGreaterThan(0);
      conversations.forEach((convo) => {
        expect(convo).toHaveProperty("id");
        expect(convo).toHaveProperty("recipient");
        expect(convo).toHaveProperty("lastMessage");
        expect(convo).toHaveProperty("unreadCount");
        expect(convo.recipient.id).not.toBe(store.state.user.id);
      });
    });
  });

  describe("getConversation", () => {
    it("returns a single conversation with formatted messages", async () => {
      const conversation = await ChatService.getConversation(1);
      expect(conversation).toHaveProperty("id", 1);
      expect(conversation).toHaveProperty("recipient");
      expect(conversation).toHaveProperty("messages");
      conversation.messages.forEach((msg) => {
        expect(msg).toHaveProperty("isOutgoing");
        if (msg.senderId === store.state.user.id) {
          expect(msg.isOutgoing).toBe(true);
        } else {
          expect(msg.isOutgoing).toBe(false);
        }
      });
    });

    it("throws an error for non-existent conversation", async () => {
      await expect(ChatService.getConversation(999)).rejects.toThrow(
        "Conversation not found"
      );
    });
  });

  describe("sendMessage", () => {
    it("adds a message to an existing conversation", async () => {
      const newMessage = await ChatService.sendMessage(1, "Test message");
      expect(newMessage).toHaveProperty("text", "Test message");
      expect(newMessage).toHaveProperty("senderId", store.state.user.id);
      expect(newMessage).toHaveProperty("isOutgoing", true);
      expect(newMessage).toHaveProperty("timestamp");
    });

    it("throws an error for non-existent conversation", async () => {
      await expect(ChatService.sendMessage(999, "Test message")).rejects.toThrow(
        "Conversation not found"
      );
    });
  });

  describe("startConversation", () => {
    it("returns existing conversation ID if conversation already exists", async () => {
      const conversationId = await ChatService.startConversation(2, "Alice Smith");
      expect(typeof conversationId).toBe("number");
      const conversations = await ChatService.getConversations();
      const matchingConversation = conversations.find(
        (c) => c.recipient.id === 2
      );
      expect(conversationId).toBe(matchingConversation.id);
    });

    it("creates a new conversation if one does not exist", async () => {
      const userId = 999;
      const userName = "New User";
      const beforeConversations = await ChatService.getConversations();
      const beforeCount = beforeConversations.length;
      await ChatService.startConversation(userId, userName, "Hello!");
      const afterConversations = await ChatService.getConversations();
      expect(afterConversations.length).toBe(beforeCount + 1);
      const newConversation = afterConversations.find(
        (c) => c.recipient.id === userId
      );
      expect(newConversation).toBeDefined();
    });
  });

  describe("getUnreadCount", () => {
    it("returns the total number of unread messages", async () => {
      const unreadCount = await ChatService.getUnreadCount();
      expect(typeof unreadCount).toBe("number");
    });
  });

  describe("simulateIncomingMessage", () => {
    it("adds a message to the conversation and marks it as unread", async () => {
      const initialUnreadCount = await ChatService.getUnreadCount();
      const conversations = await ChatService.getConversations();
      const conversationId = conversations[0].id;
      const message = await ChatService.simulateIncomingMessage(
        conversationId,
        "Test incoming message"
      );
      expect(message).toHaveProperty("text", "Test incoming message");
      expect(message).toHaveProperty("isOutgoing", false);
      expect(message).toHaveProperty("conversationId", conversationId);
      const newUnreadCount = await ChatService.getUnreadCount();
      expect(newUnreadCount).toBe(initialUnreadCount + 1);
    });
  });
});

