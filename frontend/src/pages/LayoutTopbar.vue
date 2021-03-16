<template>
  <div class="p-4 bg-purple-100" id="topbar">
    <div class="text-right text-sm">
      <span v-if="serverVersion">Server version: {{ serverVersion }}</span>
      <span v-else><m-loading/></span>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      serverVersion: '',
    }
  },
  created() {
    this.axios.get('/api/v1/healthz').then(r => {
      this.serverVersion = r.data.version
    })
  },
}
</script>

