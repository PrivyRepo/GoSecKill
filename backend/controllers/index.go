package controllers

import "github.com/astaxie/beego"

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	this.Layout = "shared/layout.html"
	this.TplName = "dashboard/index.html"
}
