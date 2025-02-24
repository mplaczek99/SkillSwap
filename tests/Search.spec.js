import { mount } from "@vue/test-utils";
import flushPromises from "flush-promises";
import Search from "@/components/Search.vue";
import axios from "axios";

beforeAll(() => {
  jest.spyOn(console, 'error').mockImplementation(() => {});
});

jest.mock("axios");

describe("Search.vue", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it("calls the API and displays results when query matches a skill", async () => {
    const dummyResults = [
      { id: 1, name: "Dummy Skill", description: "This is a dummy skill" },
    ];
    axios.get.mockResolvedValue({ data: dummyResults });

    const wrapper = mount(Search, {
      props: { forceApiCall: true },
    });

    const input = wrapper.find("input");
    await input.setValue("dummy");
    await wrapper.find("form").trigger("submit.prevent");

    // Flush debounce: force immediate execution
    wrapper.vm.debouncedSearch.flush();
    await flushPromises();

    expect(axios.get).toHaveBeenCalledWith("/api/search", { params: { q: "dummy" } });
    expect(wrapper.text()).toContain("Dummy Skill");
  });

  it("displays 'No results found.' when API returns an empty array", async () => {
    axios.get.mockResolvedValue({ data: [] });

    const wrapper = mount(Search, {
      props: { forceApiCall: true },
    });
    const input = wrapper.find("input");
    await input.setValue("nonexistent");
    await wrapper.find("form").trigger("submit.prevent");

    // Flush debounce
    wrapper.vm.debouncedSearch.flush();
    await flushPromises();

    expect(axios.get).toHaveBeenCalledWith("/api/search", { params: { q: "nonexistent" } });
    expect(wrapper.text()).toContain("No results found.");
  });

  it("handles API errors gracefully", async () => {
    axios.get.mockRejectedValue(new Error("Network Error"));

    const wrapper = mount(Search, {
      props: { forceApiCall: true },
    });
    const input = wrapper.find("input");
    await input.setValue("dummy");
    await wrapper.find("form").trigger("submit.prevent");

    // Flush debounce
    wrapper.vm.debouncedSearch.flush();
    await flushPromises();

    // In case of an error, the component should show "No results found."
    expect(wrapper.text()).toContain("No results found.");
  });
});

