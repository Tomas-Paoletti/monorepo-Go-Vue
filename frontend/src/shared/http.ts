import axios from 'axios'

// Cliente HTTP único de la app. baseURL '/api/v1' => en dev Vite proxea a
// localhost:8080; en prod (go:embed) es el mismo origen. Los componentes nunca
// hablan con axios directo: pasan por las funciones api.ts de cada feature.
export const http = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})
