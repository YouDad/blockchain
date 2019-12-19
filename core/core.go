package core

import (
	"github.com/YouDad/blockchain/app"
)

type CoreConfig struct {
	GetAppdata func() app.App
	GetGenesis func() app.App
	// DatabaseFile string
	// BlocksBucket string
}

var (
	coreConfig CoreConfig
)

func InitCore(config CoreConfig) {
	if config.GetAppdata != nil {
		coreConfig.GetAppdata = config.GetAppdata
	}

	if config.GetGenesis != nil {
		coreConfig.GetGenesis = config.GetGenesis
	}
}
