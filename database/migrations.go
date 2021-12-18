package database

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var embedded embed.FS

type migration struct {
	order   int
	name    string
	content string
}

func applyMigrations(ctx context.Context, db *Database) error {
	files, err := embedded.ReadDir("migrations")
	if err != nil {
		return err
	}
	var migrations []migration
	for _, file := range files {
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 1 {
			return fmt.Errorf("invalid migration %s", file.Name())
		}
		order, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid migration %s: %w", file.Name(), err)
		}
		content, err := embedded.ReadFile(filepath.Join("migrations/", file.Name()))
		if err != nil {
			return err
		}
		migrations = append(migrations, migration{
			order:   order,
			name:    file.Name(),
			content: strings.TrimSpace(string(content)),
		})
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].order < migrations[j].order
	})

	tx, err := db.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Release()

	var current int

	if err := tx.Exec("PRAGMA user_version;", nil, func(stmt *Statement) bool {
		current = stmt.ColumnInt(0)
		return true
	}); err != nil {
		return err
	}

	for _, m := range migrations {
		if m.order <= current {
			continue
		}
		if err := tx.Exec(m.content, nil, nil); err != nil {
			return err
		}
		if err := tx.Exec("PRAGMA user_version = $current;", func(stmt *Statement) {
			stmt.SetInt64("$current", int64(m.order))
		}, nil); err != nil {
			return err
		}
	}
	return tx.Commit()
}
