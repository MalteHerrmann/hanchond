-- name: InsertChain :one
-- TODO: make sure that version is not required here? Maybe it's important to set here?
INSERT INTO chain(
    name, chain_id, chain_info
) VALUES (
    ?,?,?
)
RETURNING *;

-- name: InsertNode :one
INSERT INTO node(
    chain_id,
    config_folder,
    moniker,
    validator_key,
    validator_key_name,
    validator_wallet,
    key_type,
    version,
    process_id,
    is_validator,
    is_archive,
    is_running
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?,?
)
RETURNING ID;

-- name: InsertPorts :exec
INSERT INTO ports(
    node_id,
	p1317,
	p8080,
	p9090,
	p9091,
	p8545,
	p8546,
	p6065,
	p26658,
	p26657,
	p6060,
	p26656,
	p26660
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?,?,?
);

-- name: SetProcessID :exec
UPDATE node SET
    process_id = ?,
    is_running = ?
WHERE (
    id = ?
);

-- name: SetNodeVersion :exec
UPDATE node SET
    version = ?
WHERE (
    id = ?
);

-- name: GetNode :one
SELECT * FROM node where id =? LIMIT 1;

-- name: GetNodePorts :one
SELECT * FROM ports where node_id =? LIMIT 1;

-- name: GetAllPorts :many
SELECT * FROM ports;

-- name: GetChain :one
SELECT * FROM chain where id =? LIMIT 1;

-- name: GetLatestChain :one
SELECT * FROM chain ORDER BY id DESC LIMIT 1;

-- name: GetAllNodes :many
SELECT * FROM node;

-- name: GetAllChainNodes :many
SELECT * FROM node n join ports p on p.node_id == n.id join chain c on n.chain_id == c.id where n.chain_id = ?;

-- name: InitRelayer :exec
INSERT INTO relayer(
    process_id, is_running
) VALUES (
    0,0
);

-- name: GetRelayer :one
SELECT * FROM relayer WHERE id = 1;

-- name: UpdateRelayer :exec
UPDATE relayer SET process_id = ?, is_running = ? WHERE id = 1;

