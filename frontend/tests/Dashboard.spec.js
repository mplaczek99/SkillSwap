import { mount, shallowMount } from "@vue/test-utils";
import Dashboard from "@/components/Dashboard.vue";
import { createStore } from "vuex";

// Mock the SVG import in Dashboard.vue
jest.mock("@/assets/images/skill-sharing.svg", () => "mock-svg", {
  virtual: true,
});

// Mock Vue Router
const mockRouter = {
  push: jest.fn(),
};

describe("Dashboard.vue", () => {
  let wrapper;
  let store;

  beforeEach(() => {
    // Create a mock store with the necessary state
    store = createStore({
      state: {
        user: {
          name: "Test User",
          skillPoints: 100,
        },
      },
      getters: {
        user: (state) => state.user,
      },
    });

    // Mount the component with mocked dependencies
    wrapper = shallowMount(Dashboard, {
      global: {
        plugins: [store],
        mocks: {
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });
  });

  it("renders welcome message with user name", () => {
    expect(wrapper.text()).toContain("Test User");
  });

  it("displays user skill points", () => {
    expect(wrapper.text()).toContain("100");
  });

  it("has four quick action buttons", () => {
    expect(wrapper.findAll(".action-button").length).toBe(4);
  });
});
