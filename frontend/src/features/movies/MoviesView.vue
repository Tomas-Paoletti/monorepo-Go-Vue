<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useMoviesStore } from './stores/movies'
import { useListsStore } from '@/features/lists/stores/lists'
import MovieCard from './components/MovieCard.vue'
import type { Movie } from './types'

// Vista buscador: input de texto + resultados. Permite agregar una película a
// la lista "activa" (la primera, para simplificar el demo de aprendizaje).
const store = useMoviesStore()
const { results, loading, error } = storeToRefs(store)

const listsStore = useListsStore()
const query = ref('')

function runSearch() {
  store.search({ q: query.value.trim(), limit: 30 })
}

async function addToList(movie: Movie) {
  const listId = await listsStore.ensureDefaultList()
  await listsStore.addMovie(listId, movie.id)
}

onMounted(runSearch)
</script>

<template>
  <section>
    <h1 class="mb-4 text-2xl font-bold">Buscar películas</h1>

    <form class="mb-6 flex gap-2" @submit.prevent="runSearch">
      <input
        v-model="query"
        type="text"
        placeholder="Título (ej: toy, matrix)…"
        class="flex-1 rounded border border-slate-700 bg-slate-900 px-3 py-2 outline-none focus:border-indigo-500"
      />
      <button class="rounded bg-indigo-600 px-4 py-2 hover:bg-indigo-500" type="submit">
        Buscar
      </button>
    </form>

    <p v-if="loading" class="text-slate-400">Cargando…</p>
    <p v-else-if="error" class="text-red-400">{{ error }}</p>
    <p v-else-if="results.length === 0" class="text-slate-400">Sin resultados.</p>

    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <MovieCard
        v-for="movie in results"
        :key="movie.id"
        :movie="movie"
        action-label="+ Lista"
        @action="addToList"
      />
    </div>
  </section>
</template>
