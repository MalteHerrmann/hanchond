package config

import (
	_ "embed"
	"fmt"

	"github.com/hanchon/hanchond/playground/types"
	"gopkg.in/yaml.v3"
)

//go:embed chain_config.yaml
var configData []byte

type BuildConfig struct {
	MakeTarget  string `yaml:"make_target"`
	MakeOptions string `yaml:"make_options"`
	BinaryPath  string `yaml:"binary_path"`
}

// TODO: check if all of the fields are being used. If not there is probably some refactoring necessary to use that. E.g. use the `make` target and options to build the binary instead of hardcoding them.
type ChainConfig struct {
	AccountPrefix string      `yaml:"account_prefix"`
	BinaryName    string      `yaml:"binary_name"`
	ChainIDBase   string      `yaml:"chain_id_base"`
	ClientName    string      `yaml:"client_name"`
	Denom         string      `yaml:"denom"`
	RepoURL       string      `yaml:"repo_url"`
	HDPath        string      `yaml:"hd_path"`
	KeyAlgo       string      `yaml:"key_algo"`
	SDKVersion    string      `yaml:"sdk_version"`
	Build         BuildConfig `yaml:"build"`
}

type Config struct {
	Chains map[string]ChainConfig `yaml:"chains"`
}

var config *Config

func init() {
	if err := loadConfig(); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
}

// loadConfig loads the chain configuration from the embedded YAML file
func loadConfig() error {
	config = &Config{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}
	return nil
}

// GetChainConfig returns the configuration for a specific chain
//
// TODO: this should actually not be the chain name aka chain ID but rather the client's name.
// e.g. for evmos it should be "evmos" instead of "evmos_9001-1" or "evmosd"
func GetChainConfig(chainName string) (*ChainConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	chainConfig, exists := config.Chains[chainName]
	if !exists {
		return nil, fmt.Errorf("chain %s not found in configuration", chainName)
	}

	return &chainConfig, nil
}

// ToChainInfo converts a ChainConfig to a types.ChainInfo
func (c *ChainConfig) ToChainInfo() types.ChainInfo {
	return types.NewChainInfo(
		c.AccountPrefix,
		c.BinaryName,
		c.ChainIDBase,
		c.ClientName,
		c.Denom,
		c.RepoURL,
		types.HDPath(c.HDPath),
		types.SignatureAlgo(c.KeyAlgo),
		types.SDKVersion(c.SDKVersion),
	)
}
