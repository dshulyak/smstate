CREATE TABLE blocks ( 
	id CHAR(20) PRIMARY KEY,
	layer INT,
	hare_output BOOL,
	verified BOOL,
	block    BLOB
);

CREATE INDEX blocks_by_layer ON blocks(layer);
CREATE INDEX blocks_by_hare_output ON blocks(hare_output, layer) WHERE hare_output = 1;
CREATE INDEX blocks_by_verified ON blocks(verified, layer) where verified = 1;