import { shallowMount } from "@vue/test-utils";
import VideosList from "@/components/VideosList.vue";
import axios from "axios";

// Mock axios
jest.mock("axios");

describe("VideosList.vue", () => {
  let wrapper;

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();

    // Mock axios response
    axios.get.mockResolvedValue({
      data: [],
    });

    // Set API URL
    axios.defaults.baseURL = "http://test-api";

    // Create wrapper with minimal mounting
    wrapper = shallowMount(VideosList, {
      attachTo: null, // Don't attach to DOM
      stubs: {
        // Stub out any child components that might cause issues
        "font-awesome-icon": true,
      },
    });
  });

  afterEach(() => {
    // Clean up
    wrapper.unmount();
  });

  it("renders correctly", () => {
    expect(wrapper.exists()).toBe(true);
  });

  it("fetches videos on created", () => {
    expect(axios.get).toHaveBeenCalledWith("/api/videos");
  });

  it("formats file size correctly", () => {
    expect(wrapper.vm.formatFileSize(500)).toBe("500 bytes");
    expect(wrapper.vm.formatFileSize(1500)).toBe("1.5 KB");
  });

  it("formats file name correctly", () => {
    expect(wrapper.vm.formatFileName("test-video.mp4")).toBe("test video");
    expect(wrapper.vm.formatFileName("my_awesome_video.mov")).toBe(
      "my awesome video",
    );
  });

  it("has playVideo method", () => {
    // Just test that the method exists, don't call it
    expect(typeof wrapper.vm.playVideo).toBe("function");
  });

  it("has downloadVideo method", () => {
    // Just test that the method exists, don't call it
    expect(typeof wrapper.vm.downloadVideo).toBe("function");
  });
});
