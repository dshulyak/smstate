package main

import (
	"context"
	"log"

	"crawshaw.io/sqlite/sqlitex"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var tables = `

CREATE TABLE blocks ( 
	block_id CHAR(20) PRIMARY KEY,
	layer_id UNSIGNED MEDIUMINT,
	in_input_vector BOOL,
	contextually_valid BOOL,
	block    BLOB
) WITHOUT ROWID;

CREATE INDEX blocks_by_layer_id ON blocks(layer_id);
CREATE INDEX blocks_by_in_input_vector ON blocks(in_input_vector,layer_id) WHERE in_input_vector = 1;
CREATE INDEX blocks_by_contextual_validity ON blocks(contextually_valid,layer_id) WHERE contextually_valid = 1;

---

CREATE TABLE layers (
	layer_id UNSIGNED MEDIUMINT PRIMARY KEY,
	/* PROCESSED, SYNCED, STATE */
	label SMALLINT,
	hash CHAR(32),
	aggregated_hash CHAR(32)
) WITHOUT ROWID;

CREATE UNIQUE INDEX layer_id_by_label ON layers(label);

---

CREATE TABLE activations (
	atx_id CHAR(32) PRIMARY KEY,
	epoch_id UNSIGNED MEDIUMINT,
	node_id VARCHAR,
	timestamp UNSIGNED BIGINT, 
	is_top BOOL,
	header BLOB,
	body BLOB
) WITHOUT ROWID;

CREATE UNIQUE INDEX atx_by_epoch_node ON activations(epoch_id,node_id);
CREATE UNIQUE INDEX top_atx ON activations(is_top) WHERE is_top = 1;

---


CREATE TABLE poets (
	poet_id VARCHAR PRIMARY KEY,
	poet BLOB
) WITHOUT ROWID;

--- 

CREATE TABLE identieis (
	node_key VARCHAR PRIMARY KEY,
	vrf_key VARCHAR
) WITHOUT ROWID;
`

func main() {
	dbpool, err := sqlitex.Open("file:memory:?mode=memory", 0, 10)
	must(err)
	ctx := context.Background()
	conn := dbpool.Get(ctx)
	must(sqlitex.ExecScript(conn, tables))
}
