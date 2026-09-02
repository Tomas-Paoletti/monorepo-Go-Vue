<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { RouterLink } from 'vue-router'
import { useListsStore } from './stores/lists'

// Vista índice de listas: crear una nueva y navegar a cada una.
const store = useListsStore()
const { lists, loading } = storeToRefs(store)
const newName = ref('')

async function submit() {
  const name = newName.value.trim()
  if (!name) return
  await store.create(name)
  newName.value = ''
}

onMounted(store.fetchAll)
</script>

<template>
  <section>
    <h1 class="mb-4 text-2xl font-bold">Mis listas</h1>

    <form class="mb-6 flex gap-2" @submit.prevent="submit">
      <input
        v-model="newName"
        type="text"
        placeholder="Nombre de la lista…"
        class="flex-1 rounded border border-slate-700 bg-slate-900 px-3 py-2 outline-none focus:border-indigo-500"
      />
      <button class="rounded bg-indigo-600 px-4 py-2 hover:bg-indigo-500" type="submit">
        Crear
      </button>
    </form>

    <p v-if="loading" class="text-slate-400">Cargando…</p>
    <p v-else-if="lists.length === 0" class="text-slate-400">
      Todavía no tenés listas. Creá una arriba.
    </p>

    <ul v-else class="space-y-2">
      <li v-for="list in lists" :key="list.id">
        <div
          class="flex items-center justify-between rounded-lg border border-slate-800 bg-slate-900 px-4 py-3"
        >
          <RouterLink
            :to="`/lists/${list.id}`"
            class="font-semibold hover:text-indigo-400"
          >
            {{ list.name }}
          </RouterLink>
          <button
            class="text-sm text-red-400 hover:text-red-300"
            @click="store.remove(list.id)"
          >
            Borrar
          </button>
        </div>
      </li>
    </ul>
  </section>
</template>
