package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/YouDad/blockchain/global"
	"github.com/YouDad/blockchain/log"
)

var (
	walletFilename string
)

func Register() {
	gob.Register(elliptic.P256())
	walletFilename = fmt.Sprintf("wallet%s.dat", global.Port)
}

// Wallets stores a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// NewWallets creates Wallets and fills it from a file if it exists
func NewWallets() (*Wallets, error) {
	wallets := Wallets{Wallets: make(map[string]*Wallet)}

	_, err := os.Stat(walletFilename)
	if os.IsNotExist(err) {
		return &wallets, err
	}

	fileContent, err := ioutil.ReadFile(walletFilename)
	if err != nil {
		return &wallets, err
	}

	reader := bytes.NewReader(fileContent)
	decoder := gob.NewDecoder(reader)
	decoder.Decode(&wallets)

	return &wallets, err
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet returns a Wallet by its address
func (ws Wallets) GetWallet(address string) (*Wallet, bool) {
	wallet, ok := ws.Wallets[address]
	return wallet, ok
}

// SaveToFile saves wallets to a file
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	log.Err(err)

	walletFile := fmt.Sprintf("wallet%s.dat", global.Port)
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	log.Err(err)
}
