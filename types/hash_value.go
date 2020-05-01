package types

import (
	"bytes"
	"encoding/hex"
)

type HashValue []byte

func (h HashValue) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.String() + `"`), nil
}

func (h *HashValue) UnmarshalJSON(bytes []byte) error {
	var err error
	*h, err = hex.DecodeString(string(bytes[1 : len(bytes)-1]))
	return err
}

func (h HashValue) String() string {
	return hex.EncodeToString(h)
}

func (h HashValue) Equal(other HashValue) bool {
	return bytes.Compare(h, other) == 0
}

func (h HashValue) Key() [32]byte {
	var ret [32]byte
	for i := range h {
		if i >= 32 {
			break
		}
		ret[i] = h[i]
	}
	return ret
}
