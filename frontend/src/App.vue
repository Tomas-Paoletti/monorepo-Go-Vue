<script setup lang="ts">
// Layout base: header con navegación + el <router-view> donde se pinta la vista
// activa. Todo el estilo con clases de Tailwind.
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100">
    <header class="border-b border-slate-800">
      <nav class="mx-auto flex max-w-5xl items-center gap-6 px-4 py-4">
        <span class="text-lg font-bold">🎬 Movie Roulette</span>
        <RouterLink
          to="/movies"
          class="text-sm text-slate-300 hover:text-white"
          active-class="text-white font-semibold"
        >
          Películas
        </RouterLink>
        <RouterLink
          to="/lists"
          class="text-sm text-slate-300 hover:text-white"
          active-class="text-white font-semibold"
        >
          Listas
        </RouterLink>
      </nav>
    </header>

    <main class="mx-auto max-w-5xl px-4 py-8">
      <RouterView />
    </main>
  </div>
</template>
