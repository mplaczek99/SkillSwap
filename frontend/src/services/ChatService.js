import store from "@/store";

// Simulated chat database (in a real app, this would come from the server)
let conversations = [
  {
    id: 1,
    participants: [
      { id: 1, name: "Test User", avatar: null },
      { id: 2, name: "Alice Smith", avatar: null },
    ],
    messages: [
      {
        id: 1,
        senderId: 2,
        text: "Hello! I saw you were interested in learning JavaScript.",
        timestamp: new Date(Date.now() - 86400000),
      },
      {
        id: 2,
        senderId: 1,
        text: "Yes, I'd love to learn more about Vue.js specifically.",
        timestamp: new Date(Date.now() - 85400000),
      },
      {
        id: 3,
        senderId: 2,
        text: "Great! I can help with that. When would you like to start?",
        timestamp: new Date(Date.now() - 84400000),
      },
    ],
    lastMessageTime: new Date(Date.now() - 84400000),
    unreadCount: 0,
  },
  {
    id: 2,
    participants: [
      { id: 1, name: "Test User", avatar: null },
      { id: 3, name: "Bob Johnson", avatar: null },
    ],
    messages: [
      {
        id: 1,
        senderId: 3,
        text: "Hi there! I noticed you're offering cooking lessons.",
        timestamp: new Date(Date.now() - 172800000),
      },
      {
        id: 2,
        senderId: 1,
        text: "Yes, I specialize in Italian cuisine. Would you be interested?",
        timestamp: new Date(Date.now() - 171800000),
      },
    ],
    lastMessageTime: new Date(Date.now() - 171800000),
    unreadCount: 1,
  },
];

class ChatService {
  /**
   * Get all conversations for the current user
   */
  async getConversations() {
    // Simulate API call
    await this.simulateNetworkDelay(300);

    // In a real app, we would fetch this from the server
    // For now, filter conversations that include the current user
    const currentUserId = store.state.user ? store.state.user.id : 1;

    return conversations
      .map((convo) => {
        // Find the other participant (not the current user)
        const otherParticipant = convo.participants.find(
          (p) => p.id !== currentUserId,
        ) || { id: 0, name: "Unknown User", avatar: null };

        const lastMessage =
          convo.messages.length > 0
            ? convo.messages[convo.messages.length - 1]
            : { text: "", timestamp: new Date() };

        return {
          id: convo.id,
          recipient: otherParticipant,
          lastMessage: lastMessage,
          unreadCount: convo.unreadCount || 0,
        };
      })
      .sort(
        (a, b) =>
          new Date(b.lastMessage.timestamp) - new Date(a.lastMessage.timestamp),
      );
  }

  /**
   * Get a specific conversation by ID with pagination support
   */
  async getConversation(conversationId, page = 1, messagesPerPage = 20) {
    // Simulate API call
    await this.simulateNetworkDelay(200);

    const currentUserId = store.state.user ? store.state.user.id : 1;
    const conversation = conversations.find((c) => c.id == conversationId);

    if (!conversation) {
      throw new Error("Conversation not found");
    }

    // Mark messages as read
    conversation.unreadCount = 0;

    // Format the conversation for the UI
    const otherParticipant = conversation.participants.find(
      (p) => p.id !== currentUserId,
    ) || { id: 0, name: "Unknown User", avatar: null };

    // Apply pagination to messages - more efficient approach
    const allMessages = [...conversation.messages];

    // Sort once in ascending order (oldest to newest)
    allMessages.sort((a, b) => {
      const dateA = new Date(a.timestamp);
      const dateB = new Date(b.timestamp);
      return dateA - dateB;
    });

    // Get the correct slice of messages based on page number
    const totalMessages = allMessages.length;
    // This maintains chronological order (oldest to newest)
    const startIndex = (page - 1) * messagesPerPage;
    const endIndex = Math.min(startIndex + messagesPerPage, totalMessages);
    const paginatedMessages = allMessages.slice(startIndex, endIndex);

    return {
      id: conversation.id,
      recipient: otherParticipant,
      messages: paginatedMessages.map((msg) => ({
        ...msg,
        isOutgoing: msg.senderId === currentUserId,
      })),
    };
  }

  /**
   * Get conversation messages with pagination
   */
  async getConversationMessages(
    conversationId,
    page = 1,
    messagesPerPage = 20,
  ) {
    // Simulate API call
    await this.simulateNetworkDelay(200);

    const currentUserId = store.state.user ? store.state.user.id : 1;
    const conversation = conversations.find((c) => c.id == conversationId);

    if (!conversation) {
      throw new Error("Conversation not found");
    }

    // All messages
    const allMessages = conversation.messages;

    // For pagination when retrieving older messages, we want:
    // 1. Latest messages to be at the end (for normal display)
    // 2. Calculate the correct page of older messages

    // Sort all messages in ascending order (oldest first)
    const sortedMessages = [...allMessages].sort(
      (a, b) => new Date(a.timestamp) - new Date(b.timestamp),
    );

    // Calculate indices for pagination
    const totalMessages = sortedMessages.length;
    const startIndex = Math.max(0, totalMessages - page * messagesPerPage);
    const endIndex = Math.max(0, totalMessages - (page - 1) * messagesPerPage);

    // Get the correct slice - when page=1, this gets the latest messages
    const paginatedMessages = sortedMessages.slice(startIndex, endIndex);

    // Return the messages with isOutgoing flag
    return paginatedMessages.map((msg) => ({
      ...msg,
      isOutgoing: msg.senderId === currentUserId,
    }));
  }

  /**
   * Send a message in a conversation
   */
  async sendMessage(conversationId, text) {
    // Simulate API call
    await this.simulateNetworkDelay(300);

    const currentUserId = store.state.user ? store.state.user.id : 1;
    const conversation = conversations.find((c) => c.id == conversationId);

    if (!conversation) {
      throw new Error("Conversation not found");
    }

    // Create new message
    const newMessage = {
      id: conversation.messages.length + 1,
      senderId: currentUserId,
      text,
      timestamp: new Date(),
    };

    // Add to conversation
    conversation.messages.push(newMessage);
    conversation.lastMessageTime = newMessage.timestamp;

    return {
      ...newMessage,
      isOutgoing: true,
    };
  }

  /**
   * Start a new conversation with a user
   */
  async startConversation(userId, userName, initialMessage) {
    // Simulate API call
    await this.simulateNetworkDelay(500);

    if (!userId) {
      throw new Error("User ID is required to start a conversation");
    }

    const currentUserId = store.state.user ? store.state.user.id : 1;
    const currentUserName = store.state.user
      ? store.state.user.name
      : "Test User";

    // Check if conversation already exists
    const existingConvo = conversations.find(
      (c) =>
        c.participants.some((p) => p.id === currentUserId) &&
        c.participants.some((p) => p.id === userId),
    );

    if (existingConvo) {
      // If conversation exists, send message to existing conversation
      if (initialMessage) {
        await this.sendMessage(existingConvo.id, initialMessage);
      }
      return existingConvo.id;
    }

    // Create new conversation
    const newConversation = {
      id: conversations.length + 1,
      participants: [
        { id: currentUserId, name: currentUserName, avatar: null },
        { id: userId, name: userName || "Unknown User", avatar: null },
      ],
      messages: [],
      lastMessageTime: new Date(),
      unreadCount: 0,
    };

    // Add initial message if provided
    if (initialMessage) {
      const newMessage = {
        id: 1,
        senderId: currentUserId,
        text: initialMessage,
        timestamp: new Date(),
      };
      newConversation.messages.push(newMessage);
      // Update the lastMessageTime with the new message timestamp
      newConversation.lastMessageTime = newMessage.timestamp;
    }

    // Add to conversations
    conversations.push(newConversation);

    return newConversation.id;
  }

  /**
   * Get total unread messages count across all conversations
   */
  async getUnreadCount() {
    const convos = await this.getConversations();
    return convos.reduce((total, convo) => total + (convo.unreadCount || 0), 0);
  }

  /**
   * Helper to simulate network delay
   */
  simulateNetworkDelay(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  /**
   * Simulate receiving a new message (for demo purposes)
   */
  async simulateIncomingMessage(conversationId, text) {
    const conversation = conversations.find((c) => c.id == conversationId);
    if (!conversation) return null;

    const currentUserId = store.state.user ? store.state.user.id : 1;
    const otherParticipant = conversation.participants.find(
      (p) => p.id !== currentUserId,
    );

    if (!otherParticipant) return null;

    const newMessage = {
      id: conversation.messages.length + 1,
      senderId: otherParticipant.id,
      text,
      timestamp: new Date(),
    };

    conversation.messages.push(newMessage);
    conversation.lastMessageTime = newMessage.timestamp;
    conversation.unreadCount = (conversation.unreadCount || 0) + 1;

    // Return the formatted message
    return {
      ...newMessage,
      isOutgoing: false,
      conversationId,
    };
  }
}

export default new ChatService();
