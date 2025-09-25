// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  ssr: false,
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  vite: { plugins: [tailwindcss(),], },
  app: {
    head: {
      script: [
        {
          src: 'https://cdn.jsdelivr.net/npm/chart.js',
          defer: true
        }
      ]
    }
  }
})
