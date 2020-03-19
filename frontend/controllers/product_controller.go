package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"homework/common"
	"homework/models/datamodels"
	"homework/models/services"
	"homework/rabbitmq"
	"strconv"
)

type ProductController struct {
	beego.Controller
	OrderService   services.IOrderService
	ProductService services.IProductService
	RabbitMQ       *rabbitmq.RabbitMQ
}

func (this *ProductController) GetList() {
	products, e := this.ProductService.GetAllProduct()
	if common.CheckErr(e) {
		return
	}
	cookie := this.Ctx.GetCookie("uid")
	var IsLogin bool
	if cookie != "" {
		IsLogin = true
	}
	this.Data["IsLogin"] = IsLogin
	this.Data["products"] = products
	this.Layout = "shared/productLayout.html"
	this.TplName = "product/listview.html"
}

func (this *ProductController) GetDetail() {
	idstring := this.GetString("id")
	id, _ := strconv.Atoi(idstring)
	ok := this.ProductService.IncPorductReview(int64(id))
	if ok {
		logs.Info("加一成功")
	} else {
		logs.Error("加一失败")
	}

	product, err := this.ProductService.GetProductByID(int64(id))
	if common.CheckErr(err) {
		return
	}
	cookie := this.Ctx.GetCookie("uid")
	var IsLogin bool
	if cookie != "" {
		IsLogin = true
	}
	this.Data["IsLogin"] = IsLogin
	this.Data["product"] = product
	this.Layout = "shared/productLayout.html"
	this.TplName = "product/view.html"

}

func (this *ProductController) GetOrder() {

	productstring := this.GetString("productID")
	userstring := this.Ctx.GetCookie("uid")
	productID, e := strconv.Atoi(productstring)
	if common.CheckErr(e) {
		this.Abort("401")
	}
	userID, e := strconv.Atoi(userstring)
	if common.CheckErr(e) {
		this.Abort("401")

	}

	//创建消息体
	message := datamodels.NewMessage(int64(userID), int64(productID))
	bytes, e := json.Marshal(message)
	if common.CheckErr(e) {
		this.Abort("401")
	}

	e = this.RabbitMQ.PublishSimple(string(bytes))
	if common.CheckErr(e) {
		this.Abort("401")
	}
	this.Ctx.WriteString("successful")
	/*
		product, e := this.ProductService.GetProductByID(int64(productID))
		if common.CheckErr(e) {
			return
		}
		var orderID int64
		var showMessage string = "抢购失败"
		//判断商品数量是否满足需求
		if product.ProductNum > 0 {
			//扣除商品数量
			product.ProductNum -= 1
			err := this.ProductService.UpdateProduct(product)
			if common.CheckErr(err) {
				return
			}
			//创建订单
			userId, err := strconv.Atoi(userstring)
			if common.CheckErr(err) {
				return
			}
			order := &datamodels.Order{
				UserId:      int64(userId),
				ProductId:   int64(productID),
				OrderStatus: datamodels.OrderWait,
			}
			orderID, err = this.OrderService.InsertOrder(order)
			if common.CheckErr(err) {
				return
			} else {
				showMessage = "抢购成功"
			}
		}
		this.Data["orderID"] = orderID
		this.Data["showMessage"] = showMessage
		this.Layout = "shared/productLayout.html"
		this.TplName = "product/result.html"
	*/
}
