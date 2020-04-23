package types

type Signature []byte

func (h Signature) String() string {
	ret := make([]byte, len(h)*2)
	i := 0
	for _, v := range h {
		ret[i] = hexTable[v>>4]
		ret[i+1] = hexTable[v&0x0f]
		i += 2
	}
	return string(ret)
}
