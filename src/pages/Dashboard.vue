<template>
  <div>
    <div>Config:</div>
    <prism v-if="config" language="yaml" class="rounded !text-sm">{{ config }}</prism>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
  data() {
    return {
      config: '',
    }
  },
  created() {
    this.axios.get('/api/v1/config').then(r => {
      this.config = r.data
    }, r => {
      console.log(r)
    })
  },
  methods: {},
})
</script>
