# ai-tools

Herramientas de IA del monorepo. Hoy: **scaffolding de specs** para
*spec-driven development*.

## Instalar

```bash
cd ai-tools
npm install
```

## Uso

Desde `ai-tools/`:

```bash
npm run spec -- new checkout                 # crea docs/specs/NNN-checkout.md (variante project)
npm run spec -- new pricing --variant=bdd    # otra variante
npm run spec -- new checkout --force         # sobrescribe si ya existe
npm run spec -- list                         # lista los specs y su estado
```

O desde la raíz del repo con el Makefile:

```bash
make spec ARGS="new checkout"
make spec ARGS="new pricing --variant=bdd"
make spec ARGS="list"
```

Los specs se generan en `docs/specs/` (la raíz del repo se detecta subiendo hasta el `.git`).

## Variantes de plantilla

| Variante     | Filosofía                                    |
|--------------|----------------------------------------------|
| `project`    | **Default.** La forma real del repo: objetivo, contrato, modelo de datos, criterios, plan de test, notas. |
| `bdd`        | Given-When-Then para lógica de negocio.      |
| `contract`   | Contract-first: endpoints, tipos, schema.    |
| `test-first` | La spec ES la lista de tests.                |
| `full`       | Ceremonial: contexto, alternativas, riesgos. |

Las plantillas viven en `templates/spec/*.md` y son la fuente de verdad. Editalas ahí.

Los huecos se marcan **inline** con `[NEEDS CLARIFICATION: <pregunta>]` en el lugar exacto;
`spec-validate` los resuelve antes de codear (no hay sección fija de "dudas").
`spec new` numera solo: mira los `NNN-` existentes y usa el siguiente.

## Flujo completo (con Claude)

1. `/spec-new <nombre>` — crea el spec.
2. Llenás el spec.
3. "validá el spec `<nombre>`" — la skill `spec-validate` lee el código, contrasta supuestos,
   pregunta lo que falta y actualiza el spec. No codea hasta que apruebes.
4. `/spec-plan <nombre>` — pasa a plan de ejecución con el spec como fuente de verdad.

## Copiar a otro repo Go/C#/TS

Copiá la carpeta `ai-tools/` completa y `.claude/`. El motor es Node/TS puro (sin deps de
runtime), así que corre en cualquier repo con Node instalado, sin importar el lenguaje del backend.
