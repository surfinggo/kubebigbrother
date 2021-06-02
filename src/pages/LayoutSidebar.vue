<template>
  <div class="fixed text-gray-500 h-screen top-0 pt-16">
    <div class="h-full overflow-y-scroll px-3">
      <div class="h-8 pointer-events-none absolute inset-x-0 z-10 px-3">
        <div class="h-8 bg-gradient-to-b from-white"></div>
      </div>
      <div class="py-3">
        <sidebar-item icon="tachometer-alt" title="Dashboard" :to="{name: 'dashboard'}"/>

        <sidebar-divider/>
      </div>
    </div>
  </div>
</template>

<script>
import SidebarDivider from './LayoutSidebarDivider.vue'
import SidebarItem from './LayoutSidebarItem.vue'

export default {
  components: {
    SidebarDivider,
    SidebarItem,
  }
}
</script>
