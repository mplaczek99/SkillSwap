import { mount } from '@vue/test-utils'
import Search from '@/components/Search.vue'

describe('Search.vue', () => {
  it('renders the search input and button', () => {
    const wrapper = mount(Search)
    expect(wrapper.find('input').exists()).toBe(true)
    expect(wrapper.find('button').exists()).toBe(true)
  })

  it('filters dummy data based on the query and displays results', async () => {
    const wrapper = mount(Search)
    // Set query to "alice"
    await wrapper.find('input').setValue('alice')
    await wrapper.find('form').trigger('submit.prevent')
    // The dummy data includes "Alice"
    expect(wrapper.vm.results.length).toBeGreaterThan(0)
    expect(wrapper.text()).toContain('Alice')
  })

  it('shows a "No results found." message when query matches nothing', async () => {
    const wrapper = mount(Search)
    await wrapper.find('input').setValue('zzz')
    await wrapper.find('form').trigger('submit.prevent')
    expect(wrapper.text()).toContain('No results found.')
  })
})

