import { mount, flushPromises } from "@vue/test-utils";
import LoginForm from "@/components/LoginForm.vue";
import { createStore } from "vuex";

describe("LoginForm.vue", () => {
  let actions, store, routerPushMock;

  beforeEach(() => {
    actions = {
      login: jest.fn(() => Promise.resolve()),
    };
    store = createStore({
      state: {},
      actions,
    });
    routerPushMock = jest.fn();
  });

  it("renders the login form correctly", () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [store],
        mocks: {
          $router: { push: routerPushMock },
        },
      },
    });

    expect(wrapper.find('input[type="email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    expect(wrapper.find("button").text()).toBe("Login");
  });

  it("calls the login action on form submission and redirects on success", async () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
      },
    });

    // Set input values
    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("password123");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    expect(actions.login).toHaveBeenCalled();
    expect(routerPushMock).toHaveBeenCalledWith("/");
  });

  it("displays an error message when login fails", async () => {
    actions.login.mockRejectedValueOnce({
      response: { data: { error: "Invalid credentials" } },
    });
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
      },
    });

    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("wrongpassword");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    expect(wrapper.find(".error").text()).toBe("Invalid credentials");
  });
});
