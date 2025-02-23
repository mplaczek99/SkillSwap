import { mount } from "@vue/test-utils";
import Search from "@/components/Search.vue";
import flushPromises from "flush-promises";

describe("Search.vue", () => {
  it("filters dummy data based on the query and displays results", async () => {
    const wrapper = mount(Search);
    await wrapper.find("input").setValue("alice");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises(); // Wait for the async search to complete
    await wrapper.vm.$nextTick();
    expect(wrapper.vm.results.length).toBeGreaterThan(0);
    expect(wrapper.text()).toContain("Alice");
  });

  it('shows a "No results found." message when query matches nothing', async () => {
    const wrapper = mount(Search);
    await wrapper.find("input").setValue("zzz");
    await wrapper.find("form").trigger("submit.prevent");
    await flushPromises();
    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain("No results found.");
  });
});
