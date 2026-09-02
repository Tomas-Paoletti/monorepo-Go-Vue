import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// El dev server de Vite corre en :5173 y PROXEA /api hacia el backend Go en
// :8080. Así el frontend pide siempre a "/api/..." (mismo origen para el
// browser) y en dev Vite reenvía. En fase 2 (go:embed) el mismo binario Go
// sirve el frontend y este proxy deja de usarse.
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
