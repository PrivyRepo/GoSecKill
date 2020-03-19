package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	_ "homework/backend/routers"
)

var FilterUser = func(ctx *context.Context) {
	cookie := ctx.GetCookie("uid")
	logs.Info(cookie)
	if ctx.Request.RequestURI != "/shop/login" && ctx.Request.RequestURI != "/shop/register" && cookie == "" {
		ctx.Redirect(302, "/shop/login")
	}
}

func main() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
	beego.SetStaticPath("/assets", "assets")
	beego.Run()
}
