package types

type Version struct {
	Version  int
	Height   int32
	RootHash HashValue
	NowHash  HashValue
}
