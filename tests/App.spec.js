import { shallowMount } from '@vue/test-utils'
import App from '@/App.vue'
import Navbar from '@/components/Navbar.vue'

describe('App.vue', () => {
  it('renders Navbar component and router-view', () => {
    const wrapper = shallowMount(App, {
      global: {
        // Stub router-view so it renders a simple placeholder
        stubs: {
          'router-view': { template: '<div class="router-view-stub"></div>' }
        }
      }
    })
    expect(wrapper.findComponent(Navbar).exists()).toBe(true)
    expect(wrapper.find('.router-view-stub').exists()).toBe(true)
  })
})

