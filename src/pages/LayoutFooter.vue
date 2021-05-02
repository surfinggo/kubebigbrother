<template>
  <footer class="">
    <div class="p-5 bg-gray-200 text-center text-xs">
      <span>Copyright &copy; {{ year }}</span>
    </div>
  </footer>
</template>

<script>
import {DateTime} from 'luxon'

export default {
  data() {
    return {
      year: DateTime.local().year,
    }
  },
}
</script>

