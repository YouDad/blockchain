package coin

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/YouDad/blockchain/app"
	"github.com/YouDad/blockchain/core"
)

type CoinApp struct {
	app.App
	Transactions []*Transaction
}

func Init() {
	core.InitCore(core.CoreConfig{
		GetAppdata: func() app.App {
			return GetCoinApp(nil)
		},
	})
}

func GetCoinApp(txs []*Transaction) *CoinApp {
	return &CoinApp{Transactions: txs}
}

func (app *CoinApp) HashPart() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range app.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (app *CoinApp) ToString() string {
	ret := "\n"
	for i, tx := range app.Transactions {
		ret += fmt.Sprintf("Txs[%d]: %+v\n", i, tx)
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
