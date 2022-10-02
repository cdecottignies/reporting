import { shallowMount } from '@vue/test-utils'
import HomeComponent from '@/components/HomeComponent.vue'

describe('HomeComponent.vue', () => {
    it('renders props.msg when passed', () => {
        const msg = 'new message'
        const wrapper = shallowMount(HomeComponent, {
            propsData: { msg },
        })
        expect(wrapper.text()).toMatch(msg)
    })
})
