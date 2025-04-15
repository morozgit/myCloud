import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8000', // туда, где у тебя работает FastAPI
        changeOrigin: true,
        rewrite: path => path.replace(/^\/api/, '/api'), // можно и просто path => path
      },
    },
  },
})
