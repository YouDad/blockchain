package types

import (
	"crypto/sha256"
	"fmt"

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

func (pk PublicKey) String() string {
	return fmt.Sprintf("%x", []byte(pk))
}
