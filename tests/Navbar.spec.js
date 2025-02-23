import { mount } from "@vue/test-utils";
import Navbar from "@/components/Navbar.vue";
import { createStore } from "vuex";

describe("Navbar.vue", () => {
  let store, mutations;

  beforeEach(() => {
    // Create a dummy logout mutation spy
    mutations = { logout: jest.fn() };
  });

  it("renders login and register links when unauthenticated", () => {
    store = createStore({
      state: { token: null },
      getters: { isAuthenticated: (state) => !!state.token },
      mutations,
    });
    const wrapper = mount(Navbar, {
      global: {
        plugins: [store],
        // Stub router-link to render its default slot as a simple <a> tag
        stubs: {
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
        mocks: {
          // Provide a dummy $router so that calls to push do not fail
          $router: { push: jest.fn() },
        },
      },
    });

    // Now the stub renders text content from its slot
    expect(wrapper.text()).toContain("Login");
    expect(wrapper.text()).toContain("Register");
    expect(wrapper.text()).not.toContain("Logout");
  });

  it("renders logout button when authenticated and calls logout mutation", async () => {
    store = createStore({
      state: { token: "dummy-token" },
      getters: { isAuthenticated: (state) => !!state.token },
      mutations,
    });
    const routerPushMock = jest.fn();
    const wrapper = mount(Navbar, {
      global: {
        plugins: [store],
        stubs: {
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
        mocks: {
          $router: { push: routerPushMock },
        },
      },
    });

    expect(wrapper.text()).toContain("Logout");
    const logoutButton = wrapper.find("button");
    await logoutButton.trigger("click");
    expect(mutations.logout).toHaveBeenCalled();
    expect(routerPushMock).toHaveBeenCalledWith("/login");
  });
});
