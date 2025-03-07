import { mount } from "@vue/test-utils";
import FeedbackSystem from "@/components/FeedbackSystem.vue";
import eventBus from "@/utils/eventBus"; // Import the event bus

describe("FeedbackSystem.vue", () => {
  beforeEach(() => {
    // Use jest.spyOn on eventBus instead of mocking $root
    jest.spyOn(eventBus, "emit").mockImplementation(() => {});
  });

  afterEach(() => {
    // Restore the mock after tests
    eventBus.emit.mockRestore();
  });

  it("renders pending feedbacks correctly", () => {
    const wrapper = mount(FeedbackSystem, {
      global: {
        stubs: {
          RatingComponent: true,
          "font-awesome-icon": true,
        },
      },
    });

    // Should display pending feedback section
    expect(wrapper.find(".pending-feedback-section").exists()).toBe(true);
    expect(wrapper.findAll(".feedback-card").length).toBe(2); // From mock data
  });

  it("submits feedback correctly", async () => {
    const wrapper = mount(FeedbackSystem, {
      global: {
        stubs: {
          // Use a shallow stunt content for RatingComponent to prevent full rendering
          RatingComponent: {
            template: '<div class="rating-stub"></div>',
            props: ["title"],
            emits: ["submit", "cancel"],
          },
          "font-awesome-icon": true,
        },
      },
    });

    // Get the initial count of pending feedbacks
    const initialPendingCount = wrapper.vm.pendingFeedbacks.length;

    // Get first feedback card
    const firstFeedback = wrapper.vm.pendingFeedbacks[0];

    // Simulate submitting feedback
    await wrapper.vm.submitFeedback(firstFeedback.id, {
      rating: 5,
      feedback: "Great teacher!",
    });

    // Pending count should decrease
    expect(wrapper.vm.pendingFeedbacks.length).toBe(initialPendingCount - 1);

    // Given feedback should increase
    expect(wrapper.vm.givenFeedbacks.length).toBeGreaterThan(0);

    // Notification should be emitted
    expect(eventBus.emit).toHaveBeenCalled();
    expect(eventBus.emit.mock.calls[0][0]).toBe("show-notification");
  });

  it("skips feedback correctly", async () => {
    const wrapper = mount(FeedbackSystem, {
      global: {
        stubs: {
          RatingComponent: true,
          "font-awesome-icon": true,
        },
      },
    });

    // Get the initial count of pending feedbacks
    const initialPendingCount = wrapper.vm.pendingFeedbacks.length;

    // Skip first feedback
    await wrapper.vm.skipFeedback(wrapper.vm.pendingFeedbacks[0].id);

    // Pending count should decrease
    expect(wrapper.vm.pendingFeedbacks.length).toBe(initialPendingCount - 1);

    // Notification should be emitted
    expect(eventBus.emit).toHaveBeenCalled();
    expect(eventBus.emit.mock.calls[0][0]).toBe("show-notification");
  });

  it("renders received feedback correctly", () => {
    const wrapper = mount(FeedbackSystem, {
      global: {
        stubs: {
          RatingComponent: true,
          "font-awesome-icon": true,
        },
      },
    });

    // Should display received feedback section
    expect(wrapper.find(".received-feedback-section").exists()).toBe(true);
    expect(wrapper.findAll(".review-card").length).toBeGreaterThan(0);
  });

  it("renders empty state when no received feedback", async () => {
    // Create an instance with empty received feedback
    const wrapper = mount(FeedbackSystem, {
      data() {
        return {
          receivedFeedbacks: [],
        };
      },
      global: {
        stubs: {
          RatingComponent: true,
          "font-awesome-icon": true,
        },
      },
    });

    // Should display empty state
    expect(wrapper.find(".empty-feedback").exists()).toBe(true);
    expect(wrapper.text()).toContain("You haven't received any feedback yet");
  });

  it("formats dates correctly", () => {
    const wrapper = mount(FeedbackSystem, {
      global: {
        stubs: {
          RatingComponent: true,
          "font-awesome-icon": true,
        },
      },
    });

    // Test date formatting
    const testDate = new Date("2023-01-15T12:30:45");
    const formattedDate = wrapper.vm.formatDate(testDate);

    // Should include month abbreviation, day, and year
    expect(formattedDate).toContain("Jan");
    expect(formattedDate).toContain("15");
    expect(formattedDate).toContain("2023");
  });
});
