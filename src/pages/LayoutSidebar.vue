<template>
  <div class="fixed text-gray-500 h-screen top-0 pt-16 w-30 sm:w-60 xl:w-72">
    <div class="h-full overflow-y-scroll px-3">
      <div class="h-6 pointer-events-none absolute inset-x-0 z-10 px-3">
        <div class="h-6 bg-gradient-to-b from-white"></div>
      </div>
      <div class="py-6">
        <div
            class="px-3 mb-3 tracking-wide font-semibold text-sm lg:text-xs text-gray-900">
          <router-link :to="{name: 'dashboard'}">
            OVERVIEW
          </router-link>
        </div>

        <div
            class="px-3 mt-6 mb-3 tracking-wide font-semibold text-sm lg:text-xs text-gray-900">
          CLUSTER WATCHERS:
        </div>
        <router-link v-for="watcher in clusterwatchers" :key="watcher.metadata.name"
                     class="w-full px-3 mb-3 transition-colors duration-200 relative block hover:text-gray-900 text-gray-500 cursor-pointer"
                     :to="{name: 'clusterwatcher', params: {name: watcher.metadata.name}}">
          <span class="ml-1">{{ watcher.metadata.name }}</span>
        </router-link>
        <div
            class="px-3 mt-6 mb-3 tracking-wide font-semibold text-sm lg:text-xs text-gray-900">
          WATCHERS:
        </div>
        <router-link v-for="watcher in watchers" :key="watcher.metadata.name"
                     class="w-full px-3 mb-3 transition-colors duration-200 relative block hover:text-gray-900 text-gray-500 cursor-pointer"
                     :to="{name: 'watcher', params: {namespace: watcher.metadata.namespace, name: watcher.metadata.name}}">
          <span class="ml-1">{{ watcher.metadata.namespace }}/{{ watcher.metadata.name }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    watchers: {},
    clusterwatchers: {},
  },
  methods: {}
}
</script>
