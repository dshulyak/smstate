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