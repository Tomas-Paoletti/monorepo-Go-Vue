import { http } from '@/shared/http'
import type { Movie, SearchParams } from './types'

// Llamadas HTTP de la feature movies. Los componentes usan estas funciones,
// no axios directo.

export async function searchMovies(params: SearchParams): Promise<Movie[]> {
  const { data } = await http.get<{ movies: Movie[] }>('/movies', { params })
  return data.movies
}
