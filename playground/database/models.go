// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

type Chain struct {
	ID            int64
	Name          string
	ChainID       string
	BinaryVersion string
}

type Node struct {
	ID               int64
	ChainID          int64
	ConfigFolder     string
	Moniker          string
	ValidatorKey     string
	ValidatorKeyName string
	BinaryVersion    string
	ProcessID        int64
	IsValidator      int64
	IsArchive        int64
	IsRunning        int64
}

type Port struct {
	ID     int64
	NodeID int64
	P1317  int64
	P8080  int64
	P9090  int64
	P9091  int64
	P8545  int64
	P8546  int64
	P6065  int64
	P26658 int64
	P26657 int64
	P6060  int64
	P26656 int64
	P26660 int64
}

type Relayer struct {
	ID        int64
	ProcessID int64
	IsRunning int64
}
