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
			beego.NSRouter("/SendTransaction", new(api.DBController), "post:SendTransaction"),
			beego.NSRouter("/SendBlock", new(api.DBController), "post:SendBlock"),
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
