package blocks

import "gihtub.com/dshulyak/sqliteexp/database"

type ID [20]byte

type Block struct {
	ID    ID
	Layer uint32
	Data  []byte
}

func Add(db database.Executor, block *Block) error {
	return db.Exec("insert or ignore into blocks (id, layer, block) values (?1, ?2, ?3);",
		func(stmt *database.Statement) {
			stmt.BindBytes(1, block.ID[:])
			stmt.BindInt64(2, int64(block.Layer))
			stmt.BindBytes(3, block.Data) // this is actually should encode block
		}, nil)
}

func AddHareOutput(db database.Executor, output []ID) error {
	for _, id := range output {
		if err := db.Exec("update blocks set hare_output = 1 where id = ?1;", func(stmt *database.Statement) {
			stmt.BindBytes(1, id[:])
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

func AddVerified(db database.Executor, verified []ID) error {
	for _, id := range verified {
		if err := db.Exec("update blocks set verified = 1 where id = ?1;", func(stmt *database.Statement) {
			stmt.BindBytes(1, id[:])
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

func Get(db database.Executor, id ID) (*Block, error) {
	var rst Block
	if err := db.Exec("select (layer, block) from blocks where id = ?1;", func(stmt *database.Statement) {
		stmt.BindBytes(1, id[:])
	}, func(stmt *database.Statement) bool {
		rst.ID = id
		rst.Layer = uint32(stmt.ColumnInt64(1))
		lth := stmt.ColumnLen(2)
		rst.Data = make([]byte, lth)
		stmt.ColumnBytes(2, rst.Data)
		return true
	}); err != nil {
		return nil, err
	}
	return &rst, nil
}

func IsVerified(db database.Executor, id ID) (bool, error) {
	var rst bool
	if err := db.Exec("select verified from blocks where id = ?1;", func(stmt *database.Statement) {
		stmt.BindBytes(1, id[:])
	}, func(stmt *database.Statement) bool {
		rst = stmt.ColumnInt(1) != 0
		return true
	}); err != nil {
		return false, err
	}
	return rst, nil
}

func FromLayer(db database.Executor, lid uint32) ([]ID, error) {
	var rst []ID
	if err := db.Exec("select id from blocks where layer = ?1;", func(stmt *database.Statement) {
		stmt.BindInt64(1, int64(lid))
	}, func(stmt *database.Statement) bool {
		id := ID{}
		stmt.ColumnBytes(1, id[:])
		rst = append(rst, id)
		return true
	}); err != nil {
		return nil, err
	}
	return rst, nil
}

func FromLayerByStatus(db database.Executor, lid uint32) (map[ID]bool, error) {
	var rst = map[ID]bool{}
	if err := db.Exec("select (id, verified) from blocks where layer = ?1;", func(stmt *database.Statement) {
		stmt.BindInt64(1, int64(lid))
	}, func(stmt *database.Statement) bool {
		id := ID{}
		stmt.ColumnBytes(1, id[:])
		rst[id] = stmt.ColumnInt(2) == 1
		return true
	}); err != nil {
		return nil, err
	}
	return rst, nil
}

func GetHareOutput(db database.Executor, lid uint32) ([]ID, error) {
	var rst []ID
	if err := db.Exec("select id from blocks where layer = ?1 and hare_output = 1;", func(stmt *database.Statement) {
		stmt.BindInt64(1, int64(lid))
	}, func(stmt *database.Statement) bool {
		id := ID{}
		stmt.ColumnBytes(1, id[:])
		rst = append(rst, id)
		return true
	}); err != nil {
		return nil, err
	}
	return rst, nil
}

func GetVerified(db database.Executor, lid uint32) ([]ID, error) {
	var rst []ID
	if err := db.Exec("select id from blocks where layer = ?1 and verified = 1;", func(stmt *database.Statement) {
		stmt.BindInt64(1, int64(lid))
	}, func(stmt *database.Statement) bool {
		id := ID{}
		stmt.ColumnBytes(1, id[:])
		rst = append(rst, id)
		return true
	}); err != nil {
		return nil, err
	}
	return rst, nil
}
