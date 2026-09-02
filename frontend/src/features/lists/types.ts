import type { Movie } from '@/features/movies/types'

// Espejo de listResponse del backend.
export interface List {
  id: string
  name: string
  description?: string
  createdAt: string
}

// Respuesta de GET /lists/:id
export interface ListDetail {
  list: List
  movies: Movie[]
}

// Filtros opcionales del giro (spinRequest en el backend).
export interface SpinFilters {
  genre?: string
  yearFrom?: number
  yearTo?: number
  minRating?: number
  excludeIds?: number[]
}
