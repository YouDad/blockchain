package routers

import (
	"github.com/YouDad/blockchain/api"
	"github.com/astaxie/beego"
)

func init() {
	beego.AddNamespace(beego.NewNamespace("/v1",
		beego.NSNamespace("/db",
			beego.NSRouter("/GetGenesis", new(api.DBController), "post:GetGenesis"),
			beego.NSRouter("/GetBalance", new(api.DBController), "post:GetBalance"),
			beego.NSRouter("/GetBlocks", new(api.DBController), "post:GetBlocks"),
			beego.NSRouter("/GossipTxn", new(api.DBController), "post:GossipTxn"),
			beego.NSRouter("/GossipBlock", new(api.DBController), "post:GossipBlock"),
			beego.NSRouter("/GetHash", new(api.DBController), "post:GetHash"),
		),
		beego.NSNamespace("/version",
			beego.NSRouter("/SendVersion", new(api.VersionController), "post:SendVersion"),
		),
		beego.NSNamespace("/net",
			beego.NSRouter("/HeartBeat", new(api.NetController), "post:HeartBeat"),
			beego.NSRouter("/GetKnownNodes", new(api.NetController), "post:GetKnownNodes"),
		),
		beego.NSNamespace("/server",
			beego.NSRouter("/SendCMD", new(api.ServerController), "post:SendCMD"),
		),
	))
}
