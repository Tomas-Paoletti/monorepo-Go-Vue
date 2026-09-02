import { createRouter, createWebHistory } from 'vue-router'

// Rutas de la SPA. Cada vista vive en su feature. Lazy import (() => import)
// para que Vite haga code-splitting por ruta.
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/movies' },
    {
      path: '/movies',
      name: 'movies',
      component: () => import('@/features/movies/MoviesView.vue'),
    },
    {
      path: '/lists',
      name: 'lists',
      component: () => import('@/features/lists/ListsView.vue'),
    },
    {
      path: '/lists/:id',
      name: 'list-detail',
      component: () => import('@/features/lists/ListDetailView.vue'),
    },
  ],
})

export default router
