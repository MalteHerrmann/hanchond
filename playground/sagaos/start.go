package sagaos

import "fmt"

// TODO: refactor to be reuse the same method from Evmos?
// maybe not if this is requiring different settings to start, e.g. setting fee payer priv
func (s *SagaOS) Start() (int, error) {
	cmd := fmt.Sprintf("%s start --chain-id %s --home %s --json-rpc.enable true --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --grpc.enable >> %s 2>&1",
		s.GetVersionedBinaryPath(),
		s.ChainID,
		s.HomeDir,
		s.GetLogPath(),
	)
	return s.Daemon.Start(cmd)
}
