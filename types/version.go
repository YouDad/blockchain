package types

type Version struct {
	Group    int
	Version  int
	Height   int32
	RootHash HashValue
	NowHash  HashValue
}
