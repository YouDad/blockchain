package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/utils"
)

type Transaction struct {
	Vin  []TxnInput
	Vout []TxnOutput
}

func (txn Transaction) Hash() HashValue {
	return utils.SHA256(txn)
}

func (txn Transaction) String() string {
	return string(utils.Encode(txn))
}

func (txn Transaction) IsCoinbase() bool {
	return len(txn.Vin) == 1 && txn.Vin[0].VoutIndex == -1
}

func (txn Transaction) TrimmedCopy() Transaction {
	var inputs []TxnInput
	var outputs []TxnOutput

	for _, vin := range txn.Vin {
		inputs = append(inputs, TxnInput{
			VoutHash:  vin.VoutHash,
			VoutIndex: vin.VoutIndex,
			VoutValue: vin.VoutValue,
			Signature: nil,
			PubKey:    nil,
		})
	}

	for _, vout := range txn.Vout {
		outputs = append(outputs, TxnOutput{
			Value:      vout.Value,
			PubKeyHash: vout.PubKeyHash,
		})
	}

	txCopy := Transaction{
		Vin:  inputs,
		Vout: outputs,
	}

	return txCopy
}

func (txn *Transaction) Sign(sk PrivateKey, hashedTxn map[string]Transaction) error {
	if txn.IsCoinbase() {
		return nil
	}

	txnCopy := txn.TrimmedCopy()

	for inIndex, vin := range txnCopy.Vin {
		prevTxn := hashedTxn[vin.VoutHash.String()]
		if len(prevTxn.Vout) <= vin.VoutIndex {
			log.Errln(txn, hashedTxn)
		}
		txnCopy.Vin[inIndex].PubKey = PublicKey(prevTxn.Vout[vin.VoutIndex].PubKeyHash)
		dataToSign := []byte(fmt.Sprintf("%s\n", txnCopy))

		r, s, err := ecdsa.Sign(rand.Reader, &sk, dataToSign)
		if err != nil {
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		txn.Vin[inIndex].Signature = signature
		txnCopy.Vin[inIndex].PubKey = nil
	}
	return nil
}

// 验证交易是否有效，需要map[前置交易哈希]前置交易
func (txn Transaction) Verify(hashedTxn map[string]Transaction) bool {
	if txn.IsCoinbase() {
		return true
	}

	txnCopy := txn.TrimmedCopy()
	curve := elliptic.P256()

	// 遍历交易的输入
	for inIndex, vin := range txn.Vin {
		prevTxn := hashedTxn[vin.VoutHash.String()]
		txnCopy.Vin[inIndex].PubKey = PublicKey(prevTxn.Vout[vin.VoutIndex].PubKeyHash)
		dataToVerify := []byte(fmt.Sprintf("%s\n", txnCopy))

		// 验证TxnInput的签名是否正确
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, dataToVerify, &r, &s) {
			log.Debugln("ecdsa.Verify Failed", txn, hashedTxn)
			return false
		}
		txnCopy.Vin[inIndex].PubKey = nil
	}

	return true
}

func (txn *Transaction) RelayVerify(merkleRoot HashValue, relayMerklePath []MerklePath) bool {
	hash := txn.Hash()
	var sha [32]byte
	// log.Traceln("RelayVerify", hash, merkleRoot)
	for _, path := range relayMerklePath {
		if path.Left {
			sha = sha256.Sum256(append(path.HashValue, hash...))
		} else {
			sha = sha256.Sum256(append(hash, path.HashValue...))
		}
		hash = sha[:]
		// log.Traceln("RelayVerify", path, hash)
	}
	return merkleRoot.Equal(hash)
}
