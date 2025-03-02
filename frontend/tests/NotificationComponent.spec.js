import { mount } from "@vue/test-utils";
import NotificationComponent from "@/components/NotificationComponent.vue";
import eventBus from "@/utils/eventBus";

describe("NotificationComponent.vue", () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.spyOn(eventBus, "on").mockImplementation(() => {});
    jest.spyOn(eventBus, "off").mockImplementation(() => {});
  });

  afterEach(() => {
    jest.useRealTimers();
    eventBus.on.mockRestore();
    eventBus.off.mockRestore();
  });

  it("sets up and cleans up event listeners", () => {
    const wrapper = mount(NotificationComponent, {
      global: {
        stubs: { "font-awesome-icon": true }
      },
    });
    
    // Verify that the event bus registers the "show-notification" event listener.
    expect(eventBus.on).toHaveBeenCalledWith("show-notification", wrapper.vm.showNotification);
    
    // Unmount the component to trigger beforeUnmount.
    wrapper.unmount();
    
    expect(eventBus.off).toHaveBeenCalledWith("show-notification", wrapper.vm.showNotification);
  });
});

