import { mount, flushPromises } from "@vue/test-utils";
import Chat from "@/components/Chat.vue";
import ChatService from "@/services/ChatService";

// Mock the ChatService
jest.mock("@/services/ChatService", () => ({
  getConversations: jest.fn(),
  getConversation: jest.fn(),
  sendMessage: jest.fn(),
  startConversation: jest.fn(),
  getUnreadCount: jest.fn(),
  simulateIncomingMessage: jest.fn()
}));

describe("Chat.vue", () => {
  let mockRoute, mockRouter;

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Mock conversations data
    const mockConversations = [
      {
        id: 1,
        recipient: { id: 2, name: "Alice Smith", avatar: null },
        lastMessage: { text: "Hello!", timestamp: new Date() },
        unreadCount: 0
      },
      {
        id: 2,
        recipient: { id: 3, name: "Bob Johnson", avatar: null },
        lastMessage: { text: "How are you?", timestamp: new Date() },
        unreadCount: 1
      }
    ];
    
    // Mock conversation data
    const mockConversation = {
      id: 1,
      recipient: { id: 2, name: "Alice Smith", avatar: null },
      messages: [
        { id: 1, senderId: 2, text: "Hello!", timestamp: new Date(), isOutgoing: false },
        { id: 2, senderId: 1, text: "Hi!", timestamp: new Date(), isOutgoing: true }
      ]
    };
    
    // Mock ChatService methods
    ChatService.getConversations.mockResolvedValue(mockConversations);
    ChatService.getConversation.mockResolvedValue(mockConversation);
    ChatService.sendMessage.mockImplementation((id, text) => {
      return Promise.resolve({
        id: 3,
        senderId: 1,
        text,
        timestamp: new Date(),
        isOutgoing: true
      });
    });
    ChatService.startConversation.mockResolvedValue(1);
    ChatService.getUnreadCount.mockResolvedValue(1);
    
    // Mock route and router
    mockRoute = {
      query: {}
    };
    
    mockRouter = {
      push: jest.fn(),
      replace: jest.fn()
    };
  });

  it("loads conversations on creation", async () => {
    const wrapper = mount(Chat, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true
        }
      }
    });
    
    await flushPromises();
    
    expect(ChatService.getConversations).toHaveBeenCalled();
    expect(wrapper.findAll(".conversation-item").length).toBe(2);
    expect(wrapper.find(".conversation-name").text()).toContain("Alice Smith");
  });

  it("loads a conversation when clicked", async () => {
    const wrapper = mount(Chat, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true
        }
      }
    });
    
    await flushPromises();
    
    // Click the first conversation
    await wrapper.findAll(".conversation-item")[0].trigger("click");
    
    await flushPromises();
    
    expect(ChatService.getConversation).toHaveBeenCalledWith(1);
    expect(wrapper.find(".chat-header").exists()).toBe(true);
    expect(wrapper.find(".user-name").text()).toContain("Alice Smith");
    expect(wrapper.findAll(".message").length).toBe(2);
  });

  it("sends a message", async () => {
    const wrapper = mount(Chat, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true
        }
      }
    });
    
    await flushPromises();
    
    // Click the first conversation
    await wrapper.findAll(".conversation-item")[0].trigger("click");
    
    await flushPromises();
    
    // Set message text
    await wrapper.find("textarea").setValue("Test message");
    
    // Submit the form
    await wrapper.find(".message-form").trigger("submit.prevent");
    
    await flushPromises();
    
    expect(ChatService.sendMessage).toHaveBeenCalledWith(1, "Test message");
    expect(wrapper.find("textarea").element.value).toBe("");
  });

  it("loads a conversation from route query", async () => {
    // Set conversation ID in route query as a number (instead of a string)
    mockRoute.query = { conversation: 2 };
    
    mount(Chat, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true
        }
      }
    });
    
    await flushPromises();
    
    expect(ChatService.getConversation).toHaveBeenCalledWith(2);
  });

  it("starts a new conversation from route query", async () => {
    // Set user ID in route query
    mockRoute.query = { user: 4, userName: "David Brown" };
    
    mount(Chat, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true
        }
      }
    });
    
    await flushPromises();
    
    expect(ChatService.startConversation).toHaveBeenCalledWith(4, "David Brown");
    expect(ChatService.getConversations).toHaveBeenCalled();
    expect(ChatService.getConversation).toHaveBeenCalled();
  });

  it("searches for users", async () => {
    const wrapper = mount(Chat, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true
        }
      }
    });
    
    await flushPromises();
    
    // Set search query
    await wrapper.find(".chat-search input").setValue("ali");
    
    // Call the search method directly since debounce makes it hard to test
    wrapper.vm.performSearch();
    
    await wrapper.vm.$nextTick();
    
    expect(wrapper.find(".search-results").exists()).toBe(true);
    expect(wrapper.find(".search-result-item").exists()).toBe(true);
    // (Assuming the test data includes "Alice")
  });
});

