# AGENTS.md

This document provides essential context for AI agents working with this codebase. It describes the project structure, key patterns, and important conventions to help agents understand and navigate the codebase effectively.

## Project Overview

**Hanchond** is a comprehensive web3 toolkit for Cosmos SDK-based blockchain networks. It's a CLI application written in Go that provides utilities for:

- Managing local Cosmos chain networks (playground)
- Building and managing blockchain binaries (Evmos, Gaia, Noble, Orbiter, Saga OS)
- Interacting with Cosmos chains (queries, transactions, IBC transfers)
- Converting between address formats, number formats, and encodings
- Managing IBC relayers (Hermes)
- Building and deploying Solidity smart contracts
- Exploring blockchain data with a TUI (Terminal User Interface)

The project follows a modular architecture with clear separation between CLI commands, core libraries, and playground-specific functionality.

## Directory Structure

### `/cmd` - CLI Commands
Contains all Cobra CLI command definitions:
- **`cmd/root.go`**: Root command that registers all subcommands
- **`cmd/convert/`**: Conversion utilities (addresses, numbers)
- **`cmd/playground/`**: Main playground commands (tx, query, relayer, explorer)
- **`cmd/repo/`**: Repository management commands
- **`cmd/version.go`**: Version command

**Pattern**: Each command package exports a `*cobra.Command` variable (e.g., `ConvertCmd`, `PlaygroundCmd`) that is registered in `cmd/root.go`.

### `/lib` - Core Libraries
Reusable libraries used across the application:

- **`lib/converter/`**: Format conversion utilities
  - `address.go`: Address format conversions (Hex ↔ Bech32)
  - `number.go`: Number format conversions
  - `hash.go`: Hash utilities
  - `base64.go`: Base64 encoding/decoding

- **`lib/requester/`**: HTTP client libraries
  - `cosmos.go`: Cosmos SDK REST API client
  - `erc20.go`: ERC20 token interaction client
  - `web3.go`: Ethereum JSON-RPC client
  - `tendermint.go`: Tendermint RPC client
  - `get.go`, `post.go`: Generic HTTP utilities

- **`lib/txbuilder/`**: Transaction building and wallet management
  - `wallet.go`: Wallet creation from mnemonics (Cosmos and Ethereum)
  - `transaction.go`: Transaction construction
  - `transfer.go`: Token transfer transactions
  - `contract.go`, `contract_deployment.go`, `contract_interaction.go`: Smart contract operations
  - `gas.go`: Gas estimation utilities

- **`lib/smartcontract/`**: Smart contract abstractions
  - `erc20/`: ERC20 token contract interface

- **`lib/types/`**: Type definitions
  - `cosmos/`: Cosmos SDK types (blocks, transactions)
  - `evmos/`: Evmos-specific types
  - `tendermint/`: Tendermint types
  - `web3/`: Ethereum/Web3 types
  - `protocol/`: Protocol-level types

- **`lib/protoencoder/`**: Protocol buffer encoding/decoding

- **`lib/utils/`**: Utility functions (logging, exit handling, execution)

### `/playground` - Playground-Specific Code
Code specific to the local chain playground functionality:

- **`playground/cosmosdaemon/`**: Core daemon management
  - `daemon.go`: `IDaemon` interface and `Daemon` struct for chain-agnostic daemon operations
  - `init.go`: Chain initialization
  - `keys.go`: Key management
  - `genesis.go`: Genesis file generation
  - `config_files.go`: Configuration file management
  - `ports.go`: Port allocation
  - `database.go`: Database operations

- **`playground/database/`**: SQLite database layer
  - `db.go`: Database connection and setup
  - `models.go`: Data models
  - `query.go`: Database queries (uses sqlc for type-safe queries)

- **`playground/explorer/`**: Block explorer
  - `explorer.go`: Explorer backend
  - `explorerui/`: TUI components (using Charmbracelet)

- **`playground/types/`**: Playground-specific types
  - `chain_info.go`: Chain configuration and metadata
  - `coin.go`: Coin/denom types
  - `daemon_info.go`: Daemon runtime information
  - `ports.go`: Port configuration

- **Chain-specific implementations**:
  - `playground/evmos/`: Evmos chain implementation
  - `playground/gaia/`: Gaia/Cosmos Hub implementation
  - `playground/noble/`: Noble chain implementation
  - `playground/orbiter/`: Orbiter chain implementation
  - `playground/sagaos/`: Saga OS chain implementation

- **`playground/hermes/`**: Hermes IBC relayer integration
  - `hermes.go`: Hermes binary management
  - `config.go`: Hermes configuration
  - `logs.go`: Log management

- **`playground/solidity/`**: Solidity compiler integration
  - `compiler.go`: Solidity compilation
  - `erc20builder.go`: ERC20 contract builder

- **`playground/filesmanager/`**: File and binary management
  - `path.go`: Path utilities
  - `git.go`: Git operations for cloning repos
  - `list_versions.go`: Version listing
  - `build_evmos.go`: Evmos binary building

- **`playground/sql/`**: SQL schema and queries
  - `schema.sql`: Database schema
  - `query.sql`: SQL queries (used by sqlc)

### `/docs` - Documentation
MDX documentation files for the project website.

### `/examples` - Example Code
Example usage of the libraries.

### `/abis` - Contract ABIs
Ethereum contract ABIs (e.g., `erc20.abi`).

## Key Patterns and Conventions

### Command Structure
- Commands use Cobra and are defined in `/cmd`
- Each command package exports a `*cobra.Command` variable
- Commands are registered in `cmd/root.go` via `rootCmd.AddCommand()`
- Persistent flags (like `--home`) are set in `init()` functions

### Error Handling
- Use `lib/utils/ExitError()` for fatal errors
- Use `lib/utils/ExitSuccess()` for successful exits
- Functions typically return `(result, error)` tuples

### Chain Abstraction
- The `IDaemon` interface in `playground/cosmosdaemon/daemon.go` defines chain-agnostic operations
- Chain-specific implementations live in `playground/{chain}/` directories
- Chain metadata is stored in `playground/types/chain_info.go` as `ChainInfo` structs
- Chain info includes: account prefix, binary name, chain ID base, denom, HD path, key algorithm, repo URL

### Database Usage
- SQLite database stores chain, node, port, and relayer information
- Schema defined in `playground/sql/schema.sql`
- Uses sqlc for type-safe queries (see `sqlc.yaml`)
- Database models in `playground/database/models.go`

### Address and Wallet Handling
- Supports both Cosmos (secp256k1) and Ethereum (eth_secp256k1) key algorithms
- Cosmos uses Bech32 addresses with configurable prefixes
- Ethereum uses hex addresses
- Conversion utilities in `lib/converter/address.go`
- Wallet creation in `lib/txbuilder/wallet.go` supports both HD paths:
  - Cosmos: `m/44'/118'/0'/0`
  - Ethereum: `m/44'/60'/0'/0`

### Port Management
- Each node gets a unique set of ports (RPC, REST, gRPC, etc.)
- Ports stored in database and managed via `playground/types/ports.go`
- Port allocation prevents conflicts

### Binary Management
- Binaries are built from source and stored in `~/.hanchond/binaries/`
- Version management via `playground/filesmanager/list_versions.go`
- Chain-specific build logic in `playground/{chain}/` directories

## Important Dependencies

- **Cobra** (`github.com/spf13/cobra`): CLI framework
- **Cosmos SDK** (`github.com/cosmos/cosmos-sdk`): Cosmos blockchain SDK (forked by Evmos)
- **Evmos** (`github.com/evmos/evmos/v18`): Evmos chain SDK
- **go-ethereum** (`github.com/ethereum/go-ethereum`): Ethereum client (forked by Evmos)
- **Charmbracelet** (`github.com/charmbracelet/*`): TUI libraries
- **sqlite** (`modernc.org/sqlite`): SQLite database driver
- **easyjson** (`github.com/mailru/easyjson`): Fast JSON marshaling
- **fasthttp** (`github.com/valyala/fasthttp`): Fast HTTP client

## Common Tasks

### Adding a New CLI Command
1. Create a new package in `/cmd` (e.g., `cmd/newcommand/`)
2. Define a `*cobra.Command` variable
3. Register it in `cmd/root.go` via `rootCmd.AddCommand()`

### Adding Support for a New Chain
1. Create a new directory in `/playground` (e.g., `playground/newchain/`)
2. Implement chain-specific initialization in `init.go` or similar
3. Define `ChainInfo` configuration
4. Add build logic if the chain needs custom binary building
5. Update any chain-specific queries or transactions

### Working with Transactions
- Use `lib/txbuilder/` for constructing transactions
- For Cosmos transactions: use `txbuilder` with Cosmos SDK types
- For Ethereum transactions: use `txbuilder` with go-ethereum types
- Gas estimation via `lib/txbuilder/gas.go`

### Querying Chain Data
- Use `lib/requester/cosmos.go` for Cosmos SDK REST queries
- Use `lib/requester/web3.go` for Ethereum JSON-RPC queries
- Use `lib/requester/tendermint.go` for Tendermint RPC queries

### Database Operations
- Schema changes: update `playground/sql/schema.sql`
- Query changes: update `playground/sql/query.sql`
- Run `sqlc generate` to regenerate Go code
- Models are in `playground/database/models.go`

## Testing and Development

- The project uses Go modules (`go.mod`)
- Main entry point: `main.go` → `cmd.Execute()`
- Home directory: `~/.hanchond` (configurable via `--home` flag)
- Database location: `{home}/playground.db`
- Binary storage: `{home}/binaries/`

## Notes for AI Agents

1. **When modifying commands**: Always check `cmd/root.go` to see how commands are registered
2. **When working with chains**: Check `playground/types/chain_info.go` for chain metadata structure
3. **When adding database fields**: Update both `schema.sql` and regenerate with sqlc
4. **When working with addresses**: Use `lib/converter/address.go` for format conversions
5. **When building transactions**: Use `lib/txbuilder/` - don't construct raw transactions
6. **Chain-specific code**: Look in `playground/{chain}/` directories
7. **Error handling**: Prefer returning errors rather than panicking (except in `Must*` functions)
8. **JSON marshaling**: The project uses easyjson for performance - generated files have `_easyjson.go` suffix

## File Naming Conventions

- Commands: `{command}.go` in `/cmd/{command}/`
- Library files: descriptive names (e.g., `wallet.go`, `transaction.go`)
- Generated files: `*_easyjson.go` for JSON marshaling, `*.sql.go` for sqlc queries
- Test files: `*_test.go` (though not many tests present currently)

## Key Entry Points

- **CLI**: `main.go` → `cmd.Execute()` → `cmd/root.go`
- **Playground**: `cmd/playground/playground.go`
- **Daemon operations**: `playground/cosmosdaemon/daemon.go`
- **Database**: `playground/database/db.go`
- **Chain info**: `playground/types/chain_info.go`

