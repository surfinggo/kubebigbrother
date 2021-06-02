<template>
  <footer class="">
    <div class="px-3 py-6 text-center text-xs">
      <span>Copyright <a href="http://kubebigbrother.spongeprojects.com" target="_blank">SpongeProjects</a> &copy; {{ year }}</span>
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

