# 006 — Shell del frontend (Vue 3 + TS + Vite)

## Objetivo

Montar el esqueleto de la SPA: Vite + Vue 3 `<script setup lang="ts">` + Pinia + Vue Router +
Tailwind + axios, con proxy `/api` → backend Go. Todo en **TypeScript** (chequeo con `vue-tsc`).

## Contrato (estructura)

```
frontend/
├── vite.config.ts        # proxy /api → localhost:8080, alias @ → src
├── tsconfig.json         # extiende @vue/tsconfig (strict)
├── tailwind.config.ts postcss.config.js
├── index.html
└── src/
    ├── main.ts           # createApp + Pinia + Router + main.css
    ├── App.vue           # layout base (header + <RouterView/>)
    ├── router/index.ts   # rutas (lazy import por ruta)
    ├── shared/http.ts    # axios con baseURL '/api/v1'
    └── features/<dominio>/{api.ts, types.ts, stores/, components/, *View.vue}
```

## Criterios de aceptación

- [ ] `npm run build` (que corre `vue-tsc --noEmit && vite build`) pasa sin errores de tipos.
- [ ] `npm run dev` levanta Vite en :5173.
- [ ] Con el backend en :8080, una llamada a `/api/v1/...` desde el front se proxea y responde.
- [ ] El header navega entre `/movies` y `/lists` sin recargar la página (SPA).

## Notas de aprendizaje

- **Por qué el proxy y no CORS-por-todos-lados:** en dev el front pide a `/api` (mismo origen que
  Vite); Vite reenvía a Go. El browser nunca ve un origen cruzado para esas llamadas. El backend
  igual tiene CORS por si se pega directo a :8080. En fase 2 (embed) todo es el mismo origen.
- **`features/<dominio>/`:** cada dominio agrupa su `api.ts` (HTTP), `types.ts` (contrato con el
  back), `stores/` (Pinia) y `components/` + vistas. Espeja la convención del proyecto React.
- **Un solo cliente axios (`shared/http.ts`):** centraliza `baseURL` y headers. Los componentes
  no tocan axios directo; llaman a las funciones de `api.ts`. Facilita cambiar auth/errores global.
- **Pinia vs props/emits:** el estado que cruza vistas (listas, resultados de búsqueda) vive en
  stores Pinia; lo local de un componente queda en `ref`. Aprender dónde va cada cosa es el punto.
