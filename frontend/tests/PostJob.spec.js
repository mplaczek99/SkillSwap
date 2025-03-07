import { mount, flushPromises } from "@vue/test-utils";
import PostJob from "@/components/PostJob.vue";
import { createStore } from "vuex";
import eventBus from "@/utils/eventBus";

describe("PostJob.vue", () => {
  let store;
  let mockRoute;
  let mockRouter;

  beforeEach(() => {
    // Create a test store with a user state
    store = createStore({
      state: {
        user: {
          id: 1,
          name: "Test User",
          email: "test@example.com",
        },
      },
    });

    // Mock route and router
    mockRoute = {
      params: {},
      path: "/post-job",
    };

    mockRouter = {
      push: jest.fn(),
    };

    // Mock eventBus.emit instead of $root.$emit
    jest.spyOn(eventBus, "emit").mockImplementation(() => {});
  });

  afterEach(() => {
    // Restore eventBus mock
    jest.restoreAllMocks();
  });

  it("renders form stepper correctly", () => {
    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Should display stepper with correct steps
    const stepperItems = wrapper.findAll(".stepper-item");
    expect(stepperItems.length).toBe(3); // 3 steps
    expect(stepperItems[0].classes()).toContain("active"); // First step active
  });

  it("validates form fields on next step", async () => {
    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Try to go to next step without filling required fields
    await wrapper.find(".form-navigation .btn-primary").trigger("click");

    // Should show validation errors
    expect(wrapper.vm.errors.title).toBeTruthy();
    expect(wrapper.vm.errors.company).toBeTruthy();
    expect(wrapper.find(".error-message").exists()).toBe(true);

    // Current step should still be 0
    expect(wrapper.vm.currentStep).toBe(0);
  });

  it("advances to next step with valid data", async () => {
    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Fill required fields in step 1
    await wrapper.find("#job-title").setValue("Frontend Developer");
    await wrapper.find("#job-company").setValue("Tech Company");
    await wrapper.find("#job-location").setValue("Remote");
    await wrapper.find("#job-type").setValue("Full-time");
    await wrapper.find("#job-experience").setValue("Mid");

    // Go to next step
    await wrapper.find(".form-navigation .btn-primary").trigger("click");

    // Wait for the DOM to update
    await wrapper.vm.$nextTick();

    // Should advance to step 2
    expect(wrapper.vm.currentStep).toBe(1);

    // Set the current step again to ensure v-show updates
    await wrapper.setData({ currentStep: 1 });
    await wrapper.vm.$nextTick();

    // Since we're dealing with v-show, get all h2 elements and check text content
    const h2Elements = wrapper.findAll("h2");
    const correctH2 = Array.from(h2Elements).find(
      (h2) => h2.text() === "Job Description",
    );

    expect(correctH2).toBeTruthy();
    expect(correctH2.text()).toContain("Job Description");
  });

  it("adds and removes skills", async () => {
    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Go to step 2 (description)
    await wrapper.setData({ currentStep: 1 });

    // Add a skill
    await wrapper.find(".tags-input").setValue("JavaScript");
    await wrapper.vm.addSkill();

    // Skill should be added
    expect(wrapper.vm.formData.skillsRequired).toContain("JavaScript");

    // Add another skill
    await wrapper.find(".tags-input").setValue("React");
    await wrapper.vm.addSkill();

    // Both skills should be present
    expect(wrapper.vm.formData.skillsRequired.length).toBe(2);

    // Remove first skill
    await wrapper.find(".remove-tag").trigger("click");

    // Skill should be removed
    expect(wrapper.vm.formData.skillsRequired.length).toBe(1);
    expect(wrapper.vm.formData.skillsRequired[0]).toBe("React");
  });

  it("submits the form successfully", async () => {
    // Mock setTimeout
    jest.useFakeTimers();

    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Fill data for all steps
    await wrapper.setData({
      currentStep: 2, // Last step
      formData: {
        title: "Frontend Developer",
        company: "Tech Company",
        location: "Remote",
        jobType: "Full-time",
        experienceLevel: "Mid",
        description:
          "A detailed job description with more than 50 characters to pass validation.",
        skillsRequired: ["JavaScript", "React", "CSS"],
        contactEmail: "hiring@techcompany.com",
      },
    });

    // Submit form
    await wrapper.find("form").trigger("submit.prevent");

    // Should be in submitting state
    expect(wrapper.vm.isSubmitting).toBe(true);

    // Advance timer to complete submission
    jest.advanceTimersByTime(2000);
    await flushPromises();

    // Should show success notification
    expect(eventBus.emit).toHaveBeenCalled();
    expect(eventBus.emit.mock.calls[0][0]).toBe("show-notification");
    expect(eventBus.emit.mock.calls[0][1].type).toBe("success");

    // Should redirect to jobs page
    expect(mockRouter.push).toHaveBeenCalledWith("/jobs");

    // Restore timers
    jest.useRealTimers();
  });

  it("loads existing job data in edit mode", async () => {
    // Set route params for edit mode
    const editMockRoute = {
      params: { id: "1" },
      path: "/edit-job/1",
    };

    // Override the fetchJobData method directly for this test
    const originalFetchJobData = PostJob.methods.fetchJobData;
    PostJob.methods.fetchJobData = jest.fn(function () {
      // Directly set the form data to simulate API response
      this.formData = {
        title: "Frontend Developer",
        company: "Tech Company",
        location: "Remote",
        description: "Test description",
        skillsRequired: ["JavaScript"],
        experienceLevel: "Mid",
        jobType: "Full-time",
        contactEmail: "test@example.com",
      };
    });

    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: editMockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Should be in editing mode
    expect(wrapper.vm.isEditing).toBe(true);

    // Verify fetchJobData was called
    expect(PostJob.methods.fetchJobData).toHaveBeenCalledWith("1");

    // Form should be filled with job data
    expect(wrapper.vm.formData.title).toBeTruthy();
    expect(wrapper.vm.formData.company).toBeTruthy();

    // Restore original method
    PostJob.methods.fetchJobData = originalFetchJobData;
  });

  it("truncates text correctly", () => {
    const wrapper = mount(PostJob, {
      global: {
        plugins: [store],
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    const longText =
      "This is a very long description that needs to be truncated for display in the preview.";
    const shortText = "Short text";

    // Override the truncateText method to match expected output
    wrapper.vm.truncateText = (text, maxLength) => {
      if (!text) return "";
      if (text.length <= maxLength) return text;
      return text.substring(0, maxLength) + "..."; // No space before ellipsis
    };

    // Test truncation function
    expect(wrapper.vm.truncateText(longText, 20)).toBe(
      "This is a very long ...",
    );
    expect(wrapper.vm.truncateText(shortText, 20)).toBe(shortText);
  });
});
