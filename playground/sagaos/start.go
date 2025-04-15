package sagaos

import "fmt"

// TODO: refactor to be reuse the same method from Evmos?
func (s *SagaOS) Start() (int, error) {
	logFile := s.HomeDir + "/run.log"
	cmd := fmt.Sprintf("%s start --chain-id %s --home %s --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --grpc.enable >> %s 2>&1",
		s.GetVersionedBinaryPath(),
		s.ChainID,
		s.HomeDir,
		logFile,
	)
	return s.Daemon.Start(cmd)
}
