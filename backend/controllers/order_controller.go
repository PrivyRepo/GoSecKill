package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"homework/common"
	"homework/models/services"
	"strconv"
)

type OrderController struct {
	beego.Controller
	OrderService services.IOrderService
}

func (this *OrderController) GetList() {
	arr, e := this.OrderService.GetAllOrderInfo()
	if common.CheckErr(e) {
		return
	}
	logs.Info(arr)
	this.Data["orders"] = arr
	this.Layout = "shared/layout.html"
	this.TplName = "order/view.html"
}

func (this *OrderController) GetDlete() {
	idstring := this.GetString("id")
	i, _ := strconv.Atoi(idstring)
	isok := this.OrderService.DeleteOrderByID(int64(i))
	if isok {
		this.Ctx.Redirect(302, "/order/list")
	}
}
