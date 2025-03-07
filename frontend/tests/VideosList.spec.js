import { mount, flushPromises } from "@vue/test-utils";
import VideosList from "@/components/VideosList.vue";
import axios from "axios";

jest.mock("axios");

describe("VideosList.vue", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("fetches and displays videos on creation", async () => {
    // Mock API response
    axios.get.mockResolvedValue({
      data: [
        {
          id: "video1.mp4",
          name: "video1.mp4",
          size: 1024 * 1024 * 10, // 10MB
          hasThumbnail: true,
          thumbnail: "video1.mp4.jpg",
          uploadedAt: new Date().toISOString(),
        },
        {
          id: "video2.mp4",
          name: "video2.mp4",
          size: 1024 * 1024 * 5, // 5MB
          hasThumbnail: false,
          uploadedAt: new Date().toISOString(),
        },
      ],
    });

    const wrapper = mount(VideosList, {
      global: {
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
        mocks: {
          $router: {
            push: jest.fn(),
          },
        },
      },
    });

    // Initially should be in loading state
    expect(wrapper.find(".loading-state").exists()).toBe(true);

    await flushPromises();

    // After loading, videos should be displayed
    expect(wrapper.findAll(".video-card").length).toBe(2);
    expect(wrapper.text()).toContain("video1");
    expect(wrapper.text()).toContain("video2");

    // Check if correct API endpoint was called
    expect(axios.get).toHaveBeenCalledWith("/api/videos");
  });

  it("shows empty state when no videos", async () => {
    // Mock empty response
    axios.get.mockResolvedValue({ data: [] });

    const wrapper = mount(VideosList, {
      global: {
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    await flushPromises();

    // Empty state should be displayed
    expect(wrapper.find(".empty-state").exists()).toBe(true);
    expect(wrapper.text()).toContain("No Videos Found");
  });

  it("handles API errors correctly", async () => {
    // Mock API error
    axios.get.mockRejectedValue(new Error("Failed to fetch videos"));

    const wrapper = mount(VideosList, {
      global: {
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    await flushPromises();

    // Error message should be displayed
    expect(wrapper.find(".error-message").exists()).toBe(true);
    expect(wrapper.text()).toContain("Failed to load videos");
  });

  it("formats file name correctly", () => {
    axios.get.mockResolvedValue({ data: [] });

    const wrapper = mount(VideosList, {
      global: {
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Test file name formatting
    expect(wrapper.vm.formatFileName("test_video.mp4")).toBe("test video");
    expect(wrapper.vm.formatFileName("my-awesome-video.mov")).toBe(
      "my awesome video",
    );
  });

  it("formats timestamps correctly", () => {
    axios.get.mockResolvedValue({ data: [] });

    const wrapper = mount(VideosList, {
      global: {
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Create a date to test with
    const testDate = new Date("2023-01-15T12:30:45");

    // Check if the date is formatted
    const formatted = wrapper.vm.formatDate(testDate.toISOString());
    expect(formatted).toContain("Jan"); // Should include month abbreviation
    expect(formatted).toContain("15"); // Should include day
    expect(formatted).toContain("2023"); // Should include year
  });

  it("handles playVideo functionality", () => {
    // Mock window.open
    const originalOpen = window.open;
    window.open = jest.fn().mockReturnValue({
      document: {
        write: jest.fn(),
      },
    });

    axios.get.mockResolvedValue({
      data: [
        {
          id: "video1.mp4",
          name: "video1.mp4",
          size: 1024 * 1024,
          hasThumbnail: false,
          uploadedAt: new Date().toISOString(),
        },
      ],
    });

    const wrapper = mount(VideosList, {
      global: {
        stubs: {
          "font-awesome-icon": true,
          "router-link": true,
        },
      },
    });

    // Test play video function
    wrapper.vm.playVideo({
      name: "test-video.mp4",
    });

    // Check if window.open was called
    expect(window.open).toHaveBeenCalled();

    // Restore original function
    window.open = originalOpen;
  });
});
