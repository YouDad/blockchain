package types

import "encoding/hex"

type Signature []byte

func (s Signature) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *Signature) UnmarshalJSON(bytes []byte) error {
	var err error
	*s, err = hex.DecodeString(string(bytes[1 : len(bytes)-1]))
	return err
}

func (s Signature) String() string {
	return hex.EncodeToString(s)
}
