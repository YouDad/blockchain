package core

import (
	"fmt"
	"strings"

	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
)

type Transaction struct {
	Hash types.HashValue
	Vin  []TxnInput
	Vout []TxnOutput
}

// func (txn *Transaction) GobEncode() ([]byte, error) {
//     result := bytes.Buffer{}
//     err := gob.NewEncoder(&result).Encode(*txn)
//     return result.Bytes(), err
// }
//
// func (txn *Transaction) GobDecode(b []byte) error {
//     return gob.NewDecoder(bytes.NewReader(b)).Decode(txn)
// }
//
func NewCoinbaseTxn(from string) *Transaction {
	txn := Transaction{}

	txn.Vin = []TxnInput{TxnInput{VoutIndex: -1}}
	// Send $from 50BTC
	txn.Vout = []TxnOutput{*NewTxnOutput(from, 50_000_000)}

	txn.Hash = utils.SHA256(&txn)
	return &txn
}

func (txn Transaction) String() string {
	lines := []string{}

	lines = append(lines, fmt.Sprintf("\t\tTxnHash %x:", utils.SHA256(txn)))

	for i, input := range txn.Vin {

		lines = append(lines, fmt.Sprintf("\t\t+ Input %d:", i))
		lines = append(lines, fmt.Sprintf("\t\t  - VoutHash:   %x", input.VoutHash))
		lines = append(lines, fmt.Sprintf("\t\t  - VoutIndex:  %d", input.VoutIndex))
		lines = append(lines, fmt.Sprintf("\t\t  - Signature:  %x", input.Signature))
		lines = append(lines, fmt.Sprintf("\t\t  - PubKeyHash: %x", input.PubKeyHash))
	}

	for i, output := range txn.Vout {
		lines = append(lines, fmt.Sprintf("\t\t+ Output %d:", i))
		lines = append(lines, fmt.Sprintf("\t\t  - Value:      %d", output.Value))
		lines = append(lines, fmt.Sprintf("\t\t  - PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

func (txn Transaction) IsCoinbase() bool {
	return len(txn.Vin) == 1 && txn.Vin[0].VoutIndex == -1
}
