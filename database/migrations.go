package database

import (
	"bufio"
	"bytes"
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
	content *bufio.Scanner
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
		scanner := bufio.NewScanner(bytes.NewBuffer(content))
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if i := bytes.Index(data, []byte(";")); i >= 0 {
				return i + 1, data[0 : i+1], nil
			}
			return 0, nil, nil
		})
		migrations = append(migrations, migration{
			order:   order,
			name:    file.Name(),
			content: scanner,
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
		for m.content.Scan() {
			if err := tx.Exec(m.content.Text(), nil, nil); err != nil {
				return err
			}
		}
		// binding values in pragma statement is not allowed
		if err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d;", m.order), nil, nil); err != nil {
			return err
		}
	}
	return tx.Commit()
}
