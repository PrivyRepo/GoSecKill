package middleware

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"homework/common/encrypt"
)

var FilterUser = func(ctx *context.Context) {
	uid, sign := ctx.GetCookie("uid"), ctx.GetCookie("sign")
	if ctx.Request.RequestURI == "/shop/login" || ctx.Request.RequestURI == "/shop/register" {
		logs.Info("注册登录")
	} else {
		if uid == "" || sign == "" {
			logs.Info("校验cookie失败，跳转登录")
			ctx.Redirect(302, "/shop/login")
		}
		signbyte, err := encrypt.DePwdCode(sign)
		if err != nil || uid != string(signbyte) {
			logs.Info("校验cookie失败，跳转登录")
			ctx.Redirect(302, "/shop/login")
		}
	}
}

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}
