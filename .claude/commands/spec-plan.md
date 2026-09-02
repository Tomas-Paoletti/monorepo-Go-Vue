---
description: Pasa un spec validado a plan de ejecución
argument-hint: <nombre>
allowed-tools: Read, Glob, Grep
---

Tomá el spec `docs/specs/$ARGUMENTS.md` como fuente de verdad del requerimiento y preparate
para planificar su ejecución.

Pasos:
1. Leé `docs/specs/$ARGUMENTS.md`. Si no existe, listá `docs/specs/` y pedí el nombre correcto.
2. Verificá que **no queden "Dudas abiertas" ni supuestos sin confirmar**. Si quedan, PARÁ:
   avisale al usuario que primero conviene correr la skill `spec-validate` (o resolvé esos huecos
   preguntando de a una pregunta por vez). No planifiques sobre un spec incompleto.
3. Con el spec cerrado, entrá en plan mode y diseñá el plan de implementación basándote
   únicamente en lo que dice el spec + el código real del repo. No inventes alcance que no esté
   en el spec.
