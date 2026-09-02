---
name: spec-validate
description: Valida un spec de docs/specs/ contra el código real, resuelve supuestos con evidencia, pregunta lo que falta de a una pregunta por vez y actualiza el spec. Usar cuando el usuario diga "validá el spec X", "revisá la spec X" o antes de planificar una feature que tiene spec.
---

# Validación de spec (spec-driven development)

Objetivo: llevar un spec de "borrador" a "listo para planificar", contrastándolo con el código
real y cerrando toda ambigüedad **antes** de escribir una línea de código.

## Cuándo se activa

Cuando el usuario pide validar/revisar un spec (ej. "validá el spec checkout"), o cuando va a
planificar una feature que ya tiene un archivo en `docs/specs/`.

## Protocolo (seguir en orden)

1. **Leé el spec** en `docs/specs/<nombre>.md`. Si no existe, listá `docs/specs/` y pedí el
   nombre correcto. Identificá su variante mirando las secciones (dialectic, lean, bdd, etc.).

2. **Leé el código relevante** del repo antes de opinar. Buscá los módulos, handlers, modelos,
   migraciones y tests que la feature va a tocar. No asumas: verificá firmas, nombres y
   convenciones existentes.

3. **Contrastá cada supuesto y cada "Lo que NO sé"** contra la evidencia del código:
   - Lo que puedas responder con el código, respondelo vos y anotá la evidencia
     (`archivo:línea`).
   - Marcá los supuestos que el código **contradice**: eso es un hallazgo, avisalo explícito.

4. **Preguntá lo que no se pueda resolver con el código — UNA pregunta por vez.** No dispares una
   lista. Esperá la respuesta, actualizá el spec, y recién ahí hacé la siguiente pregunta.
   Si algo tiene 2+ interpretaciones razonables, mostrá las opciones; no elijas por el usuario.

5. **Actualizá el spec** a medida que se resuelve:
   - Mové ítems de "Dudas abiertas" / "Lo que NO sé" → "Decisiones cerradas" con fecha
     (`[YYYY-MM-DD] Se decidió X porque Y`).
   - Ajustá supuestos que resultaron falsos.
   - Cuando no queden huecos, cambiá `**Estado:**` a `aprobada`.

6. **No propongas plan ni escribas código** hasta que el usuario apruebe explícitamente. El
   entregable de esta skill es un spec cerrado, no una implementación.

## Reglas duras

- Si el código contradice el spec, PARÁ y preguntá cuál está bien. No resuelvas solo.
- Si algo del spec es ambiguo, mostrá las interpretaciones. No elijas.
- Una pregunta por vez. Siempre.
- El spec en `docs/specs/` es la fuente de verdad del requerimiento: mantenelo actualizado.
