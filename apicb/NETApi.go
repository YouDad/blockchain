package apicb

import "github.com/YouDad/blockchain/core"

type NETApi struct {
	set *core.UTXOSet
}

func GetNETApi() *NETApi {
	return &NETApi{core.GetUTXOSet()}
}

func (net *NETApi) HeartBeat(args *NIL, reply *NIL) error {
	return nil
}
