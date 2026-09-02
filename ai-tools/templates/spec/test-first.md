# <Feature>

**Estado:** borrador

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
make test
```

## Dudas abiertas
- ...
