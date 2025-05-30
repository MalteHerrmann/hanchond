-- TODO: update this to include the daemon info like chain info etc. -> chain info probably needs to go into the chain schema
CREATE TABLE IF NOT EXISTS chain(
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    chain_id TEXT NOT NULL,
    -- TODO: should this even be moved to a separate table? not really necessary I guess..
    chain_info TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS node(
    id INTEGER NOT NULL PRIMARY KEY,
    chain_id INTEGER NOT NULL,
    config_folder TEXT NOT NULL,
    moniker TEXT NOT NULL,
    validator_key TEXT NOT NULL,
    validator_key_name TEXT NOT NULL,
    validator_wallet TEXT NOT NULL,
    key_type TEXT NOT NULL,
    -- TODO: remove this? Not sure if storing the version is required both for chain and node here - it probably makes sense though since the chain can have nodes on different patch versions communicating for example.
    version TEXT NOT NULL,

    process_id INTEGER NOT NULL,
    is_validator INTEGER NOT NULL,
    is_archive INTEGER NOT NULL,
    is_running INTEGER NOT NULL,

    FOREIGN KEY (chain_id) REFERENCES chain(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ports(
    id INTEGER NOT NULL PRIMARY KEY,
    node_id INTEGER NOT NULL UNIQUE ,
	p1317 INTEGER NOT NULL,
	p8080 INTEGER NOT NULL,
	p9090 INTEGER NOT NULL,
	p9091 INTEGER NOT NULL,
	p8545 INTEGER NOT NULL,
	p8546 INTEGER NOT NULL,
	p6065 INTEGER NOT NULL,
	p26658 INTEGER NOT NULL,
	p26657 INTEGER NOT NULL,
	p6060  INTEGER NOT NULL,
	p26656 INTEGER NOT NULL,
	p26660 INTEGER NOT NULL,

    FOREIGN KEY (node_id) REFERENCES node(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS relayer(
    id INTEGER NOT NULL PRIMARY KEY,
    process_id INTEGER NOT NULL,
    is_running INTEGER NOT NULL
);
