<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import * as api from './api'
import type { ListDetail } from './types'
import type { Movie } from '@/features/movies/types'
import MovieCard from '@/features/movies/components/MovieCard.vue'

// Vista detalle: muestra las películas de la lista y la RULETA (girar).
const route = useRoute()
const id = route.params.id as string

const detail = ref<ListDetail | null>(null)
const loading = ref(true)
const winner = ref<Movie | null>(null)
const spinError = ref<string | null>(null)
const spinning = ref(false)

async function load() {
  loading.value = true
  detail.value = await api.getList(id)
  loading.value = false
}

async function removeMovie(movie: Movie) {
  await api.removeMovie(id, movie.id)
  await load()
}

async function spin() {
  spinning.value = true
  spinError.value = null
  winner.value = null
  try {
    winner.value = await api.spin(id)
  } catch {
    // 422 => ninguna película pasa el filtro / lista vacía.
    spinError.value = 'No hay película para elegir (¿lista vacía?).'
  } finally {
    spinning.value = false
  }
}

onMounted(load)
</script>

<template>
  <section v-if="detail">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">{{ detail.list.name }}</h1>
      <button
        class="rounded bg-emerald-600 px-4 py-2 font-semibold hover:bg-emerald-500 disabled:opacity-50"
        :disabled="spinning || detail.movies.length === 0"
        @click="spin"
      >
        🎲 {{ spinning ? 'Girando…' : 'Girar ruleta' }}
      </button>
    </div>

    <div
      v-if="winner"
      class="mb-6 rounded-lg border border-emerald-700 bg-emerald-950 p-4"
    >
      <p class="text-sm text-emerald-400">La ruleta eligió:</p>
      <p class="text-lg font-bold">
        {{ winner.title }}
        <span v-if="winner.year" class="text-slate-400">({{ winner.year }})</span>
      </p>
    </div>
    <p v-if="spinError" class="mb-6 text-amber-400">{{ spinError }}</p>

    <p v-if="detail.movies.length === 0" class="text-slate-400">
      Lista vacía. Agregá películas desde el buscador.
    </p>
    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <MovieCard
        v-for="movie in detail.movies"
        :key="movie.id"
        :movie="movie"
        action-label="Quitar"
        @action="removeMovie"
      />
    </div>
  </section>
  <p v-else class="text-slate-400">Cargando…</p>
</template>
