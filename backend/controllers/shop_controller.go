package controllers

import (
	"github.com/astaxie/beego"
	"homework/common"
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
	//ozzo-validation
	shop := &datamodels.Shop{
		UserName:     userName,
		ShopName:     shopName,
		HashPassword: password,
	}
	_, err := this.ShopService.AddShop(shop)
	if common.CheckErr(err) {
		this.Abort("401")

	}
	this.Ctx.Redirect(302, "/shop/login")
}

func (this *ShopController) GetLogin() {
	this.TplName = "shop/login.html"
}

func (this *ShopController) PostLogin() {
	//1.获取用户提交的表单信息
	var (
		userName = this.GetString("userName")
		password = this.GetString("password")
	)
	//2.验证好账号密码正确
	user, isOk := this.ShopService.IsPwdSuccess(userName, password)
	if !isOk {
		this.Layout = "shop/login.html"
		this.Abort("401")
	}
	//3.写入用户ID到cookie中
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, e := common.EnPwdCode(uidByte)
	if common.CheckErr(e) {
		this.Abort("401")

	}
	this.Ctx.SetCookie("uid", strconv.FormatInt(user.ID, 10), "/")
	this.Ctx.SetCookie("sign", uidString, "/")
	this.Ctx.Redirect(302, "/product/list")
}

func (this *ShopController) GetLogout() {
	this.Ctx.SetCookie("uid", "", -1)
	this.Ctx.SetCookie("sign", "", -1)
	this.Redirect("/shop/login", 302)
}
