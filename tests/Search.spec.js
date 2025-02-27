import { mount } from "@vue/test-utils";
import flushPromises from "flush-promises";
import Search from "@/components/Search.vue";
import axios from "axios";

jest.mock("axios");
jest.spyOn(console, 'error').mockImplementation(() => {});

describe("Search.vue", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it("displays search results when query matches", async () => {
    const dummyResults = [
      { id: 1, name: "Dummy Skill", description: "This is a dummy skill" },
    ];
    axios.get.mockResolvedValue({ data: dummyResults });

    const wrapper = mount(Search, { props: { forceApiCall: true } });
    const input = wrapper.find("input");
    await input.setValue("dummy");
    await wrapper.find("form").trigger("submit.prevent");

    wrapper.vm.debouncedSearch.flush();
    await flushPromises();

    expect(axios.get).toHaveBeenCalledWith("/api/search", { params: { q: "dummy" } });
    expect(wrapper.text()).toContain("Dummy Skill");
  });

  it("shows no results when API returns empty array", async () => {
    axios.get.mockResolvedValue({ data: [] });

    const wrapper = mount(Search, { props: { forceApiCall: true } });
    const input = wrapper.find("input");
    await input.setValue("nonexistent");
    await wrapper.find("form").trigger("submit.prevent");

    wrapper.vm.debouncedSearch.flush();
    await flushPromises();

    expect(wrapper.text()).toContain("No results found.");
  });

  it("handles API errors gracefully", async () => {
    axios.get.mockRejectedValue(new Error("Network Error"));

    const wrapper = mount(Search, { props: { forceApiCall: true } });
    const input = wrapper.find("input");
    await input.setValue("dummy");
    await wrapper.find("form").trigger("submit.prevent");

    wrapper.vm.debouncedSearch.flush();
    await flushPromises();

    expect(wrapper.text()).toContain("An error occurred while searching");
  });
});

