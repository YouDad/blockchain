package app

type App interface {
	HashPart() []byte
	ToString() string
	GobEncode() ([]byte, error)
	GobDecode([]byte) error
}
