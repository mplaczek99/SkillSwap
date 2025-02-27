import { mount } from "@vue/test-utils";
import Chat from "@/components/Chat.vue";

describe("Chat.vue", () => {
  it("renders chat messages", () => {
    const wrapper = mount(Chat);
    expect(wrapper.text()).toContain("Alice:");
    expect(wrapper.text()).toContain("Hi, how are you?");
  });

  it("adds a new message on form submit", async () => {
    const wrapper = mount(Chat);
    const input = wrapper.find("input");
    await input.setValue("New message");
    await wrapper.find("form").trigger("submit.prevent");
    expect(wrapper.text()).toContain("New message");
  });
});

