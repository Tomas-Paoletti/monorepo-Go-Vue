import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '../api'
import type { List } from '../types'

// Store Pinia de listas. Mantiene el índice de listas y ofrece helpers usados
// desde varias vistas (crear, agregar película, asegurar una lista default).
export const useListsStore = defineStore('lists', () => {
  const lists = ref<List[]>([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      lists.value = await api.getLists()
    } finally {
      loading.value = false
    }
  }

  async function create(name: string, description = ''): Promise<List> {
    const list = await api.createList(name, description)
    lists.value.unshift(list)
    return list
  }

  async function remove(id: string) {
    await api.deleteList(id)
    lists.value = lists.value.filter((l) => l.id !== id)
  }

  async function addMovie(listId: string, movieId: number) {
    await api.addMovie(listId, movieId)
  }

  // ensureDefaultList devuelve una lista existente o crea una "Mi lista".
  // Simplifica el demo: agregar desde el buscador sin elegir lista.
  async function ensureDefaultList(): Promise<string> {
    if (lists.value.length === 0) {
      await fetchAll()
    }
    if (lists.value.length === 0) {
      const list = await create('Mi lista')
      return list.id
    }
    return lists.value[0].id
  }

  return { lists, loading, fetchAll, create, remove, addMovie, ensureDefaultList }
})
