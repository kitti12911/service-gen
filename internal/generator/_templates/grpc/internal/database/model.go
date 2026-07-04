package database

import (
	"time"

	"github.com/uptrace/bun"
)

// Example is a placeholder model. Replace with your real bun models and
// register them through models() so that lib-orm picks them up at startup
// and mapgen fields can walk the package.
type Example struct {
	bun.BaseModel `bun:"table:examples,alias:e"`

	ID        string    `bun:"id,pk"`
	Name      string    `bun:"name"`
	CreatedAt time.Time `bun:"created_at"`
}

func models() []any {
	return []any{
		(*Example)(nil),
	}
}
