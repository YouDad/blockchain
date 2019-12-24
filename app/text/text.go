package text

import (
	"encoding/gob"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/core"
)

type TextApp struct {
	str string
}

func Init() {
	gob.Register(TextApp{})
	core.InitCore(core.Config{
		GetAppdata: func() app.App {
			return GetAppString("")
		},
		GetGenesis: func() app.App {
			return GetAppString("Genesis Block")
		},
	})
}

func GetAppString(str string) *TextApp {
	return &TextApp{str: str}
}

func (app *TextApp) HashPart() []byte {
	return []byte(app.str)
}

func (app *TextApp) ToString() string {
	return app.str
}

func (app *TextApp) GobEncode() ([]byte, error) {
	return []byte(app.str), nil
}

func (app *TextApp) GobDecode(data []byte) error {
	app.str = string(data)
	return nil
}
