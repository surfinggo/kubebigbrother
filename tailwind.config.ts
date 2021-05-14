import {defineConfig} from 'windicss/helpers'

export default defineConfig({
  extract: {
    include: [
      'src/**/*.vue',
      'src/**/*.css',
      'src/**/*.scss',
    ],
  },
})