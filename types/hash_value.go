package types

import "bytes"

type HashValue []byte

const hexTable = "0123456789abcdef"

func (h HashValue) String() string {
	ret := make([]byte, len(h)*2)
	i := 0
	for _, v := range h {
		ret[i] = hexTable[v>>4]
		ret[i+1] = hexTable[v&0x0f]
		i += 2
	}
	return string(ret)
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
