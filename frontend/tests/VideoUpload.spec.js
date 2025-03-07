import { mount, flushPromises } from "@vue/test-utils";
import VideoUpload from "@/components/VideoUpload.vue";
import axios from "axios";

jest.mock("axios");

describe("VideoUpload.vue", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders file selection UI correctly", () => {
    const wrapper = mount(VideoUpload);

    expect(wrapper.find(".file-label").exists()).toBe(true);
    expect(wrapper.find('input[type="file"]').exists()).toBe(true);
    expect(wrapper.find(".file-info").text()).toContain("Maximum size");
  });

  it("displays selected file details", async () => {
    const wrapper = mount(VideoUpload);

    // Mock file selection
    const file = new File(["dummy content"], "test-video.mp4", {
      type: "video/mp4",
    });

    // Simulate file selection
    const input = wrapper.find('input[type="file"]');
    Object.defineProperty(input.element, "files", {
      value: [file],
      writable: false,
    });

    await input.trigger("change");

    // File details should be displayed
    expect(wrapper.find(".selected-file-details").exists()).toBe(true);
    expect(wrapper.text()).toContain("test-video.mp4");
  });

  it("validates file size", async () => {
    const wrapper = mount(VideoUpload);

    // Create a file that exceeds size limit (100MB)
    const largeFile = new File([""], "large-video.mp4", { type: "video/mp4" });
    Object.defineProperty(largeFile, "size", { value: 105 * 1024 * 1024 });

    const input = wrapper.find('input[type="file"]');
    Object.defineProperty(input.element, "files", {
      value: [largeFile],
      writable: false,
    });

    await input.trigger("change");

    // Error message should be displayed
    expect(wrapper.find(".error-message").exists()).toBe(true);
    expect(wrapper.text()).toContain("File is too large");
  });

  it("submits file upload successfully", async () => {
    // Mock successful API response
    axios.post.mockResolvedValue({
      data: {
        message: "Video uploaded successfully!",
        file: "test-video.mp4",
      },
    });

    const wrapper = mount(VideoUpload);

    // Set a file
    const file = new File(["dummy content"], "test-video.mp4", {
      type: "video/mp4",
    });
    await wrapper.setData({ selectedFile: file });

    // Click upload button
    const uploadButton = wrapper.find(".upload-button");
    await uploadButton.trigger("click");

    // FormData assertions
    expect(axios.post).toHaveBeenCalledWith(
      "/api/videos/upload",
      expect.any(FormData),
      expect.objectContaining({
        headers: { "Content-Type": "multipart/form-data" },
        onUploadProgress: expect.any(Function),
      }),
    );

    await flushPromises();

    // Success message should be displayed
    expect(wrapper.find(".success-message").exists()).toBe(true);
    expect(wrapper.text()).toContain("Video uploaded successfully");
  });

  it("handles upload errors correctly", async () => {
    // Mock API error
    axios.post.mockRejectedValue({
      response: {
        data: {
          error: "Upload failed due to server error",
        },
      },
    });

    const wrapper = mount(VideoUpload);

    // Set a file
    const file = new File(["dummy content"], "test-video.mp4", {
      type: "video/mp4",
    });
    await wrapper.setData({ selectedFile: file });

    // Click upload button
    const uploadButton = wrapper.find(".upload-button");
    await uploadButton.trigger("click");

    await flushPromises();

    // Error message should be displayed
    expect(wrapper.find(".error-message").exists()).toBe(true);
    expect(wrapper.text()).toContain("Upload failed");
  });

  it("formats file size correctly", () => {
    const wrapper = mount(VideoUpload);

    // Test different file sizes
    expect(wrapper.vm.formatFileSize(500)).toBe("500 bytes");
    expect(wrapper.vm.formatFileSize(1500)).toBe("1.5 KB");
    expect(wrapper.vm.formatFileSize(1500000)).toBe("1.4 MB");
  });

  it("disables upload button when no file is selected", () => {
    const wrapper = mount(VideoUpload);

    // Upload button should be disabled
    const uploadButton = wrapper.find(".upload-button");
    expect(uploadButton.attributes("disabled")).toBeDefined();

    // Select a file
    wrapper.setData({
      selectedFile: new File([""], "test.mp4", { type: "video/mp4" }),
    });

    // Button should be enabled
    expect(uploadButton.attributes("disabled")).toBeFalsy();
  });
});
