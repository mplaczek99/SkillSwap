import { mount } from "@vue/test-utils";
import { createStore } from "vuex";
import Profile from "@/components/Profile.vue";

describe("Profile.vue", () => {
  let store, actions, state, getters, routerPushMock;
  const userData = {
    name: "Test User",
    email: "test@example.com",
    bio: "Hello world",
    avatar: "",
    skillPoints: 10,
  };

  beforeEach(() => {
    actions = {
      updateProfile: jest.fn(() => Promise.resolve()),
    };
    state = { user: userData };
    getters = { user: (state) => state.user };
    store = createStore({ state, getters, actions });
    routerPushMock = jest.fn();
  });

  it("renders profile overview with user data", () => {
    const wrapper = mount(Profile, {
      global: {
        plugins: [store],
        mocks: { 
          $router: { push: routerPushMock },
          $route: { params: {} } // Provide an empty params object
        },
        stubs: { ProfileCard: true },
      },
    });

    expect(wrapper.text()).toContain(userData.name);
    expect(wrapper.text()).toContain(userData.email);
    expect(wrapper.text()).toContain(userData.bio);
    expect(wrapper.text()).toContain("SkillPoints: 10");
  });

  it("toggles edit mode when clicking edit button", async () => {
    const wrapper = mount(Profile, {
      global: {
        plugins: [store],
        mocks: { 
          $router: { push: routerPushMock },
          $route: { params: {} }
        },
        stubs: { ProfileCard: true },
      },
    });

    expect(wrapper.find('[data-test="edit-profile-form"]').exists()).toBe(false);

    await wrapper.find('[data-test="edit-button"]').trigger("click");
    expect(wrapper.find('[data-test="edit-profile-form"]').exists()).toBe(true);

    const nameInput = wrapper.find("input#name");
    expect(nameInput.element.value).toBe(userData.name);

    await wrapper.find('[data-test="edit-button"]').trigger("click");
    expect(wrapper.find('[data-test="edit-profile-form"]').exists()).toBe(false);
  });

  it("dispatches updateProfile on form submission", async () => {
    const wrapper = mount(Profile, {
      global: {
        plugins: [store],
        mocks: { 
          $router: { push: routerPushMock },
          $route: { params: {} }
        },
        stubs: { ProfileCard: true },
      },
    });

    await wrapper.find('[data-test="edit-button"]').trigger("click");

    const bioTextarea = wrapper.find("textarea#bio");
    await bioTextarea.setValue("Updated bio");

    await wrapper.find('[data-test="edit-profile-form"]').trigger("submit.prevent");

    expect(actions.updateProfile).toHaveBeenCalled();
    expect(actions.updateProfile.mock.calls[0][1]).toEqual(
      expect.objectContaining({ bio: "Updated bio" })
    );
  });

  it("renders My Skills section with dummy skills", () => {
    const wrapper = mount(Profile, {
      global: {
        plugins: [store],
        mocks: { 
          $router: { push: routerPushMock },
          $route: { params: {} }
        },
        stubs: {
          ProfileCard: {
            template: '<div class="profile-card-stub">{{ title }} - {{ description }}</div>',
            props: ["title", "description"],
          },
        },
      },
    });

    const skillsSection = wrapper.find(".my-skills");
    expect(skillsSection.exists()).toBe(true);
    expect(skillsSection.text()).toContain("Go Programming");
    expect(skillsSection.text()).toContain("Vue.js Development");
  });
});

