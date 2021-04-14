import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    server: {
        port: 1984,
        proxy: {
            '/api': {
                target: 'http://localhost:1949',
                changeOrigin: true
            },
        }, // end proxy
    }, // end server
})
