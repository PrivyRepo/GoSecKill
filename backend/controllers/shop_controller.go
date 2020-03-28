package controllers

import (
	"github.com/astaxie/beego"
	"homework/common/encrypt"
	"homework/models/datamodels"
	"homework/models/services"
	"strconv"
)

type ShopController struct {
	beego.Controller
	ShopService services.IShopService
}

func (this *ShopController) GetRegister() {
	this.TplName = "shop/register.html"
}

func (this *ShopController) PostRegister() {
	var (
		shopName = this.GetString("shopName")
		userName = this.GetString("userName")
		password = this.GetString("password")
	)
	shop := &datamodels.Shop{
		UserName:     userName,
		ShopName:     shopName,
		HashPassword: password,
	}
	_, err := this.ShopService.AddShop(shop)
	if err != nil {
		this.Abort("500")
	}
	this.Redirect("/shop/login", 302)
}

func (this *ShopController) GetLogin() {
	this.TplName = "shop/login.html"
}

func (this *ShopController) PostLogin() {
	var (
		userName = this.GetString("userName")
		password = this.GetString("password")
	)
	user, isOk := this.ShopService.IsPwdSuccess(userName, password)
	if !isOk {
		this.Abort("401")
	}
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, e := encrypt.EnPwdCode(uidByte)
	if e != nil {
		this.Abort("501")
	}
	this.Ctx.SetCookie("uid", strconv.FormatInt(user.ID, 10), "/")
	this.Ctx.SetCookie("sign", uidString, "/")
	this.Redirect("/", 302)
}

func (this *ShopController) GetLogout() {
	this.Ctx.SetCookie("uid", "", -1)
	this.Ctx.SetCookie("sign", "", -1)
	this.Redirect("/shop/login", 302)
}
