<template>
  <footer class="">
    <div class="fcc px-3 py-6 text-center text-xs">
      <a class="mr-1" href="https://github.com/spongeprojects/kubebigbrother" target="_blank">SpongeProjects</a>
      <span class="mr-1">&copy; {{ year }}</span>
      <github-icon/>
    </div>
  </footer>
</template>

<script>
import GithubIcon from '../components/GithubIcon.vue'
import {DateTime} from 'luxon'

export default {
  components: {GithubIcon},
  data() {
    return {
      year: DateTime.local().year,
    }
  },
}
</script>

