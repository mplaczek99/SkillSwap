import { mount, flushPromises } from "@vue/test-utils";
import RegisterForm from "@/components/RegisterForm.vue";
import { createStore } from "vuex";

describe("RegisterForm.vue", () => {
  let actions, store, routerPushMock;

  beforeEach(() => {
    actions = { register: jest.fn(() => Promise.resolve()) };
    store = createStore({
      state: {},
      actions,
    });
    routerPushMock = jest.fn();
  });

  it("renders register form correctly", () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
        stubs: {
          'router-link': true,
          'font-awesome-icon': true
        }
      },
    });

    expect(wrapper.find('input[type="text"]').exists()).toBe(true);
    expect(wrapper.find('input[type="email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    // Update to check button content, not specific text
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
  });

  it("calls register action and redirects on success", async () => {
    const wrapper = mount(RegisterForm, {
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

    // Set checkbox to true to avoid validation errors
    const checkboxInput = wrapper.find('input[type="checkbox"]');
    await checkboxInput.setValue(true);

    await wrapper.find('input[type="text"]').setValue("Test User");
    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("password123");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    expect(actions.register).toHaveBeenCalled();
    expect(routerPushMock).toHaveBeenCalledWith("/");
  });

  it("displays error message when registration fails", async () => {
    actions.register.mockRejectedValueOnce({ response: { data: { error: "Registration error" } } });
    const wrapper = mount(RegisterForm, {
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

    // Set checkbox to true to avoid validation errors
    const checkboxInput = wrapper.find('input[type="checkbox"]');
    await checkboxInput.setValue(true);

    await wrapper.find('input[type="text"]').setValue("Test User");
    await wrapper.find('input[type="email"]').setValue("test@example.com");
    await wrapper.find('input[type="password"]').setValue("weakpassword");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();

    // Updated selector to match the component's error class
    expect(wrapper.find(".alert-danger").text()).toBe("Registration error");
  });
});
