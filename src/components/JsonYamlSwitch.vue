<template>
  <div class="fcc text-sm">
    <Switch
        :model-value="modelValue"
        @update:model-value="$emit('update:modelValue', $event)"
        class="bg-blue-400 relative inline-flex flex-shrink-0 h-[20px] w-[36px] border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none"
    >
      <span
          :class="modelValue ? 'translate-x-4' : 'translate-x-0'"
          class="pointer-events-none inline-block h-[16px] w-[16px] rounded-full bg-white shadow-lg transform ring-0 transition ease-in-out duration-200"
      />
    </Switch>
    <span class="ml-1"
          :class="modelValue ? 'text-gray-500' : 'text-gray-900'">JSON</span>/<span
      :class="modelValue ? 'text-gray-900' : 'text-gray-500'">YAML</span>
  </div>
</template>

<script lang="ts">
import {ref} from 'vue'
import {Switch} from '@headlessui/vue'

export default {
  components: {Switch},
  props: ['modelValue'],
  emits: ['update:modelValue'],
}
</script>