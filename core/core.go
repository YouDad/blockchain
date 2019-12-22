package core

import (
	"github.com/YouDad/blockchain/app"
)

type Config struct {
	GetAppdata   func() app.App
	GetGenesis   func() app.App
	DatabaseFile string
	WalletFile   string
}

var (
	CoreConfig Config
)

func InitCore(config Config) {
	if config.GetAppdata != nil {
		CoreConfig.GetAppdata = config.GetAppdata
	}

	if config.GetGenesis != nil {
		CoreConfig.GetGenesis = config.GetGenesis
	}

	if config.DatabaseFile != "" {
		CoreConfig.DatabaseFile = config.DatabaseFile
	}

	if config.WalletFile != "" {
		CoreConfig.WalletFile = config.WalletFile
	}
}
