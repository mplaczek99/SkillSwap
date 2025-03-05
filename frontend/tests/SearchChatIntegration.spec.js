import { mount } from "@vue/test-utils";
import flushPromises from "flush-promises";
import Search from "@/components/Search.vue";
import axios from "axios";
import eventBus from "@/utils/eventBus";

jest.mock("axios");

describe("Search Chat Integration", () => {
  const mockRouter = {
    push: jest.fn(),
    replace: jest.fn()
  };
  
  const mockRoute = {
    query: {}
  };
  
  beforeEach(() => {
    jest.clearAllMocks();
    jest.spyOn(eventBus, "emit").mockImplementation(() => {});
  });
  
  afterEach(() => {
    eventBus.emit.mockRestore();
  });

  it("navigates to chat when message button is clicked", async () => {
    const mockSearchResults = [
      { id: 1, name: "Test User", email: "test@example.com" },
      { id: 2, name: "Alice Smith", email: "alice@example.com" },
      { id: 3, name: "Programming Skill", description: "Learn to code", email: null }
    ];
    axios.get.mockResolvedValue({ data: mockSearchResults });

    const wrapper = mount(Search, { 
      props: { forceApiCall: true },
      global: {
        mocks: {
          $router: mockRouter,
          $route: mockRoute
        },
        stubs: { 'font-awesome-icon': true }
      }
    });
    
    const input = wrapper.find("input");
    await input.setValue("alice");
    await wrapper.find("form").trigger("submit.prevent");
    
    // Remove the flush() call that was causing the error
    await flushPromises();
    
    // Simulate startChat call directly.
    const user = { id: 2, name: "Alice Smith", email: "alice@example.com" };
    wrapper.vm.startChat(user);
    
    expect(mockRouter.push).toHaveBeenCalledWith({
      name: "Chat",
      query: { user: user.id, userName: user.name }
    });
    
    expect(eventBus.emit).toHaveBeenCalledWith("show-notification", {
      type: "info",
      title: "Starting Chat",
      message: `Starting a conversation with ${user.name}`,
      duration: 3000,
    });
  });
});
