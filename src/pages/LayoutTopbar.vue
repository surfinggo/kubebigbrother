<template>
  <div class="fixed top-0 w-full z-1 flex border-b border-gray-200 bg-white h-16">
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
                 class="outline-none" @keydown.enter="search"/>
          <span class="text-sm py-0.5 px-1.5 border border-gray-300 rounded-md">
          <span class="sr-only">Press</span>
          <kbd class="font-sans"><abbr title="Command" class="no-underline">⌘</abbr></kbd>
          <span class="sr-only">and</span>
          <kbd class="font-sans">K</kbd>
          <span class="sr-only">to search</span>
        </span>
        </div>
      </div>
      <font-awesome-icon icon="file-alt"
                         class="text-gray-400 hover:text-gray-600 transition-colors duration-200 cursor-pointer"
                         @click="openModal"/>
    </div>
  </div>
  <TransitionRoot appear :show="isOpen" as="template">
    <Dialog as="div" @close="closeModal">
      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="min-h-screen px-4 text-center">
          <TransitionChild
              as="template"
              enter="duration-300 ease-out"
              enter-from="opacity-0"
              enter-to="opacity-50"
              leave="duration-200 ease-in"
              leave-from="opacity-50"
              leave-to="opacity-0">
            <DialogOverlay class="fixed inset-0 bg-black opacity-30"/>
          </TransitionChild>

          <span class="inline-block h-screen align-middle" aria-hidden="true">
            &#8203;
          </span>

          <TransitionChild
              as="template"
              enter="duration-300 ease-out"
              enter-from="opacity-0 scale-95"
              enter-to="opacity-100 scale-100"
              leave="duration-200 ease-in"
              leave-from="opacity-100 scale-100"
              leave-to="opacity-0 scale-95">
            <div
                class="inline-block w-full max-w-4xl p-4 my-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
              <DialogTitle as="h3" class="fcb text-lg font-medium leading-6 text-gray-900">
                <span>Informers Config</span>
                <JsonYamlSwitch v-model="useYaml"/>
              </DialogTitle>
              <div class="mt-2">
                <prism v-if="useYaml" language="yaml"
                       class="overflow-y-scroll rounded !text-sm"
                       style="max-height: calc(100vh - 12rem)">{{ configYaml }}
                </prism>
                <prism v-else language="json"
                       class="overflow-y-scroll rounded !text-sm"
                       style="max-height: calc(100vh - 12rem)">{{ config }}
                </prism>
              </div>
              <div class="mt-4 text-right">
                <button
                    type="button"
                    class="inline-flex justify-center px-4 py-2 text-sm font-medium text-blue-900 bg-blue-100 border border-transparent rounded-md hover:bg-blue-200 focus:outline-none"
                    @click="closeModal">
                  Done
                </button>
              </div>
            </div>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script lang="ts">
import Icon from '/icon-text-right.png'
import {ref} from 'vue'
import {Dialog, DialogOverlay, DialogTitle, TransitionChild, TransitionRoot} from '@headlessui/vue'
import JsonYamlSwitch from '../components/JsonYamlSwitch.vue'

export default {
  components: {
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogOverlay,
    DialogTitle,
    JsonYamlSwitch,
  },
  props: {
    config: {},
    configYaml: {},
  },
  setup() {
    const isOpen = ref(false)

    return {
      isOpen,
      closeModal() {
        isOpen.value = false
      },
      openModal() {
        isOpen.value = true
      },
    }
  },
  data() {
    return {
      q: '',
      icon: Icon,
      useYaml: false,
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
      if (event.key === 'k') {
        if (event.metaKey || event.ctrlKey) {
          this.$refs.q && this.$refs.q.focus()
        }
      }
    },
    search() {
      this.$router.push({
        name: 'search', query: {
          'q': this.q,
        }
      })
    }
  }
}
</script>

