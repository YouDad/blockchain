package types

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/YouDad/blockchain/log"
	"golang.org/x/crypto/ripemd160"
)

type PublicKey []byte

func (pk PublicKey) Hash() HashValue {
	publicSHA256 := sha256.Sum256(pk)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	log.Err(err)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func (pk PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pk.String() + `"`), nil
}

func (pk *PublicKey) UnmarshalJSON(bytes []byte) error {
	var err error
	*pk, err = hex.DecodeString(string(bytes[1 : len(bytes)-1]))
	return err
}

func (pk PublicKey) String() string {
	return hex.EncodeToString(pk)
}
