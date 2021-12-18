package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	db, err := Open("file:memory:?mode=memory")
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, applyMigrations(ctx, db))
}
