package cmd

type Config struct {
	AuthPassEnvVar string
	AuthUserEnvVar string
	DefaultEnv     string
	ProjectRoot    string
	GitBinPath     string
	UploadPack     bool
	ReceivePack    bool
	RemoteProxyURL string

	RemoteRustStaticURL string
}

var (
	DefaultConfig *Config
)
