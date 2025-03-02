import { mount } from "@vue/test-utils";
import ChatNotificationPreview from "@/components/ChatNotificationPreview.vue";

const mockRouter = { push: jest.fn() };
const mockRoute = { name: 'Dashboard', query: {} };
const mockRoot = { $on: jest.fn(), $off: jest.fn(), $emit: jest.fn() };

describe("ChatNotificationPreview.vue", () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });
  
  afterEach(() => {
    jest.useRealTimers();
  });

  it("displays message preview when a message event is emitted", async () => {
    const wrapper = mount(ChatNotificationPreview, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: mockRoute
        },
        config: {
          globalProperties: { $root: mockRoot }
        },
        stubs: { 'font-awesome-icon': true }
      }
    });
    
    // Simulate incoming message by calling the method directly
    wrapper.vm.handleIncomingMessage({
      conversationId: 1,
      senderName: 'Alice Smith',
      senderAvatar: null,
      text: 'Hello there!',
      timestamp: new Date()
    });
    await wrapper.vm.$nextTick();
    
    expect(wrapper.find('.message-preview').exists()).toBe(true);
    expect(wrapper.find('.message-sender').text()).toBe('Alice Smith');
    expect(wrapper.find('.message-text').text()).toBe('Hello there!');
  });

  it("closes message preview after timeout", async () => {
    const wrapper = mount(ChatNotificationPreview, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: mockRoute
        },
        config: {
          globalProperties: { $root: mockRoot }
        },
        stubs: { 'font-awesome-icon': true }
      }
    });
    
    wrapper.vm.handleIncomingMessage({
      conversationId: 1,
      senderName: 'Alice Smith',
      senderAvatar: null,
      text: 'Hello there!',
      timestamp: new Date()
    });
    await wrapper.vm.$nextTick();
    
    expect(wrapper.find('.message-preview').exists()).toBe(true);
    
    jest.advanceTimersByTime(5000);
    await wrapper.vm.$nextTick();
    
    expect(wrapper.find('.message-preview').exists()).toBe(false);
  });

  it("navigates to conversation when clicked", async () => {
    const wrapper = mount(ChatNotificationPreview, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: mockRoute
        },
        config: {
          globalProperties: { $root: mockRoot }
        },
        stubs: { 'font-awesome-icon': true }
      }
    });
    
    wrapper.vm.handleIncomingMessage({
      conversationId: 1,
      senderName: 'Alice Smith',
      senderAvatar: null,
      text: 'Hello there!',
      timestamp: new Date()
    });
    await wrapper.vm.$nextTick();
    
    await wrapper.find('.message-preview').trigger('click');
    
    expect(mockRouter.push).toHaveBeenCalledWith({
      name: 'Chat',
      query: { conversation: 1 }
    });
    expect(wrapper.find('.message-preview').exists()).toBe(false);
  });

  it("closes preview when close button is clicked", async () => {
    const wrapper = mount(ChatNotificationPreview, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: mockRoute
        },
        config: {
          globalProperties: { $root: mockRoot }
        },
        stubs: { 'font-awesome-icon': true }
      }
    });
    
    wrapper.vm.handleIncomingMessage({
      conversationId: 1,
      senderName: 'Alice Smith',
      senderAvatar: null,
      text: 'Hello there!',
      timestamp: new Date()
    });
    await wrapper.vm.$nextTick();
    
    await wrapper.find('.message-close').trigger('click');
    
    expect(wrapper.find('.message-preview').exists()).toBe(false);
  });
});

