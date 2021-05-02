<template>
  <div class="bg-purple-600 px-4 py-2 text-white" id="sidebar">
    <router-link class="" :to="{name: 'dashboard'}">
      <div class="">{{ $store.state.app.siteName }}</div>
      <div class="text-xs">{{ $store.state.app.siteDescription }}</div>
    </router-link>

    <sidebar-divider/>

    <!-- Nav Item - Dashboard -->
    <sidebar-item icon="tachometer-alt" title="Dashboard" :to="{name: 'dashboard'}"/>

    <sidebar-divider/>


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
