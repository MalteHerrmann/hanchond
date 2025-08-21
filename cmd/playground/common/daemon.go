package commoncmd

import (
	"fmt"

	"github.com/hanchon/hanchond/playground/cosmosdaemon"
	"github.com/hanchon/hanchond/playground/evmos"
	"github.com/hanchon/hanchond/playground/gaia"
	"github.com/hanchon/hanchond/playground/noble"
	"github.com/hanchon/hanchond/playground/orbiter"
	"github.com/hanchon/hanchond/playground/sagaos"
	"github.com/hanchon/hanchond/playground/types"
)

func GetDaemonForNode(n types.DaemonInfo, ports *types.Ports) (cosmosdaemon.IDaemon, error) {
	binaryName := n.GetChainInfo().GetBinaryName()

	var d cosmosdaemon.IDaemon
	switch binaryName {
	case gaia.ChainInfo.GetBinaryName():
		d = gaia.NewGaia(
			n.GetMoniker(),
			n.GetConfigFolder(),
			n.GetChainID(),
			n.GetValidatorKeyName(),
			ports,
		)
	case evmos.ChainInfo.GetBinaryName():
		d = evmos.NewEvmos(
			n.GetMoniker(),
			n.GetVersion(),
			n.GetConfigFolder(),
			n.GetChainID(),
			n.GetValidatorKeyName(),
			ports,
		)
	case sagaos.ChainInfo.GetBinaryName():
		d = sagaos.NewSagaOS(
			n.GetMoniker(),
			n.GetVersion(),
			n.GetConfigFolder(),
			n.GetChainID(),
			n.GetValidatorKeyName(),
			ports,
		)
	// TODO: instead of simd abstract this away to enable multiple things that are called
	// simd
	case orbiter.ChainInfo.GetBinaryName():
		d = orbiter.NewOrbiter(
			n.GetMoniker(),
			n.GetVersion(),
			n.GetConfigFolder(),
			n.GetChainID(),
			n.GetValidatorKeyName(),
			ports,
		)
	case noble.ChainInfo.GetBinaryName():
		d = noble.NewNoble(
			n.GetMoniker(),
			n.GetVersion(),
			n.GetConfigFolder(),
			n.GetChainID(),
			n.GetValidatorKeyName(),
			ports,
		)
	default:
		return nil, fmt.Errorf("binary %s not configured", binaryName)
	}

	return d, nil
}
