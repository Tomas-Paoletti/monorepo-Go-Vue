import { defineStore } from 'pinia'
import { ref } from 'vue'
import { searchMovies } from '../api'
import type { Movie, SearchParams } from '../types'

// Store Pinia del buscador: guarda resultados, estado de carga y error.
export const useMoviesStore = defineStore('movies', () => {
  const results = ref<Movie[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function search(params: SearchParams) {
    loading.value = true
    error.value = null
    try {
      results.value = await searchMovies(params)
    } catch {
      error.value = 'No se pudo buscar. ¿Está levantado el backend?'
      results.value = []
    } finally {
      loading.value = false
    }
  }

  return { results, loading, error, search }
})
