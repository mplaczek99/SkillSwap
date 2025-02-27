import { mount } from "@vue/test-utils";
import Dashboard from "@/components/Dashboard.vue";
import { createStore } from "vuex";

describe("Dashboard.vue", () => {
  let store;
  
  beforeEach(() => {
    // Create a simple Vuex store with a user state for testing.
    store = createStore({
      state: {
        user: { name: "Test User", skillPoints: 20 },
      },
    });
  });

  it("renders user greeting and SkillPoints", () => {
    const wrapper = mount(Dashboard, {
      global: {
        plugins: [store],
        stubs: {
          'router-link': true,
          'font-awesome-icon': true
        }
      },
    });
    expect(wrapper.text()).toContain("Welcome to SkillSwap, Test User!");
    // Updated assertion for the new UI format
    expect(wrapper.text()).toContain("Your SkillPoints");
    expect(wrapper.text()).toContain("20");
  });

  it("renders featured skills, recent activities, and announcements", () => {
    const wrapper = mount(Dashboard, {
      global: {
        plugins: [store],
        stubs: {
          'router-link': true,
          'font-awesome-icon': true
        }
      },
    });
    expect(wrapper.text()).toContain("Featured Skills");
    expect(wrapper.text()).toContain("Recent Activity");
    expect(wrapper.text()).toContain("Announcements");
  });
});
