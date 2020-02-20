package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type CoinApp struct {
	Transactions []*Transaction
}

func (app *CoinApp) ToString() string {
	ret := "\n"
	for i, tx := range app.Transactions {
		ret += fmt.Sprintf("Txs[%d].%+v\n", i, tx)
	}
	return ret
}

func (app *CoinApp) GobEncode() ([]byte, error) {
	var result bytes.Buffer
	err := gob.NewEncoder(&result).Encode(app.Transactions)
	return result.Bytes(), err
}

func (app *CoinApp) GobDecode(d []byte) error {
	return gob.NewDecoder(bytes.NewReader(d)).Decode(&app.Transactions)
}
