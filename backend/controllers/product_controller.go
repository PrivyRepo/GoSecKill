package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
	uidstr := this.Ctx.GetCookie("uid")
	shopid, _ := strconv.Atoi(uidstr)
	s := this.GetString("pagenum")
	var indexpage int
	indexpage, err := strconv.Atoi(s)
	if err != nil || indexpage == 0 {
		indexpage = 1
	}
	arr, count, e := this.ProductService.GetProductByshop(int64(shopid), indexpage, 10)
	endpage := count/10 + 1
	if e != nil {
		logs.Error(e)
		this.Abort("501")
	}
	arrpages := []int{(indexpage-1)/5*5 + 1}
	for i := 1; i < 5; i++ {
		if arrpages[i-1]+1 >= endpage {
			break
		}
		arrpages = append(arrpages, arrpages[i-1]+1)
	}
	//count
	//count/10+1 总页数
	//(pagenum-1)*10,10

	//start:= (pagenum/5)*5	end := start+5>endpage ? start+5:endpage
	this.Data["Count"] = count
	this.Data["IndexPage"] = indexpage
	this.Data["NextPage"] = indexpage + 1
	this.Data["PrePage"] = indexpage - 1
	this.Data["ArrayPages"] = arrpages
	this.Data["EndPage"] = endpage

	this.Data["productArray"] = arr
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
		this.Abort("501")
	}
	logs.Info(product)
	_, e := this.ProductService.InsertProduct(&product)
	if e != nil {
		this.Abort("501")
	}
	this.Ctx.Redirect(302, "/product/list")
}

//编辑商品界面
func (this *ProductController) GetManager() {
	idString := this.GetString("id")
	i, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		this.Abort("501")

	}
	product, err := this.ProductService.GetProductByID(i)
	if err != nil {
		this.Abort("501")

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
		this.Abort("501")
	}
	logs.Info(product)
	err := this.ProductService.UpdateProduct(&product)
	if err != nil {
		this.Abort("501")
	}
	this.Ctx.Redirect(302, "/product/list")
}

//删除商品
func (this *ProductController) GetDelete() {
	idString := this.GetString("id")
	i, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
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
