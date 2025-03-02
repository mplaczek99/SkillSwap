import { mount } from "@vue/test-utils";
import Navbar from "@/components/Navbar.vue";
import { createStore } from "vuex";

describe("Navbar.vue", () => {
  let store, mutations;
  const routerPushMock = jest.fn();
  const mockRoot = { $on: jest.fn(), $off: jest.fn() };

  beforeEach(() => {
    mutations = { logout: jest.fn() };
    store = createStore({
      state: { token: null },
      getters: { isAuthenticated: (state) => !!state.token },
      mutations,
    });
  });

  it("renders login and register links when unauthenticated", () => {
    const wrapper = mount(Navbar, {
      global: {
        plugins: [store],
        mocks: { 
          $router: { push: routerPushMock }
        },
        config: {
          globalProperties: { $root: mockRoot }
        },
        stubs: {
          "router-link": { template: "<a><slot /></a>" }
        },
      },
    });

    expect(wrapper.text()).toContain("Login");
    expect(wrapper.text()).toContain("Register");
    expect(wrapper.text()).not.toContain("Logout");
  });

  it("renders logout button when authenticated and calls logout", async () => {
    store = createStore({
      state: { token: "dummy-token" },
      getters: { isAuthenticated: (state) => !!state.token },
      mutations,
    });
    const wrapper = mount(Navbar, {
      global: {
        plugins: [store],
        mocks: { 
          $router: { push: routerPushMock }
        },
        config: {
          globalProperties: { $root: mockRoot }
        },
        stubs: { "router-link": { template: "<a><slot /></a>" } },
      },
    });

    expect(wrapper.text()).toContain("Logout");
    const logoutButton = wrapper.find("button.btn-outline");
    await logoutButton.trigger("click");
    expect(mutations.logout).toHaveBeenCalled();
    expect(routerPushMock).toHaveBeenCalledWith("/login");
  });
});

