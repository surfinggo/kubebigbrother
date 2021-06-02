<template>
  <router-link :to="to" class="">
    <div class="w-full px-3 py-2 transition-colors duration-200 relative block hover:text-gray-900 text-gray-500">
      <font-awesome-icon v-if="icon" fixed-width :icon="icon"/>
      <span class="ml-1">{{ title }}</span>
    </div>
  </router-link>
</template>

<script>
export default {
  props: {
    icon: {
      type: String, // fontawesome icon class name
    },
    title: {
      type: String,
      required: true,
    },
    to: {
      type: [Object, String], // route name, required if has no collapse items, ignored when collapseId appears
    },
  },
  computed: {}
}
</script>
