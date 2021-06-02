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
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
//@ts-ignore
import {DateTime} from "luxon"

export default {
  data() {
    return {
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
    }
  }
}
</script>
