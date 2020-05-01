package types

import "github.com/YouDad/blockchain/utils"

type Version struct {
	Group    int
	Version  int
	Height   int32
	RootHash HashValue
	NowHash  HashValue
}

func (v Version) String() string {
	return string(utils.Encode(v))
}
