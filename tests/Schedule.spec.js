import { mount, flushPromises } from "@vue/test-utils";
import Schedule from "@/components/Schedule.vue";
import axios from "axios";

jest.mock("axios");

beforeAll(() => {
  jest.spyOn(console, "error").mockImplementation(() => {});
});

afterAll(() => {
  // Optionally restore console.error if needed
  console.error.mockRestore();
});

describe("Schedule.vue", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it("creates a schedule successfully", async () => {
    const dummySchedule = {
      id: 1,
      skill_id: 1,
      startTime: new Date().toISOString(),
      endTime: new Date(Date.now() + 3600000).toISOString(),
    };
    axios.post.mockResolvedValue({ data: dummySchedule });
    const wrapper = mount(Schedule);
    await wrapper.find('input[type="datetime-local"]').setValue(dummySchedule.startTime);
    const inputs = wrapper.findAll('input[type="datetime-local"]');
    await inputs[1].setValue(dummySchedule.endTime);
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();
    expect(wrapper.text()).toContain(`Session on Skill ID: ${dummySchedule.skill_id}`);
  });

  it("displays error message on schedule creation failure", async () => {
    axios.post.mockRejectedValue(new Error("Creation error"));
    const wrapper = mount(Schedule);
    await wrapper.find('input[type="datetime-local"]').setValue(new Date().toISOString());
    const inputs = wrapper.findAll('input[type="datetime-local"]');
    await inputs[1].setValue(new Date(Date.now() + 3600000).toISOString());
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();
    expect(wrapper.text()).toContain("Failed to create schedule");
  });
});
