<template>
  <div>
    <div v-for="event in events" :key="event.id">
      <div class="py-8 flex flex-wrap md:flex-nowrap">
        <div class="md:w-64 md:mb-0 mb-6 flex-shrink-0 flex flex-col">
          <span class="font-semibold title-font text-gray-700">{{ event.event_type }}</span>
          <span class="text-sm text-gray-500">ID: {{ event.id }}</span>
          <span class="text-sm text-gray-500">Time: {{ lux(event.create_time) }}</span>
        </div>
        <div class="md:flex-grow">
          <h2 class="text-xl font-medium text-gray-900 title-font mb-2">
            {{ event.kind }}
            <span class="font-bold text-indigo-500">{{ event.name }}</span>
            has been
            <span class="font-bold text-indigo-500">{{ event.event_type }}</span>
            in namespace
            <span class="font-bold text-indigo-500">{{ event.namespace }}</span>
          </h2>
          <span class="cursor-pointer text-indigo-500 inline-flex items-center mt-4"
                @click="openModal(event)">Detail
            <svg class="w-4 h-4 ml-2" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" fill="none"
                 stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 12h14"></path>
              <path d="M12 5l7 7-7 7"></path>
            </svg>
          </span>
        </div>
      </div>
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
                class="inline-block w-full max-w-xl p-4 my-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
              <DialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">
                Event Detail
              </DialogTitle>
              <div class="mt-2">
                <prism v-if="evt" language="json" class="rounded !text-sm">{{ evt }}</prism>
              </div>
              <div class="mt-2 text-sm">Resource object:</div>
              <div class="mt-2">
                <prism v-if="evt" language="json" class="rounded !text-sm">{{ base64decode(evt.obj) }}</prism>
              </div>
              <div class="mt-2 text-sm" v-if="evt && evt.old_obj">Old object (when updated):</div>
              <div class="mt-2">
                <prism v-if="evt && evt.old_obj" language="json" class="rounded !text-sm">{{
                    base64decode(evt.old_obj)
                  }}
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
//@ts-ignore
import {DateTime} from "luxon"
import {ref} from "vue";
import {Dialog, DialogOverlay, DialogTitle, TransitionChild, TransitionRoot,} from '@headlessui/vue'

export default {
  components: {
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogOverlay,
    DialogTitle,
  },
  setup() {
    const isOpen = ref(false)

    return {
      isOpen,
    }
  },
  data() {
    return {
      evt: null,
      events: [],
      refresher: null,
    }
  },
  watch: {
    '$route.params.informerName': function () {
      clearInterval(this.refresher)
      this.refresh()
      this.refresher = setInterval(this.refresh, 3000)
    }
  },
  created() {
    this.refresh()
    this.refresher = setInterval(this.refresh, 3000)
  },
  beforeUnmount() {
    clearInterval(this.refresher)
  },
  methods: {
    closeModal() {
      this.isOpen = false
    },
    openModal(e) {
      this.evt = e
      this.isOpen = true
    },
    refresh() {
      this.$http.get('/api/v1/events', {
        params: {
          informerName: this.$route.params.informerName
        }
      }).then(r => {
        this.events = r.data.events
      })
    },
    lux(t) {
      return DateTime.fromISO(t).toFormat("yyyy-MM-dd HH:mm:ss")
    },
    base64decode(bytes) {
      if (!bytes) {
        return ''
      }
      return JSON.parse(atob(bytes))
    }
  }
}
</script>
