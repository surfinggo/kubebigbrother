<template>
  <div class="flex flex-col" id="main">
    <layout-topbar/>
    <div class="flex-grow">
      <router-view class="p-4" id="content"/>
    </div>
    <layout-footer/>
  </div>
</template>

<script>
import LayoutFooter from './LayoutFooter.vue'
import LayoutTopbar from './LayoutTopbar.vue'

export default {
  components: {
    LayoutFooter,
    LayoutTopbar,
  }
}
</script>
