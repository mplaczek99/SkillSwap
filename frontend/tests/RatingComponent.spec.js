import { mount } from "@vue/test-utils";
import RatingComponent from "@/components/RatingComponent.vue";

describe("RatingComponent.vue", () => {
  it("renders correctly in view-only mode", () => {
    const wrapper = mount(RatingComponent, {
      props: {
        value: 4,
        count: 10,
        viewOnly: true,
      },
      global: {
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Should be in rating display mode
    expect(wrapper.find(".rating-display").exists()).toBe(true);
    expect(wrapper.find(".rating-input").exists()).toBe(false);

    // Should display count
    expect(wrapper.find(".rating-count").text()).toBe("(10)");
  });

  it("renders correctly in interactive mode", () => {
    const wrapper = mount(RatingComponent, {
      props: {
        value: 3,
        viewOnly: false,
        title: "Rate this skill",
      },
      global: {
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Should be in interactive mode
    expect(wrapper.find(".rating-input").exists()).toBe(true);
    expect(wrapper.find(".rating-display").exists()).toBe(false);

    // Should display title
    expect(wrapper.find(".rating-title").text()).toBe("Rate this skill");
  });

  it("emits rating when stars are clicked", async () => {
    const wrapper = mount(RatingComponent, {
      props: {
        value: 0,
        viewOnly: false,
      },
      global: {
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Directly call the setRating method instead of trying to click stars
    // which may be problematic due to the font-awesome-icon stub
    await wrapper.vm.setRating(3);

    // Check emitted events
    expect(wrapper.emitted()).toHaveProperty("input");
    expect(wrapper.emitted().input[0]).toEqual([3]);
  });

  it("emits submit event with rating and feedback", async () => {
    const wrapper = mount(RatingComponent, {
      props: {
        value: 0,
        viewOnly: false,
        feedbackEnabled: true,
      },
      global: {
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Set rating
    await wrapper.setData({ currentRating: 5 });

    // Set feedback
    const textarea = wrapper.find("textarea");
    await textarea.setValue("Great experience!");

    // Click submit button
    const submitButton = wrapper.find("button.btn-primary");
    await submitButton.trigger("click");

    // Check emitted events
    expect(wrapper.emitted()).toHaveProperty("submit");
    expect(wrapper.emitted().submit[0][0]).toEqual({
      rating: 5,
      feedback: "Great experience!",
    });
  });

  it("disables submit button when required fields are missing", async () => {
    const wrapper = mount(RatingComponent, {
      props: {
        value: 0,
        viewOnly: false,
        feedbackEnabled: true,
        feedbackRequired: true,
      },
      global: {
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Set rating but no feedback
    await wrapper.setData({ currentRating: 4, feedback: "" });

    // Submit button should be disabled
    const submitButton = wrapper.find("button.btn-primary");
    expect(submitButton.attributes("disabled")).toBeDefined();

    // Add feedback
    const textarea = wrapper.find("textarea");
    await textarea.setValue("Nice work!");

    // Submit button should be enabled now
    expect(submitButton.attributes("disabled")).toBeUndefined();
  });

  it("shows correct star states based on rating", () => {
    const wrapper = mount(RatingComponent, {
      props: {
        value: 3.5, // Testing half star
        viewOnly: true,
      },
      global: {
        stubs: {
          "font-awesome-icon": true,
        },
      },
    });

    // Check classes (since we're using a stub, check the method instead)
    expect(wrapper.vm.getStarClass(1)).toBe("filled");
    expect(wrapper.vm.getStarClass(2)).toBe("filled");
    expect(wrapper.vm.getStarClass(3)).toBe("filled");
    expect(wrapper.vm.getStarClass(4)).toBe("half-filled");
    expect(wrapper.vm.getStarClass(5)).toBe("empty");
  });
});
