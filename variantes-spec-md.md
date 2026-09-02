# Variantes de plantilla de spec (`docs/specs/<feature>.md`)

Seis variantes con filosofías distintas, para probar cuál te da mejor output con el agente.
Al final hay un protocolo de evaluación para no decidir "a ojo".

---

## Variante A — Mínima (lean)

**Hipótesis:** menos ruido = más señal. El agente completa los huecos con el código real en vez de con lo que vos imaginaste.
**Cuándo:** features chicas, codebase que el agente ya conoce bien, vos disponible para responder rápido.

```markdown
# <Feature>

## Qué
Una frase. Qué hace y para quién.

## Por qué
Una frase. Qué problema resuelve.

## Alcance
- Incluye: ...
- NO incluye: ...

## Criterios de aceptación
- [ ] ...
- [ ] ...

## Dudas abiertas
- ...
```

---

## Variante B — Comportamiento (Gherkin / Given-When-Then)

**Hipótesis:** describir comportamiento observable en vez de implementación reduce las alucinaciones de diseño y se traduce casi 1:1 a tests.
**Cuándo:** lógica de negocio con reglas y casos borde (validaciones, pricing, reglas de cancelación).

```markdown
# <Feature>

## Contexto
2-3 líneas.

## Escenarios

### Escenario 1: <nombre>
- **Dado** que ...
- **Cuando** ...
- **Entonces** ...

### Escenario 2: <caso borde>
- **Dado** que ...
- **Cuando** ...
- **Entonces** ...

### Escenario 3: <error>
- **Dado** que ...
- **Cuando** ...
- **Entonces** ... (mensaje/código exacto: `...`)

## Fuera de alcance
- ...

## Dudas abiertas
- ...
```

---

## Variante C — Contract-first

**Hipótesis:** si fijás las interfaces (endpoints, tipos, esquema de DB) antes, el agente no inventa firmas ni renombra cosas a mitad de camino.
**Cuándo:** APIs, integraciones con terceros, trabajo cross-package en monorepo donde otro paquete consume lo que hacés.

```markdown
# <Feature>

## Objetivo
1-2 frases.

## Contrato

### Endpoint
```
POST /api/v1/<recurso>
```

### Request
```ts
{ ... }
```

### Response 200
```ts
{ ... }
```

### Errores
| Código | Cuándo | Body |
|---|---|---|
| 422 | ... | ... |
| 409 | ... | ... |

### Cambios de schema
```sql
-- migración
```

## Invariantes
- ...

## Dudas abiertas
- ...
```

---

## Variante D — Test-first

**Hipótesis:** si la spec ES la lista de tests, el agente tiene un criterio de "terminado" verificable y se auto-corrige sin que vos revises.
**Cuándo:** refactors, bugfixes, cualquier cosa con buena cobertura existente.

```markdown
# <Feature>

## Qué
Una frase.

## Tests que deben pasar

### Unit
- [ ] `<Clase>::<metodo>` devuelve X cuando Y
- [ ] `<Clase>::<metodo>` tira `<Excepcion>` cuando Z

### Integración
- [ ] `POST /...` con payload válido → 201 y fila en `<tabla>`
- [ ] `POST /...` sin auth → 401

### Regresión (no romper)
- [ ] `<test existente>` sigue verde

## Comando de verificación
```bash
<pest/phpunit/vitest/dotnet test ...>
```

## Dudas abiertas
- ...
```

---

## Variante E — Completa (ceremonial)

**Hipótesis:** captura el *por qué* y las alternativas descartadas, así el agente no re-propone lo que ya rechazaste.
**Cuándo:** features grandes, decisiones arquitectónicas, cosas que otro (o vos en 6 meses) va a tener que entender.

```markdown
# <Feature>

**Estado:** borrador | aprobada | implementada
**Ticket:** CU-xxxx
**Fecha:** YYYY-MM-DD

## 1. Contexto
Qué existe hoy, qué duele.

## 2. Objetivo
Qué queremos lograr. Métrica de éxito si aplica.

## 3. No-objetivos
Explícitamente fuera de alcance.

## 4. Diseño propuesto
Archivos/módulos que se tocan. Flujo de datos.

## 5. Alternativas descartadas
| Opción | Por qué no |
|---|---|
| ... | ... |

## 6. Casos borde
- ...

## 7. Criterios de aceptación
- [ ] ...

## 8. Riesgos
- ...

## 9. Decisiones cerradas
- [YYYY-MM-DD] Se decidió X porque Y.

## 10. Dudas abiertas
- [ ] ...
```

---

## Variante F — Dialéctica (huecos explícitos)

**Hipótesis:** en vez de intentar cubrirlo todo, la spec declara qué NO sabe. El agente ataca esos huecos primero y pregunta antes de codear.
**Cuándo:** features donde vos mismo no tenés todo definido. Es la que mejor calza con tu flujo "pregunta si hay dudas".

```markdown
# <Feature>

## Lo que sé
- ...
- ...

## Lo que NO sé (resolver antes de codear)
| # | Pregunta | Quién decide | Estado |
|---|---|---|---|
| 1 | ... | yo / PM / negocio | abierta |
| 2 | ... | ... | abierta |

## Supuestos que estoy haciendo
> Si alguno es falso, PARÁ y avisá.
- Asumo que ...
- Asumo que ...

## Restricciones duras
- No tocar `<archivo/módulo>`
- Debe funcionar con `<versión/lib>`

## Definición de terminado
- [ ] ...

## Protocolo
1. Leé el código relevante y contrastá los supuestos.
2. Respondé lo que puedas de "lo que NO sé" con evidencia del código.
3. Lo que no puedas, preguntámelo (una pregunta por vez).
4. Actualizá esta spec con las respuestas.
5. Recién ahí proponé plan. No escribas código hasta que apruebe.
```

---

## Cómo elegir sin engañarte

No compares leyendo. Corré el mismo trabajo real con cada variante y medí.

### Protocolo

1. Elegí **3 features ya implementadas** del repo (una chica, una mediana, una con lógica de negocio fea). Ya sabés cómo terminaron: ese es tu ground truth.
2. Escribí la spec de cada una en 2-3 variantes distintas.
3. Corré cada spec en una **sesión limpia** (`/clear`), mismo modelo, mismo prompt de arranque: *"Leé `docs/specs/X.md` y seguí el proceso."*
4. Anotá los resultados en la tabla.

### Qué medir

| Métrica | Cómo | Qué buscás |
|---|---|---|
| Preguntas útiles | ¿Cuántas hizo? ¿Cuántas eran las correctas? | Alto en las correctas, bajo en ruido |
| Asunciones silenciosas | Cosas que decidió sin avisar | Lo más bajo posible |
| Iteraciones hasta verde | Cuántas vueltas para que pasen los tests | Menos = mejor |
| Fidelidad al codebase | ¿Respetó convenciones, capas, naming? | Sí/parcial/no |
| Scope creep | ¿Tocó cosas que no debía? | Ninguna |
| Costo de escribir la spec | Minutos que te llevó a vos | El punto es que ahorre, no que sume |

### Planilla

```markdown
| Feature | Variante | Preguntas útiles | Asunciones silenciosas | Iteraciones | Scope creep | Min. de spec | Nota /10 |
|---|---|---|---|---|---|---|---|
| Import PNR | A | | | | | | |
| Import PNR | F | | | | | | |
| ...     | | | | | | | |
```

### Sospechas previas (para falsar, no para creer)

- **A** va a rendir bien en cosas chicas y mal apenas hay reglas de negocio.
- **B** debería bajar mucho las asunciones silenciosas en lógica de negocio.
- **C** debería eliminar el problema de firmas inventadas en trabajo cross-package.
- **D** debería dar las menos iteraciones, pero solo si el repo ya tiene tests decentes.
- **E** probablemente sea overkill: mucho costo de escritura para poca ganancia salvo en features grandes.
- **F** debería ganar en "preguntas útiles", que es justo lo que buscás.

Lo más probable es que termines con **dos** plantillas, no una: una liviana (A o F) para el día a día y una pesada (C o E) para lo grande. Eso está bien.

---

## Reglas de soporte en `CLAUDE.md`

Sea cual sea la variante que gane, esto va aparte, en el `CLAUDE.md` raíz:

```markdown
## Specs
- La spec en `docs/specs/` es la fuente de verdad del requerimiento.
- Si el código contradice la spec, PARÁ y preguntá cuál está bien. No resuelvas solo.
- Si algo de la spec es ambiguo o tiene 2+ interpretaciones, mostrámelas. No elijas.
- Antes de codear: confirmá el plan y esperá aprobación.
- Al terminar: actualizá la sección "Dudas abiertas" y "Decisiones cerradas" de la spec.
```
