package dummy

// We need an import of 99designs/gqlgen to trigger the analyzer
import (
	_ "github.com/99designs/gqlgen/graphql"
)

type Resolver struct {
	db *DB
}

// Missing dataloader, direct DB call
func (r *Resolver) Users() {
	r.db.Query("SELECT id FROM users")
}
