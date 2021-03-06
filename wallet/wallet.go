package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/YouDad/blockchain/log"
	"github.com/YouDad/blockchain/types"
	"github.com/YouDad/blockchain/utils"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey types.PrivateKey
	PublicKey  types.PublicKey
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{private, public}
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	return address
}

func (w Wallet) String() string {
	return fmt.Sprintf("%s", w.GetAddress())
}

// HashPubKey hashes public key
func HashPubKey(pubKey types.PublicKey) types.HashValue {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	log.Err(err)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))

	if len(pubKeyHash) < 5 {
		return false
	}
	version := pubKeyHash[0:1]
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]

	targetChecksum := checksum(append(version, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func newKeyPair() (types.PrivateKey, types.PublicKey) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	log.Err(err)

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
