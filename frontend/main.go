package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"homework/frontend/controllers"
	_ "homework/frontend/routers"
)

var FilterUser = func(ctx *context.Context) {
	cookie := ctx.GetCookie("uid")
	logs.Info(cookie)
	if ctx.Request.RequestURI != "/user/login" && ctx.Request.RequestURI != "/user/register" && cookie == "" {
		ctx.Redirect(302, "/user/login")
	}
}

func main() {
	beego.InsertFilter("/product/order", beego.BeforeRouter, FilterUser)
	beego.ErrorController(&controllers.ErrorController{})
	beego.SetStaticPath("/public", "public")
	beego.Run()
}
