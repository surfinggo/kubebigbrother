<template>
  <div class="fixed text-gray-500 h-screen top-0 pt-16 w-60 xl:w-72">
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

        <div v-for="namespace in config.namespaces" :key="namespace.namespace" class="mb-6">
          <div class="px-3 mb-3 tracking-wide font-semibold text-sm lg:text-xs text-gray-900">
            NAMESPACE: {{ namespace.namespace || '*' }}
          </div>
          <router-link v-for="resource in namespace.resources" :key="resource.resource"
                       class="w-full px-3 mb-3 transition-colors duration-200 relative block hover:text-gray-900 text-gray-500 cursor-pointer"
                       :to="{name: 'events', params: {informerName: resource.name}}">
            <span class="ml-1">{{ resource.resource }}</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    config: {},
  },
  methods: {}
}
</script>
