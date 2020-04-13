package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	gob.Register(elliptic.P256())
}

// Wallets stores a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// NewWallets creates Wallets and fills it from a file if it exists
func NewWallets(walletFile string) (*Wallets, error) {
	wallets := Wallets{Wallets: make(map[string]*Wallet)}
	walletFile = fmt.Sprintf("wallet%s.dat", walletFile)

	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		return &wallets, err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
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
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// SaveToFile saves wallets to a file
func (ws Wallets) SaveToFile(walletFile string) {
	var content bytes.Buffer

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	walletFile = fmt.Sprintf("wallet%s.dat", walletFile)
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
