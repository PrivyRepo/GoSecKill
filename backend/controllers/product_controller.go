package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"homework/common"
	"homework/models/datamodels"
	"homework/models/services"
	"strconv"
)

type ProductController struct {
	beego.Controller
	ProductService services.IProductService
}

//列出所有商品
func (this *ProductController) GetAll() {
	cookiestr := this.Ctx.GetCookie("uid")
	if cookiestr == "" {
		this.Redirect("/shop/login", 302)
	}
	shopid, _ := strconv.ParseInt(cookiestr, 10, 64)
	products, e := this.ProductService.GetProductByshop(shopid)
	logs.Info(cookiestr, products)
	if common.CheckErr(e) {
		this.Abort("401")
	}
	this.Data["shopName"] = "test"
	this.Data["productArray"] = products
	this.Layout = "shared/layout.html"
	this.TplName = "product/view.html"
}

//添加商品界面
func (this *ProductController) GetInsert() {
	cookiestr := this.Ctx.GetCookie("uid")
	this.Data["shopid"] = cookiestr
	this.Layout = "shared/layout.html"
	this.TplName = "product/add.html"
}

//添加商品
func (this *ProductController) PostInsert() {
	product := datamodels.Product{}
	if err := this.ParseForm(&product); err != nil {
		logs.Error(err)
		this.Abort("401")
	}
	logs.Info(product)
	_, e := this.ProductService.InsertProduct(&product)
	if common.CheckErr(e) {
		this.Abort("401")

	}
	this.Ctx.Redirect(302, "/product/list")
}

//编辑商品界面
func (this *ProductController) GetManager() {
	idString := this.GetString("id")
	i, err := strconv.ParseInt(idString, 10, 16)
	if common.CheckErr(err) {
		this.Abort("401")

	}
	product, err := this.ProductService.GetProductByID(i)
	if common.CheckErr(err) {
		this.Abort("401")

	}
	this.Data["product"] = product
	this.Layout = "shared/layout.html"
	this.TplName = "product/manager.html"
}

//编辑商品
func (this *ProductController) PostManager() {
	product := datamodels.Product{}
	if err := this.ParseForm(&product); err != nil {
		logs.Info(this.Ctx.Request)
		logs.Error(err)
		this.Abort("401")
	}
	logs.Info(product)
	err := this.ProductService.UpdateProduct(&product)
	if common.CheckErr(err) {
		this.Abort("401")
	}
	this.Ctx.Redirect(302, "/product/list")
}

//删除商品
func (this *ProductController) GetDelete() {
	idString := this.GetString("id")
	i, err := strconv.ParseInt(idString, 10, 16)
	if common.CheckErr(err) {
		this.Abort("401")
	}
	isok := this.ProductService.DeleteProductByID(i)
	if isok {
		logs.Info("删除商品成功，id为", i)
	} else {
		logs.Info("删除商品失败，id为", i)
	}
	this.Ctx.Redirect(302, "/product/list")
}
