import { mount, flushPromises } from "@vue/test-utils";
import RegisterForm from "@/components/RegisterForm.vue";
import { createStore } from "vuex";

describe("RegisterForm.vue", () => {
  let actions, store, routerPushMock;

  beforeEach(() => {
    actions = {
      register: jest.fn(() => Promise.resolve()),
    };
    store = createStore({
      state: {},
      actions,
    });
    routerPushMock = jest.fn();
  });

  it("renders the register form correctly", () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
      },
    });

    expect(wrapper.find('input[type="text"]').exists()).toBe(true);
    expect(wrapper.find('input[type="email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    expect(wrapper.find("button").text()).toBe("Register");
  });

  it("calls the register action on form submission and redirects on success", async () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
      },
    });

    await wrapper.find('input[type="text"]').setValue("Test User");
    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("password123");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    expect(actions.register).toHaveBeenCalled();
    expect(routerPushMock).toHaveBeenCalledWith("/");
  });

  it("displays an error message when registration fails", async () => {
    actions.register.mockRejectedValueOnce({
      response: { data: { error: "Registration error" } },
    });
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
      },
    });

    await wrapper.find('input[type="text"]').setValue("Test User");
    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("weakpassword");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    expect(wrapper.find(".error").text()).toBe("Registration error");
  });
});
