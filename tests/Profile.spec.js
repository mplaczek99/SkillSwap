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
    state = {
      user: userData,
    };
    getters = {
      user: (state) => state.user,
    };
    store = createStore({
      state,
      getters,
      actions,
    });
    routerPushMock = jest.fn();
  });

  it("renders profile overview with user data", () => {
    const wrapper = mount(Profile, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
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
        mocks: { $router: { push: routerPushMock } },
        stubs: { ProfileCard: true },
      },
    });

    // Initially, the edit profile form should not be visible.
    expect(wrapper.find('[data-test="edit-profile-form"]').exists()).toBe(false);

    // Click the edit button to enter edit mode.
    await wrapper.find('[data-test="edit-button"]').trigger("click");
    expect(wrapper.find('[data-test="edit-profile-form"]').exists()).toBe(true);

    // Verify that the form inputs have the current user info.
    const nameInput = wrapper.find("input#name");
    expect(nameInput.element.value).toBe(userData.name);

    // Toggle edit mode off.
    await wrapper.find('[data-test="edit-button"]').trigger("click");
    expect(wrapper.find('[data-test="edit-profile-form"]').exists()).toBe(false);
  });

  it("dispatches updateProfile on form submission", async () => {
    const wrapper = mount(Profile, {
      global: {
        plugins: [store],
        mocks: { $router: { push: routerPushMock } },
        stubs: { ProfileCard: true },
      },
    });

    // Enter edit mode.
    await wrapper.find('[data-test="edit-button"]').trigger("click");

    // Change the bio field.
    const bioTextarea = wrapper.find("textarea#bio");
    await bioTextarea.setValue("Updated bio");

    // Submit the form.
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
        mocks: { $router: { push: routerPushMock } },
        stubs: {
          // Stub ProfileCard to render a simplified version of the skill.
          ProfileCard: {
            template:
              '<div class="profile-card-stub">{{ title }} - {{ description }}</div>',
            props: ["title", "description"],
          },
        },
      },
    });

    // The created hook in Profile.vue sets dummy skills: "Go Programming" and "Vue.js"
    const skillsSection = wrapper.find(".my-skills");
    expect(skillsSection.exists()).toBe(true);
    expect(skillsSection.text()).toContain("Go Programming");
    expect(skillsSection.text()).toContain("Vue.js");
  });
});
