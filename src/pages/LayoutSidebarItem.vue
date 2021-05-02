<template>
  <div class="">
    <router-link :to="to" class="">
      <font-awesome-icon fixed-width :icon="icon"/>
      <span class="ml-1">{{ title }}</span>
    </router-link>
  </div>
</template>

<script>
export default {
  props: {
    icon: {
      type: String, // fontawesome icon class name
      required: true,
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
