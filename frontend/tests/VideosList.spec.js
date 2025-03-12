import { mount, flushPromises } from "@vue/test-utils";
import VideosList from "@/components/VideosList.vue";
import axios from "axios";

// Mock axios
jest.mock("axios");

// Mock the eventBus
jest.mock("@/utils/eventBus", () => ({
  emit: jest.fn(),
}));

describe("VideosList.vue", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("fetches videos on created", () => {
    // Mock successful response
    axios.get.mockResolvedValue({ data: [] });

    // Mount component
    mount(VideosList);

    // Check that the API was called with correct parameters
    expect(axios.get).toHaveBeenCalledWith("/api/videos", { timeout: 15000 });
  });

  it("formats file size correctly", () => {
    // Mock successful response
    axios.get.mockResolvedValue({ data: [] });

    // Mount component to access methods
    const wrapper = mount(VideosList);

    // Test bytes formatting (under 1024 bytes)
    expect(wrapper.vm.formatFileSize(1000)).toBe("1000 bytes");

    // Test KB formatting
    expect(wrapper.vm.formatFileSize(1500)).toBe("1.5 KB");

    // Test MB formatting
    expect(wrapper.vm.formatFileSize(1500000)).toBe("1.4 MB");
  });

  it("displays videos when data is loaded", async () => {
    // Mock successful response with video data
    axios.get.mockResolvedValue({
      data: [
        {
          id: "test.mp4",
          name: "test.mp4",
          size: 1024,
          uploadedAt: new Date().toISOString(),
        },
      ],
    });

    // Mount component
    const wrapper = mount(VideosList);

    // Wait for promises to resolve
    await flushPromises();

    // Check that videos are displayed
    expect(wrapper.find(".videos-grid").exists()).toBe(true);
  });

  it("displays empty state when no videos found", async () => {
    // Mock empty response
    axios.get.mockResolvedValue({ data: [] });

    // Mount component
    const wrapper = mount(VideosList);

    // Wait for promises to resolve
    await flushPromises();

    // Check for empty state content (check for specific text)
    expect(wrapper.text()).toContain("No Videos Found");
  });

  it("displays error message when API call fails", async () => {
    // Mock error response
    const errorMessage = "Server error";
    axios.get.mockRejectedValue({
      response: {
        status: 500,
        data: { error: errorMessage },
      },
    });

    // Mount component
    const wrapper = mount(VideosList);

    // Wait for promises to resolve
    await flushPromises();

    // Check that error message is displayed (check for error element)
    expect(wrapper.find(".error-message").exists()).toBe(true);

    // Check that error state is set in component data
    expect(wrapper.vm.error).not.toBeNull();
  });
});
