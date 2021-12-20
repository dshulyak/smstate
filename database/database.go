package database

import (
	"context"
	"errors"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

var (
	ErrNoConnection = errors.New("database: no free connection")
)

type Executor interface {
	Exec(string, Encoder, Decoder) error
}

type Database struct {
	pool *sqlitex.Pool
}

type Statement = sqlite.Stmt

type Encoder func(*Statement)

type Decoder func(*Statement) bool

func defaultConf() *conf {
	return &conf{
		connections: 16,
	}
}

type conf struct {
	flags       sqlite.OpenFlags
	connections int
}

func WithConnections(n int) Opt {
	return func(c *conf) {
		c.connections = n
	}
}

type Opt func(c *conf)

func InMemory() *Database {
	db, err := Open("file::memory:?mode=memory", WithConnections(1))
	if err != nil {
		panic(err)
	}
	return db
}

func Open(uri string, opts ...Opt) (*Database, error) {
	config := defaultConf()
	for _, opt := range opts {
		opt(config)
	}
	pool, err := sqlitex.Open(uri, config.flags, config.connections)
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool}, nil
}

func (db *Database) Tx(ctx context.Context) (*Tx, error) {
	conn := db.pool.Get(ctx)
	if conn == nil {
		return nil, ErrNoConnection
	}
	tx := &Tx{db: db, conn: conn}
	return tx, tx.begin()
}

func (db *Database) ConcurrentTx(ctx context.Context) (*Tx, error) {
	conn := db.pool.Get(ctx)
	if conn == nil {
		return nil, ErrNoConnection
	}
	tx := &Tx{db: db, conn: conn}
	return tx, tx.beginConcurrent()
}

func (db *Database) Exec(query string, encoder Encoder, decoder Decoder) error {
	conn := db.pool.Get(context.Background())
	if conn == nil {
		return ErrNoConnection
	}
	defer db.pool.Put(conn)
	return exec(conn, query, encoder, decoder)
}

func exec(conn *sqlite.Conn, query string, encoder Encoder, decoder Decoder) error {
	stmt, err := conn.Prepare(query)
	if err != nil {
		return err
	}
	if encoder != nil {
		encoder(stmt)
	}
	defer stmt.ClearBindings()

	for {
		row, err := stmt.Step()
		if !row || err != nil {
			return err
		}
		if decoder != nil && !decoder(stmt) {
			return nil
		}
	}
}

type Tx struct {
	db       *Database
	conn     *sqlite.Conn
	commited bool
	err      error
}

func (tx *Tx) begin() error {
	stmt := tx.conn.Prep("BEGIN;")
	_, err := stmt.Step()
	return err
}

// by default sqlite will support single writer, with concurrent it will optimistically support many
// but the code should expect that transaction may be aborted if another concurrent transaction
// modified same rows.
//
// https://sqlite.org/src/doc/begin-concurrent/doc/begin_concurrent.md
func (tx *Tx) beginConcurrent() error {
	stmt := tx.conn.Prep("BEGIN CONCURRENT;")
	_, err := stmt.Step()
	return err
}

func (tx *Tx) Commit() error {
	stmt := tx.conn.Prep("COMMIT;")
	_, tx.err = stmt.Step()
	if tx.err != nil {
		return tx.err
	}
	tx.commited = true
	return nil
}

func (tx *Tx) Release() error {
	defer tx.db.pool.Put(tx.conn)
	if tx.commited {
		return nil
	}
	stmt := tx.conn.Prep("ROLLBACK")
	_, tx.err = stmt.Step()
	return tx.err
}

func (tx *Tx) Exec(query string, encoder Encoder, decoder Decoder) error {
	return exec(tx.conn, query, encoder, decoder)
}
