---
description: Crea un documento de spec en docs/specs/ desde una plantilla
argument-hint: <nombre> [--variant=dialectic|lean|bdd|contract|test-first|full]
allowed-tools: Bash(cd:*), Bash(npm:*), Bash(make:*), Read
---

Creá un nuevo spec ejecutando el comando de `ai-tools`.

Argumentos recibidos: `$ARGUMENTS`

Pasos:
1. Ejecutá desde `ai-tools/`: `npm run spec -- new $ARGUMENTS`
   (si el usuario no pasó `--variant`, se usa la variante `dialectic` por default).
2. Leé el archivo `docs/specs/<slug>.md` que se acaba de crear y mostráselo al usuario.
3. Recordale que lo llene y que después puede pedir **"validá el spec &lt;nombre&gt;"** para
   que la skill `spec-validate` contraste supuestos contra el código y pregunte lo que falte.

No edites el código del proyecto en este comando: solo generás y mostrás el spec.
