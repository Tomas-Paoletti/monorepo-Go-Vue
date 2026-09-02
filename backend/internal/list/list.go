// Package list es el dominio de las listas de películas. Es data plana (sin
// comportamiento por ahora), pero vive acá para que storage y api la compartan
// sin que ninguno dependa del otro.
package list

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound se devuelve cuando la lista pedida no existe. Vive en el dominio
// para que la capa HTTP la detecte sin importar el paquete de storage.
var ErrNotFound = errors.New("list not found")

type List struct {
	ID          uuid.UUID
	Name        string
	Description string // "" si no tiene
	CreatedAt   time.Time
}
