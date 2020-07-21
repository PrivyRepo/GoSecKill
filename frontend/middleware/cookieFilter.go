package middleware

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"homework/common/encrypt"
)

var FilterUser = func(ctx *context.Context) {
	uid, sign := ctx.GetCookie("uid"), ctx.GetCookie("sign")
	if uid == "" || sign == "" {
		logs.Info("校验cookie失败，跳转登录")
		ctx.Redirect(302, "/user/login")
	}
	signbyte, err := encrypt.DePwdCode(sign)
	if err != nil || uid != string(signbyte) {
		logs.Info("校验cookie失败，跳转登录")
		ctx.Redirect(302, "/user/login")
	}
}

func init() {
	beego.InsertFilter("/order/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/product/kill", beego.BeforeRouter, FilterUser)

}
