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
	db, err := Open("file:test.sql")
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, applyMigrations(ctx, db))

	require.NoError(t, err)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]byte, 5)
	rng.Read(id)
	dst := make([]byte, 10)
	hex.Encode(dst, id)
	require.NoError(t, db.Exec("INSERT INTO blocks (block_id) VALUES (?1);", func(stmt *Statement) {
		stmt.BindBytes(1, dst)
	}, nil))
}
