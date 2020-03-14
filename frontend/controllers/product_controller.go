package controllers

import (
	"github.com/astaxie/beego"
	"homework/common"
	"homework/models/datamodels"
	"homework/models/services"
	"strconv"
)

type ProductController struct {
	beego.Controller
	OrderService   services.IOrderService
	ProductService services.IProductService
}

func (this *ProductController) GetList() {
	products, e := this.ProductService.GetAllProduct()
	if common.CheckErr(e) {
		return
	}
	this.Data["products"] = products
	this.TplName = "product/listview.html"
}

func (this *ProductController) GetDetail() {
	idstring := this.GetString("id")
	id, _ := strconv.Atoi(idstring)
	product, err := this.ProductService.GetProductByID(int64(id))
	if common.CheckErr(err) {
		return
	}
	this.Data["product"] = product
	this.Layout = "shared/productLayout.html"
	this.TplName = "product/view.html"

}

func (this *ProductController) GetOrder() {

	productstring := this.GetString("productID")
	userstring := this.Ctx.GetCookie("uid")
	productID, e := strconv.Atoi(productstring)
	if common.CheckErr(e) {
		return
	}
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

}
