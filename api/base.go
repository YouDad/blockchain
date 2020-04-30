package api

import (
	"encoding/json"
	"os"

	"github.com/YouDad/blockchain/log"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

type SimpleJSONResult struct {
	Message string
	Data    interface{}
}

func (c *BaseController) ParseParameter(data interface{}) {
	log.Debugln(log.Funcname(1), c.GetString("address"))
	if data != nil {
		err := json.Unmarshal(c.Ctx.Input.RequestBody, data)
		if err != nil {
			log.Warnln(string(c.Ctx.Input.RequestBody))
			log.Warnln(err)
			c.ReturnErr(err)
		}
	}
}

func (c *BaseController) GetCookie(key string) string {
	return c.Ctx.GetCookie(key)
}

func (c *BaseController) GetToken() string {
	return c.GetCookie("token")
}

func (c *BaseController) ReturnJson(data SimpleJSONResult) {
	c.Data["json"] = data
	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) ReturnErr(err error) {
	if err != nil {
		c.ReturnJson(SimpleJSONResult{err.Error(), nil})
	}
}

func (c *BaseController) Return(data interface{}) {
	c.ReturnJson(SimpleJSONResult{"", data})
}

func (c *BaseController) Param(key string) string {
	return c.Ctx.Input.Param(key)
}

func ReadJsonFromFile(jsonPath string, dataRef interface{}) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	log.Err(err)

	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	log.Err(err)

	return json.Unmarshal(data, &dataRef)
}
