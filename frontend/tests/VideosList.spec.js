import { mount, flushPromises } from "@vue/test-utils";
import VideosList from "@/components/VideosList.vue";
import axios from "axios";

// Mock axios
jest.mock("axios");

// Mock router
const mockRouter = {
  push: jest.fn(),
};

// Mock event bus
jest.mock("@/utils/eventBus", () => ({
  emit: jest.fn(),
  on: jest.fn(),
  off: jest.fn(),
}));

describe("VideosList.vue", () => {
  let wrapper;

  // Test data
  const mockVideos = [
    {
      id: "video1.mp4",
      name: "video1.mp4",
      originalFilename: "Test Video 1.mp4",
      hasThumbnail: true,
      thumbnail: "video1.mp4.jpg",
      size: 1024 * 1024 * 10, // 10MB
      uploadedAt: "2023-01-15T12:00:00Z",
    },
    {
      id: "video2.mp4",
      name: "video2.mp4",
      originalFilename: "Test Video 2.mp4",
      hasThumbnail: false,
      size: 1024 * 1024 * 5, // 5MB
      uploadedAt: "2023-01-10T10:00:00Z",
    },
  ];

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
  });

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount();
    }
  });

  it("displays video list when videos are loaded", async () => {
    // Mock successful API response
    axios.get.mockResolvedValueOnce({ data: mockVideos });

    // Mount component with mocked router
    wrapper = mount(VideosList, {
      global: {
        mocks: {
          $router: mockRouter,
        },
      },
    });

    // Wait for API call to resolve
    await flushPromises();

    // Check that videos are displayed
    expect(wrapper.findAll(".video-card")).toHaveLength(mockVideos.length);
  });

  it("displays empty state when no videos found", async () => {
    // Mock empty response
    axios.get.mockResolvedValueOnce({ data: [] });

    // Mount component
    wrapper = mount(VideosList, {
      global: {
        mocks: {
          $router: mockRouter,
        },
      },
    });

    await flushPromises();

    // Check for empty state content (update expected text to match the actual component)
    expect(wrapper.text()).toContain("Your Video Library is Empty");
  });

  it("displays error message when API call fails", async () => {
    // Mock API error
    axios.get.mockRejectedValueOnce(new Error("API Error"));

    // Mount component
    wrapper = mount(VideosList, {
      global: {
        mocks: {
          $router: mockRouter,
        },
      },
    });

    await flushPromises();

    // Check that error message is displayed
    expect(wrapper.find(".error-message").exists()).toBe(true);
  });
});
