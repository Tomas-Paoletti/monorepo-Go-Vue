import { http } from '@/shared/http'
import type { Movie } from '@/features/movies/types'
import type { List, ListDetail, SpinFilters } from './types'

// Llamadas HTTP de la feature lists.

export async function getLists(): Promise<List[]> {
  const { data } = await http.get<{ lists: List[] }>('/lists')
  return data.lists
}

export async function createList(name: string, description = ''): Promise<List> {
  const { data } = await http.post<List>('/lists', { name, description })
  return data
}

export async function getList(id: string): Promise<ListDetail> {
  const { data } = await http.get<ListDetail>(`/lists/${id}`)
  return data
}

export async function deleteList(id: string): Promise<void> {
  await http.delete(`/lists/${id}`)
}

export async function addMovie(listId: string, movieId: number): Promise<void> {
  await http.post(`/lists/${listId}/movies`, { movieId })
}

export async function removeMovie(listId: string, movieId: number): Promise<void> {
  await http.delete(`/lists/${listId}/movies/${movieId}`)
}

export async function spin(listId: string, filters: SpinFilters = {}): Promise<Movie> {
  const { data } = await http.post<{ winner: Movie }>(`/lists/${listId}/spin`, filters)
  return data.winner
}
