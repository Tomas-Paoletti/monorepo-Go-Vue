import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'

// Wiring del frontend: monta App con Pinia (estado) y Router (navegación).
createApp(App).use(createPinia()).use(router).mount('#app')
