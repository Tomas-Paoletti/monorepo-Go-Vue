<script setup lang="ts">
import type { Movie } from '../types'

// Card de una película. Si se pasa actionLabel, muestra un botón que emite
// "action" con la película. El padre decide qué significa (agregar, quitar…).
// El componente no conoce la API: solo presenta.
defineProps<{ movie: Movie; actionLabel?: string }>()
defineEmits<{ action: [movie: Movie] }>()
</script>

<template>
  <div class="rounded-lg border border-slate-800 bg-slate-900 p-4">
    <div class="flex items-start justify-between gap-2">
      <h3 class="font-semibold">
        {{ movie.title }}
        <span v-if="movie.year" class="text-slate-400">({{ movie.year }})</span>
      </h3>
      <button
        v-if="actionLabel"
        class="shrink-0 rounded bg-indigo-600 px-2 py-1 text-xs hover:bg-indigo-500"
        @click="$emit('action', movie)"
      >
        {{ actionLabel }}
      </button>
    </div>
    <p class="mt-1 flex flex-wrap gap-1">
      <span
        v-for="g in movie.genres"
        :key="g"
        class="rounded bg-slate-800 px-2 py-0.5 text-xs text-slate-300"
      >
        {{ g }}
      </span>
    </p>
    <p v-if="movie.rating > 0" class="mt-2 text-sm text-amber-400">
      ★ {{ movie.rating.toFixed(1) }}
      <span class="text-slate-500">({{ movie.ratingCount }})</span>
    </p>
  </div>
</template>
