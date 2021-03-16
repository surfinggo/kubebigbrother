<template>
  <button class="inline-block rounded px-2 py-1 text-center text-white"
          :class="disabled ? 'bg-blue-300 hover:bg-blue-300 cursor-not-allowed'
       : 'bg-blue-500 hover:bg-blue-600 cursor-pointer'"
          @click="handleClick">
    <slot/>
  </button>
</template>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
  props: {
    disabled: {
      type: Boolean,
      default: false,
    }
  },
  emits: ['click'],
  methods: {
    handleClick(e: MouseEvent) {
      if (this.disabled) {
        return false
      }
      this.$emit('click', e)
    }
  }
})
</script>