package cosmosdaemon

// StartOptions defines the extra options that can
// be set when starting a network or an individual node.
type StartOptions struct {
	LogLevel string `json:"log_level"`
}
