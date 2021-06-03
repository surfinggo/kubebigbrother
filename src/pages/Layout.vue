<!--
------------------wrapper-------------------------
|   --------------------topbar-----------------  |
|   |                                         |  |
|   --------------------topbar-----------------  |
|                                                |
|   --sidebar--   ------------main------------   |
|   |         |   |                          |   |
|   |         |   |   -div (flex-grow: 1)-   |   |
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
  <div id="wrapper" class="text-gray-600">
    <layout-topbar :config="config" :config-yaml="configYaml"/>
    <div class="w-full">
      <layout-sidebar :config="config"/>
      <div class="min-h-screen flex flex-col pl-60 xl:pl-72 pt-16">
        <div class="flex-grow p-3">
          <router-view/>
        </div>
        <layout-footer/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import LayoutFooter from './LayoutFooter.vue'
import LayoutSidebar from './LayoutSidebar.vue'
import LayoutTopbar from './LayoutTopbar.vue'

export default {
  components: {
    LayoutFooter,
    LayoutSidebar,
    LayoutTopbar,
  },
  created() {
    this.$store.commit('SET_SITE_NAME', 'Unknown')
    this.$store.commit('SET_SITE_DESCRIPTION', 'description undefined')

    this.axios.get('/api/v1/config').then(r => {
      this.config = r.data.json
      this.configYaml = atob(r.data.yaml)
    })
  },
  data() {
    return {
      config: '',
      configYaml: '',
    }
  },
}
</script>
