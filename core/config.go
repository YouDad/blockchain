package core

type Config struct {
	DatabaseFile string
	WalletFile   string
}

var (
	CoreConfig Config
)

func InitCore(config Config) {
	if config.DatabaseFile != "" {
		CoreConfig.DatabaseFile = config.DatabaseFile
	}

	if config.WalletFile != "" {
		CoreConfig.WalletFile = config.WalletFile
	}
}
