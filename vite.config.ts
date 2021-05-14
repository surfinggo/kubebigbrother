import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import WindiCSS from 'vite-plugin-windicss'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    WindiCSS(),
  ],
  server: {
    port: 1984,
    proxy: {
      '/api': {
        target: 'http://localhost:8984',
        changeOrigin: true
      },
    }, // end proxy
  }, // end server
})
