import { mount, flushPromises } from "@vue/test-utils";
import JobDetail from "@/components/JobDetail.vue";
import JobPost from "@/models/JobPost";

jest.mock("@/models/JobPost");

describe("JobDetail.vue", () => {
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
    };

    mockRoute = {
      params: { id: "1" },
    };

    // Mock setTimeout to make tests run faster
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("renders job details correctly", async () => {
    const wrapper = mount(JobDetail, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
      },
    });

    // Initially should be in loading state
    expect(wrapper.find(".loading-state").exists()).toBe(true);

    // Fast-forward timer to trigger the setTimeout in fetchJob
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // After loading, job details should be displayed
    expect(wrapper.find(".job-detail-page").exists()).toBe(true);
    expect(wrapper.find(".job-header").exists()).toBe(true);
    expect(wrapper.find(".job-description").exists()).toBe(true);
  });

  it("displays error when job is not found", async () => {
    // Set up route with non-existent job ID
    mockRoute.params.id = "999";

    // Mock getMockJobs to return an empty array
    const originalGetMockJobs = JobDetail.methods.getMockJobs;
    JobDetail.methods.getMockJobs = jest.fn().mockReturnValue([]);

    const wrapper = mount(JobDetail, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Should display error message
    expect(wrapper.find(".error-message").exists()).toBe(true);

    // Restore original method
    JobDetail.methods.getMockJobs = originalGetMockJobs;
  });

  it("refreshes when route changes", async () => {
    const wrapper = mount(JobDetail, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Mock fetchJob method
    const fetchJobMock = jest.fn();
    wrapper.vm.fetchJob = fetchJobMock;

    // Simulate route change
    await wrapper.vm.$options.watch.$route.call(
      wrapper.vm,
      { params: { id: "2" } }, // new route
      { params: { id: "1" } }, // old route
    );

    expect(fetchJobMock).toHaveBeenCalled();
  });

  it("renders similar jobs correctly", async () => {
    const wrapper = mount(JobDetail, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Check similar jobs section
    expect(wrapper.find(".similar-jobs-section").exists()).toBe(true);
    expect(
      wrapper.findAll(".similar-jobs-grid .job-card").length,
    ).toBeGreaterThan(0);
  });

  it("truncates description correctly", async () => {
    const wrapper = mount(JobDetail, {
      global: {
        mocks: {
          $route: mockRoute,
          $router: mockRouter,
        },
        stubs: {
          "font-awesome-icon": true,
          "router-link": {
            template: "<a><slot /></a>",
          },
        },
      },
    });

    // Fast-forward timer
    jest.advanceTimersByTime(1000);
    await flushPromises();

    // Test the truncateDescription method
    const shortDesc = "Short description";
    const longDesc =
      "This is a very long description that should be truncated because it exceeds the maximum length specified for the truncation function";

    expect(wrapper.vm.truncateDescription(shortDesc, 50)).toBe(shortDesc);
    expect(wrapper.vm.truncateDescription(longDesc, 50).length).toBeLessThan(
      longDesc.length,
    );
    expect(wrapper.vm.truncateDescription(longDesc, 50).endsWith("...")).toBe(
      true,
    );
  });
});
