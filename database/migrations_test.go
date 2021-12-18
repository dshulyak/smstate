package database

import (
	"context"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	// schema changes for some reason is not stored with other connections in memory mode.
	db := InMemory()

	ctx := context.Background()
	require.NoError(t, Apply(ctx, db))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]byte, 10)
	rng.Read(id)
	dst := make([]byte, 20)
	hex.Encode(dst, id)

	require.NoError(t, db.Exec("INSERT INTO blocks (id) VALUES (?1);", func(stmt *Statement) {
		stmt.BindBytes(1, dst)
	}, nil))
}
