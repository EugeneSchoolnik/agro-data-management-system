import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      jspdf: fileURLToPath(new URL('./node_modules/jspdf/dist/jspdf.es.min.js', import.meta.url)),
    },
  },
})
