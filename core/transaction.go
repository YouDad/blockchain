package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/YouDad/blockchain/log"
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

func BytesToTransaction(bytes []byte) *Transaction {
	txn := Transaction{}
	log.Err(utils.GetDecoder(bytes).Decode(&txn))
	return &txn
}

func (txn Transaction) IsCoinbase() bool {
	return len(txn.Vin) == 1 && txn.Vin[0].VoutIndex == -1
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxnInput
	var outputs []TxnOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TxnInput{
			VoutHash:   vin.VoutHash,
			VoutIndex:  vin.VoutIndex,
			Signature:  nil,
			PubKeyHash: nil,
		})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TxnOutput{
			Value:      vout.Value,
			PubKeyHash: vout.PubKeyHash,
		})
	}

	txCopy := Transaction{
		Hash: tx.Hash,
		Vin:  inputs,
		Vout: outputs,
	}

	return txCopy
}

func (txn *Transaction) Sign(sk types.PrivateKey, hashedTxn map[string]Transaction) {
	if txn.IsCoinbase() {
		return
	}

	for _, vin := range txn.Vin {
		if hashedTxn[hex.EncodeToString(vin.VoutHash)].Hash == nil {
			log.Errln("Prev Hashed Txn is not correct")
		}
	}

	txnCopy := txn.TrimmedCopy()

	for inIndex, vin := range txnCopy.Vin {
		prevTxn := hashedTxn[hex.EncodeToString(vin.VoutHash)]
		txnCopy.Vin[inIndex].PubKeyHash = prevTxn.Vout[vin.VoutIndex].PubKeyHash
		dataToSign := []byte(fmt.Sprintf("%x\n", txnCopy))

		r, s, err := ecdsa.Sign(rand.Reader, &sk, dataToSign)
		log.Err(err)
		signature := append(r.Bytes(), s.Bytes()...)

		txn.Vin[inIndex].Signature = signature
		txnCopy.Vin[inIndex].PubKeyHash = nil
	}
}

func (txn *Transaction) Verify(hashedTxn map[string]Transaction) bool {
	if txn.IsCoinbase() {
		return true
	}

	for _, vin := range txn.Vin {
		if hashedTxn[hex.EncodeToString(vin.VoutHash)].Hash == nil {
			log.Errln("Previous transaction is not correct")
		}
	}

	txnCopy := txn.TrimmedCopy()
	curve := elliptic.P256()

	for inIndex, vin := range txn.Vin {
		prevTxn := hashedTxn[hex.EncodeToString(vin.VoutHash)]
		txnCopy.Vin[inIndex].PubKeyHash = prevTxn.Vout[vin.VoutIndex].PubKeyHash
		dataToVerify := []byte(fmt.Sprintf("%x\n", txnCopy))

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKeyHash)
		x.SetBytes(vin.PubKeyHash[:(keyLen / 2)])
		y.SetBytes(vin.PubKeyHash[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, dataToVerify, &r, &s) {
			return false
		}
		txnCopy.Vin[inIndex].PubKeyHash = nil
	}

	return true
}
