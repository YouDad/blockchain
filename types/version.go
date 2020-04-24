package types

import "fmt"

type Version struct {
	Group    int
	Version  int
	Height   int32
	RootHash HashValue
	NowHash  HashValue
}

func (v Version) String() (ret string) {
	ret += fmt.Sprintf("Version: %d, Group: %d, Height: %d\n", v.Version, v.Group, v.Height)
	ret += fmt.Sprintf("RootHash: %s\n", v.RootHash)
	ret += fmt.Sprintf("NowHash:  %s\n", v.NowHash)
	return ret
}
