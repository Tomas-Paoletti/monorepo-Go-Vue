# CLAUDE.md

Guía para agentes trabajando en este monorepo.

## Specs

- La spec en `docs/specs/` es la fuente de verdad del requerimiento.
- Si el código contradice la spec, PARÁ y preguntá cuál está bien. No resuelvas solo.
- Si algo de la spec es ambiguo o tiene 2+ interpretaciones, mostrámelas. No elijas.
- Antes de codear: confirmá el plan y esperá aprobación.
- Al terminar: actualizá las secciones "Dudas abiertas" y "Decisiones cerradas" de la spec.
- Para crear specs: `make spec ARGS="new <nombre>"` (o `/spec-new`). Para validarlas, la skill
  `spec-validate`. Ver `ai-tools/README.md`.
