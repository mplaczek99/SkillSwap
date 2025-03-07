import { mount, flushPromises } from "@vue/test-utils";
import JobPostings from "@/components/JobPostings.vue";
import JobPost from "@/models/JobPost";

jest.mock("@/models/JobPost");

describe("JobPostings.vue", () => {
  let mockRouter, mockRoute;

  beforeEach(() => {
    // Reset mock implementation
    jest.clearAllMocks();

    // Mock the JobPost model
    JobPost.mockImplementation((data) => {
      return {
        ...data,
        skillsArray: jest.fn().mockReturnValue(data.skillsRequired || []),
        formattedDate: jest.fn().mockReturnValue("January 1, 2023"),
        daysSincePosting: jest.fn().mockReturnValue(5),
      };
    });

    // Set up mock route and router
    mockRouter = {
      push: jest.fn(),
      replace: jest.fn(),
    };

    mockRoute = {
      query: {},
    };

    // Mock setTimeout to make tests run faster
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("renders job listings and filters correctly", async () => {
    const wrapper = mount(JobPostings, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Initially should be in loading state
    expect(wrapper.find(".loading-state").exists()).toBe(true);

    // Fast-forward timer to trigger the setTimeout in fetchJobs
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // After loading, jobs should be displayed
    expect(wrapper.find(".job-card").exists()).toBe(true);

    // Test filtering
    const searchInput = wrapper.find(".search-input");
    await searchInput.setValue("frontend");

    // Directly call filterJobs instead of relying on debounced event
    await wrapper.vm.filterJobs();

    // Check that the search query was applied
    expect(mockRouter.replace).toHaveBeenCalled();
    expect(mockRouter.replace).toHaveBeenCalledWith(
      expect.objectContaining({
        query: expect.objectContaining({
          q: "frontend",
        }),
      }),
    );
  });

  it("applies filters correctly", async () => {
    const wrapper = mount(JobPostings, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Get the select element and set a value
    const filterValue = "Full-time";

    // Use setData to set the filter value directly
    await wrapper.setData({
      filters: {
        ...wrapper.vm.filters,
        jobType: filterValue,
      },
    });

    // Call filterJobs directly
    await wrapper.vm.filterJobs();

    // Check if the filter was applied correctly
    expect(mockRouter.replace).toHaveBeenCalled();
    expect(mockRouter.replace).toHaveBeenCalledWith(
      expect.objectContaining({
        query: expect.objectContaining({
          type: filterValue,
        }),
      }),
    );
  });

  it("resets filters correctly", async () => {
    // Create wrapper with existing filters
    mockRoute.query = { q: "developer", type: "Full-time" };

    const wrapper = mount(JobPostings, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Call resetFilters method directly
    await wrapper.vm.resetFilters();

    // Wait for reset to apply
    await flushPromises();

    // Verify filters were reset
    expect(mockRouter.replace).toHaveBeenCalled();
    const lastCall =
      mockRouter.replace.mock.calls[mockRouter.replace.mock.calls.length - 1];
    expect(lastCall[0].query.q).toBeUndefined();
    expect(lastCall[0].query.type).toBeUndefined();
  });

  it("handles empty search results", async () => {
    // Mount component
    const wrapper = mount(JobPostings, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Wait for initial load
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Manually set filtered jobs to empty array
    await wrapper.setData({
      filteredJobs: [],
      loading: false,
      jobs: [],
      searchQuery: "nonexistent",
    });

    // Verify empty state is shown
    expect(wrapper.find(".no-jobs").exists()).toBe(true);
    expect(wrapper.find(".job-card").exists()).toBe(false);
  });

  it("has the correct job type classes", async () => {
    const wrapper = mount(JobPostings, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Test the getJobTypeClass method
    expect(wrapper.vm.getJobTypeClass("Full-time")).toBe("full-time");
    expect(wrapper.vm.getJobTypeClass("Part-time")).toBe("part-time");
    expect(wrapper.vm.getJobTypeClass("Contract")).toBe("contract");
    expect(wrapper.vm.getJobTypeClass("Freelance")).toBe("freelance");
    expect(wrapper.vm.getJobTypeClass("Unknown")).toBe("");
  });
});
