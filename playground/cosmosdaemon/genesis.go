package cosmosdaemon

import (
	"fmt"
	"os/exec"
)

func (d *Daemon) GetGenesisSubcommand(args []string) ([]string, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments to get genesis subcommand")
	}

	subcommand := args[0]

	// check if `$BIN subcommand -h` is successful
	command := exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		subcommand,
		"-h",
	)

	if err := command.Run(); err == nil {
		return args, nil
	}

	// check if `$BIN genesis subcommand -h` is successful
	command = exec.Command( //nolint:gosec
		d.GetVersionedBinaryPath(),
		"genesis",
		subcommand,
		"-h",
	)

	if err := command.Run(); err == nil {
		args = append([]string{"genesis"}, args...)

		return args, nil
	}

	return nil, fmt.Errorf("no corresponding cli command found for %s", subcommand)
}

func (d *Daemon) UpdateGenesisFile() error {
	genesis, err := d.OpenGenesisFile()
	if err != nil {
		return err
	}
	// Update the genesis
	d.setStaking(genesis)
	d.setEvm(genesis)
	d.setInflation(genesis)
	d.setCrisis(genesis)
	d.setMint(genesis)
	d.setProvider(genesis)
	d.setConsensusParams(genesis)
	d.setFeeMarket(genesis)
	d.setGovernance(genesis, true)

	return d.SaveGenesisFile(genesis)
}

func (d *Daemon) setStaking(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["staking"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					// Base Denom
					if _, ok := v["base_denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["staking"].(map[string]any)["params"].(map[string]any)["bond_denom"] = d.BaseDenom
					}

					// Bond denom
					if _, ok := v["bond_denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["staking"].(map[string]any)["params"].(map[string]any)["bond_denom"] = d.BaseDenom
					}
				}
			}
		}
	}
}

func (d *Daemon) setEvm(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["evm"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					if _, ok := v["evm_denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["evm"].(map[string]any)["params"].(map[string]any)["evm_denom"] = d.BaseDenom
					}
				}
			}
		}
	}
}

func (d *Daemon) setInflation(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["inflation"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					if _, ok := v["mint_denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["inflation"].(map[string]any)["params"].(map[string]any)["mint_denom"] = d.BaseDenom
					}
				}
			}
		}
	}
}

func (d *Daemon) setCrisis(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}
	if v, ok := appState["crisis"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["constant_fee"]; ok {
				if v, ok := v.(map[string]any); ok {
					if _, ok := v["denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["crisis"].(map[string]any)["constant_fee"].(map[string]any)["denom"] = d.BaseDenom
					}
				}
			}
		}
	}
}

func (d *Daemon) setMint(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["mint"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					if _, ok := v["mint_denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["mint"].(map[string]any)["params"].(map[string]any)["mint_denom"] = d.BaseDenom
					}
				}
			}
		}
	}
}

func (d *Daemon) setProvider(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["provider"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					if v, ok := v["consumer_reward_denom_registration_fee"]; ok {
						if v, ok := v.(map[string]any); ok {
							if _, ok := v["denom"]; ok {
								//nolint:forcetypeassert // ok like this
								appState["provider"].(map[string]any)["params"].(map[string]any)["consumer_reward_denom_registration_fee"].(map[string]any)["denom"] = d.BaseDenom
							}
						}
					}
				}
			}
		}
	}
}

func (d *Daemon) setConsensusParams(genesis map[string]any) {
	var consensusParams map[string]any
	if _, ok := genesis["consensus_params"]; ok {
		//nolint:forcetypeassert // ok like this
		consensusParams = genesis["consensus_params"].(map[string]any)
	}

	// SDKv0.50 support
	if _, ok := genesis["consensus"]; ok {
		//nolint:forcetypeassert // ok like this
		consensusParams = genesis["consensus"].(map[string]any)["params"].(map[string]any)
	}

	if v, ok := consensusParams["block"]; ok {
		if v, ok := v.(map[string]any); ok {
			if _, ok := v["max_gas"]; ok {
				//nolint:forcetypeassert // ok like this
				consensusParams["block"].(map[string]any)["max_gas"] = d.GasLimit
			}
		}
	}
}

func (d *Daemon) setFeeMarket(genesis map[string]any) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["feemarket"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					// Evmos FeeMarket
					if _, ok := v["base_fee"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["feemarket"].(map[string]any)["params"].(map[string]any)["base_fee"] = d.BaseFee
						// FeeMarket using static base fee
						//
						//nolint:forcetypeassert // ok like this
						appState["feemarket"].(map[string]any)["params"].(map[string]any)["base_fee_change_denominator"] = 1
						//nolint:forcetypeassert // ok like this
						appState["feemarket"].(map[string]any)["params"].(map[string]any)["elasticity_multiplier"] = 1
						//nolint:forcetypeassert // ok like this
						appState["feemarket"].(map[string]any)["params"].(map[string]any)["min_gas_multiplier"] = "0.0"
					}
					// SDK FeeMarket
					if _, ok := v["fee_denom"]; ok {
						//nolint:forcetypeassert // ok like this
						appState["feemarket"].(map[string]any)["params"].(map[string]any)["fee_denom"] = d.BaseDenom
					}
				}
			}
		}
	}
}

func (d *Daemon) setGovernance(genesis map[string]any, fastProposals bool) {
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	if v, ok := appState["gov"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["params"]; ok {
				if v, ok := v.(map[string]any); ok {
					// Proposals
					if fastProposals {
						if _, ok := v["max_deposit_period"]; ok {
							//nolint:forcetypeassert // ok like this
							appState["gov"].(map[string]any)["params"].(map[string]any)["max_deposit_period"] = "10s"
						}

						if _, ok := v["voting_period"]; ok {
							//nolint:forcetypeassert // ok like this
							appState["gov"].(map[string]any)["params"].(map[string]any)["voting_period"] = "15s"
						}

						if _, ok := v["expedited_voting_period"]; ok {
							//nolint:forcetypeassert // ok like this
							appState["gov"].(map[string]any)["params"].(map[string]any)["expedited_voting_period"] = "14s"
						}
					}

					//  Expedited_min_deposit
					if v, ok := v["expedited_min_deposit"]; ok {
						if v, ok := v.([]any); ok {
							if len(v) > 0 {
								if v, ok := v[0].(map[string]any); ok {
									if _, ok := v["denom"]; ok {
										//nolint:forcetypeassert // ok like this
										appState["gov"].(map[string]any)["params"].(map[string]any)["expedited_min_deposit"].([]any)[0].(map[string]any)["denom"] = d.BaseDenom
									}
								}
							}
						}
					}

					// Min Deposit
					if v, ok := v["min_deposit"]; ok {
						if v, ok := v.([]any); ok {
							if len(v) > 0 {
								if v, ok := v[0].(map[string]any); ok {
									if _, ok := v["denom"]; ok {
										//nolint:forcetypeassert // ok like this
										appState["gov"].(map[string]any)["params"].(map[string]any)["min_deposit"].([]any)[0].(map[string]any)["denom"] = d.BaseDenom
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (d *Daemon) setBank() error {
	genesis, err := d.OpenGenesisFile()
	if err != nil {
		return err
	}
	appState, ok := genesis["app_state"].(map[string]any)
	if !ok {
		panic("unexpected app_state")
	}

	//nolint:forcetypeassert // ok like this
	appState["bank"].(map[string]any)["supply"].([]any)[0].(map[string]any)["amount"] = d.ValidatorInitialSupply

	if err := d.SaveGenesisFile(genesis); err != nil {
		return err
	}

	return nil
}
