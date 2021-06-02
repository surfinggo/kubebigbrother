<template>
  <div class="sticky top-0 w-full z-50 flex border-b border-gray-200 bg-white h-16">
    <div class="fcc w-60 xl:w-72 px-3">
      <a class="overflow-hidden" href="/">
        <img :src="icon"/>
      </a>
    </div>
    <div class="fcb flex-grow px-3">
      <div>
        <div class="font-medium flex items-center space-x-4
            text-gray-400 hover:text-gray-600 focus-within:text-gray-600 transition-colors duration-200 w-full py-2">
          <font-awesome-icon icon="search"/>
          <input v-model="q" ref="q" placeholder="Search events"
                 class="outline-none"/>
          <span class="text-sm py-0.5 px-1.5 border border-gray-300 rounded-md">
          <span class="sr-only">Press</span>
          <kbd class="font-sans"><abbr title="Command" class="no-underline">⌘</abbr></kbd>
          <span class="sr-only">and</span>
          <kbd class="font-sans">K</kbd>
          <span class="sr-only">to search</span>
        </span>
        </div>
      </div>
      <github-icon/>
    </div>
  </div>
</template>

<script lang="ts">
import GithubIcon from '../components/GithubIcon.vue'
import Icon from '../../public/icon-text-right.png'

export default {
  components: {GithubIcon},
  data() {
    return {
      q: '',
      icon: Icon,
    }
  },
  created() {
    window.addEventListener('keydown', this.onkey)
  },
  beforeUnmount() {
    window.removeEventListener('keydown', this.onkey)
  },
  methods: {
    onkey(event: KeyboardEvent) {
      console.log(event)
      if (event.key === 'k') {
        if (event.metaKey || event.ctrlKey) {
          this.$refs.q && this.$refs.q.focus()
        }
      }
    },
  }
}
</script>

