package coin_core

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/core"
)

type CoinApp struct {
	Transactions []*Transaction
}

func Init() {
	gob.Register(CoinApp{
		Transactions: []*Transaction{
			&Transaction{},
		},
	})
	core.InitCore(core.Config{
		GetAppdata: func() app.App {
			return GetCoinApp(nil)
		},
	})
}

func GetCoinApp(txs []*Transaction) *CoinApp {
	return &CoinApp{Transactions: txs}
}

func (app *CoinApp) HashPart() []byte {
	var txs [][]byte

	for _, tx := range app.Transactions {
		txs = append(txs, tx.Serialize())
	}
	mTree := NewMerkleTree(txs)

	return mTree.RootNode.Data
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
