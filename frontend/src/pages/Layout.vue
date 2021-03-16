<!--
------------------wrapper-------------------------
|   --sidebar--   ------------main------------   |
|   |         |   |                          |   |
|   |         |   |   -div (flex-grow: 1)-   |   |
|   |         |   |   |                  |   |   |
|   |         |   |   |   ---topbar---   |   |   |
|   |         |   |   |   |          |   |   |   |
|   |         |   |   |   ---topbar---   |   |   |
|   |         |   |   |                  |   |   |
|   s         s   |   |   ---content--   |   |   |
|   i         i   |   |   |          |   |   |   |
|   d         d   m   |   |  HELLO,  |   |   m   |
|   e         e   a   |   |  WORLD!  |   |   a   |
|   b         b   i   |   |          |   |   i   |
|   a         a   n   |   ---content--   |   n   |
|   r         r   |   |                  |   |   |
|   |         |   |   ---------div---------  |   |
|   |         |   |                          |   |
|   |         |   |   -------footer--------  |   |
|   |         |   |   |                   |  |   |
|   |         |   |   -------footer--------  |   |
|   |         |   |                          |   |
|   --sidebar--   ------------main------------   |
|                                                |
------------------wrapper-------------------------
-->

<template>
  <div id="wrapper" class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 h-screen">
    <layout-sidebar class="col-span-1"/>
    <layout-main class="col-span-2 md:col-span-3 lg:col-span-4"/>
  </div>
</template>

<script>
import LayoutMain from './LayoutMain.vue'
import LayoutSidebar from './LayoutSidebar.vue'

export default {
  components: {
    LayoutMain,
    LayoutSidebar,
  },
  created() {
    this.$store.commit('SET_SITE_NAME', 'Unknown')
    this.$store.commit('SET_SITE_DESCRIPTION', 'description undefined')
  },
}
</script>
