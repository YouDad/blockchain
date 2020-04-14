package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

type Wallets map[string]*Wallet

var (
	onceWallets     sync.Once
	instanceWallets Wallets
	errWallets      error
)

func init() {
	gob.Register(elliptic.P256())
}

func getWallets() (Wallets, error) {
	wallets := make(Wallets)
	filename := fmt.Sprintf("wallet%s.dat", global.Port)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return wallets, nil
	}

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return wallets, err
	}

	reader := bytes.NewReader(fileContent)
	decoder := gob.NewDecoder(reader)
	decoder.Decode(&wallets)

	return wallets, err
}

func GetWallets() (Wallets, error) {
	onceWallets.Do(func() {
		instanceWallets, errWallets = getWallets()
	})
	return instanceWallets, errWallets
}

func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	encoder := gob.NewEncoder(&content)
	log.Err(encoder.Encode(ws))

	walletFile := fmt.Sprintf("wallet%s.dat", global.Port)
	log.Err(ioutil.WriteFile(walletFile, content.Bytes(), 0644))
}
