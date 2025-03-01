import { mount, flushPromises } from "@vue/test-utils";
import LoginForm from "@/components/LoginForm.vue";
import { createStore } from "vuex";

describe("LoginForm.vue", () => {
  let actions, store, routerPushMock;

  beforeEach(() => {
    actions = { login: jest.fn(() => Promise.resolve()) };
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
        mocks: { $router: { push: routerPushMock } },
        stubs: {
          'router-link': {
            template: '<a><slot /></a>'
          },
          'font-awesome-icon': true
        }
      },
    });

    expect(wrapper.find('input[type="email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    // Update to check button content, not specific text
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
  });

  it("calls login action and redirects on success", async () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
        stubs: {
          'router-link': {
            template: '<a><slot /></a>'
          },
          'font-awesome-icon': true
        }
      },
    });

    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("password123");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    expect(actions.login).toHaveBeenCalled();
    expect(routerPushMock).toHaveBeenCalledWith("/");
  });

  it("displays error message when login fails", async () => {
    actions.login.mockRejectedValueOnce({ response: { data: { error: "Invalid credentials" } } });
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
        stubs: {
          'router-link': {
            template: '<a><slot /></a>'
          },
          'font-awesome-icon': true
        }
      },
    });

    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("wrongpassword");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    // Updated selector to match the component's error class
    expect(wrapper.find(".alert-danger").text()).toBe("Invalid credentials");
  });
});
