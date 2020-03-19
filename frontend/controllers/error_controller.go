package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error501() {
	c.Data["content"] = "server error"
	c.TplName = "shared/error.html"
}
