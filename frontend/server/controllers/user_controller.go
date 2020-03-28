package controllers

import (
	"github.com/astaxie/beego"
	"homework/common/encrypt"
	"homework/models/datamodels"
	"homework/models/services"
	"strconv"
)

type UserController struct {
	beego.Controller
	UserService services.IUserService
}

func (this *UserController) GetRegister() {
	this.TplName = "user/register.html"
}

func (this *UserController) PostRegister() {
	var (
		nickName = this.GetString("nickName")
		userName = this.GetString("userName")
		password = this.GetString("password")
	)
	//ozzo-validation
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}
	_, err := this.UserService.AddUser(user)
	if err != nil {
		this.Abort("500")
	}
	this.Ctx.Redirect(302, "/user/login")
}

func (this *UserController) PostLogin() {
	//1.获取用户提交的表单信息
	var (
		userName = this.GetString("userName")
		password = this.GetString("password")
	)
	//2.验证好账号密码正确
	user, isOk := this.UserService.IsPwdSuccess(userName, password)
	if !isOk {
		this.Data["json"] = "failed"
	}
	//3.写入用户ID到cookie中
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, e := encrypt.EnPwdCode(uidByte)
	if e != nil {
		this.Abort("500")
	}
	this.Ctx.SetCookie("uid", strconv.FormatInt(user.ID, 10), "/")
	this.Ctx.SetCookie("sign", uidString, "/")
	this.Data["json"] = "登录成功"
	this.ServeJSON()
}

func (this *UserController) GetLogout() {
	this.Ctx.SetCookie("uid", "", -1)
	this.Ctx.SetCookie("sign", "", -1)
	this.Ctx.Redirect(302, "/product/list")
}
