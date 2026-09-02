// Espejo del movieResponse del backend (internal/api/dto.go). Los nombres van
// en camelCase, igual que el JSON que devuelve Go.
export interface Movie {
  id: number
  title: string
  altTitle?: string
  year?: number
  genres: string[]
  rating: number
  ratingCount: number
}

export interface SearchParams {
  q?: string
  genre?: string
  yearFrom?: number
  yearTo?: number
  minRating?: number
  limit?: number
  offset?: number
}
