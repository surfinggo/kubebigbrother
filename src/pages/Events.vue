<template>
  <div>
    <div
        v-if="broken"
        class="transition-all duration-1000 text-right mt-3"
        :class="brokenBlinkState ? 'text-red-400' : 'text-red-700'"
    >
      <font-awesome-icon icon="unlink"/>
      sync error, please refresh the page
    </div>
    <div v-for="event in events" :key="event.id">
      <div class="py-8 flex flex-wrap md:flex-nowrap">
        <div class="md:w-64 md:mb-0 mb-6 flex-shrink-0 flex flex-col">
          <span class="font-semibold title-font" :class="typeColor(event.event_type)">{{ event.event_type }}</span>
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
                class="inline-block w-full max-w-4xl p-4 my-6 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
              <DialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">
                <div class="fcb">
                  <span>Event Detail</span>
                  <JsonYamlSwitch v-model="useYaml"/>
                </div>
              </DialogTitle>
              <div class="mt-2">
                <prism v-if="evtData && !useYaml" language="json"
                       class="rounded !text-sm"
                       style="max-height: calc(12rem)">
                  {{ evtData.event }}
                </prism>
                <prism v-if="evtData && useYaml" language="yaml"
                       class="rounded !text-sm"
                       style="max-height: calc(12rem)">
                  {{ atob(evtData.event_yaml) }}
                </prism>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div class="" :class="evtData && evtData.old_obj ? '' : 'col-span-2'">
                  <div class="mt-2 text-sm">Resource object:</div>
                  <div class="mt-2">
                    <prism v-if="evtData && !useYaml" language="json"
                           class="rounded !text-sm"
                           style="max-height: calc(100vh - 28rem)">{{ evtData.obj }}
                    </prism>
                    <prism v-if="evtData && useYaml" language="yaml"
                           class="rounded !text-sm"
                           style="max-height: calc(100vh - 28rem)">{{ atob(evtData.obj_yaml) }}
                    </prism>
                  </div>
                </div>
                <div>
                  <div class="mt-2 text-sm" v-if="evtData && evtData.old_obj">Old object (when updated):</div>
                  <div class="mt-2">
                    <prism v-if="evtData && evtData.old_obj && !useYaml" language="json"
                           class="rounded !text-sm"
                           style="max-height: calc(100vh - 28rem)">{{ evtData.old_obj }}
                    </prism>
                    <prism v-if="evtData && evtData.old_obj && useYaml" language="yaml"
                           class="rounded !text-sm"
                           style="max-height: calc(100vh - 28rem)">{{ atob(evtData.old_obj_yaml) }}
                    </prism>
                  </div>
                </div>
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
import JsonYamlSwitch from "../components/JsonYamlSwitch.vue";
import {Dialog, DialogOverlay, DialogTitle, TransitionChild, TransitionRoot,} from '@headlessui/vue'

export default {
  components: {
    JsonYamlSwitch,
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogOverlay,
    DialogTitle,
  },
  setup() {
    const isOpen = ref(false)
    const useYaml = ref(false)

    return {
      isOpen,
      useYaml,
    }
  },
  data() {
    return {
      evtData: null,
      getting: false,
      broken: false,
      brokenBlinkState: false,
      brokenBlink: null, // setInterval to change brokenBlink
      events: [],
      refresher: null,
    }
  },
  watch: {
    '$route.query.q': function () {
      clearInterval(this.refresher)
      this.refresh()
      this.refresher = setInterval(this.refresh, 3000)
    },
    '$route.params.informerName': function () {
      clearInterval(this.refresher)
      this.refresh()
      this.refresher = setInterval(this.refresh, 3000)
    }
  },
  created() {
    this.refresh()
    this.refresher = setInterval(this.refresh, 3000)
    this.brokenBlink = setInterval(() => {
      this.brokenBlinkState = !this.brokenBlinkState
    }, 800)
  },
  beforeUnmount() {
    clearInterval(this.brokenBlink)
    clearInterval(this.refresher)
  },
  methods: {
    closeModal() {
      this.isOpen = false
      this.evtData = null
    },
    openModal(e) {
      this.isOpen = true
      this.axios.get('/api/v1/events/' + e.id).then(r => {
        console.log(r)
        this.evtData = r.data
      })
    },
    refresh() {
      if (this.getting) {
        return
      }
      this.getting = true
      this.$http.get('/api/v1/events', {
        params: {
          informerName: this.$route.params.informerName
        }
      }).then(r => {
        this.events = r.data.events
      }, () => {
        this.broken = true
        clearInterval(this.refresher)
      }).finally(() => {
        this.getting = false
      })
    },
    typeColor(t) {
      return {
        'ADDED': 'text-green-500',
        'DELETED': 'text-yellow-500',
        'UPDATED': 'text-blue-500'
      }[t] || 'text-gray-500'
    },
    lux(t) {
      return DateTime.fromISO(t).toFormat("yyyy-MM-dd HH:mm:ss")
    },
    atob(v) {
      return atob(v)
    }
  }
}
</script>
