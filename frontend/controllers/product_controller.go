package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"homework/common/encrypt"
	"homework/models/datamodels"
	"homework/models/services"
	"strconv"
)

type ProductController struct {
	beego.Controller
	ProductService services.IProductService
}

type ProductDetail struct {
	IsLogin bool
	Product *datamodels.Product
}

type ProductList struct {
	ProductInfo map[int]map[string]string
	PageInfo    *PageInfo
	IsLogin     bool
}

type PageInfo struct {
	Count      int
	IndexPage  int
	NextPage   int
	PrePage    int
	ArrayPages []int
	EndPage    int
}

func (this *ProductController) GetList() {
	s := this.GetString("pagenum")
	var indexpage int
	indexpage, err := strconv.Atoi(s)
	if err != nil || indexpage == 0 {
		indexpage = 1
	}
	arr, count, e := this.ProductService.GetAllProductInfo(indexpage, 12)
	logs.Info(arr)
	endpage := count/12 + 1
	if e != nil {
		logs.Error(e)
		this.Abort("501")
	}
	arrpages := []int{(indexpage-1)/5*5 + 1}
	for i := 1; i < 5; i++ {
		if arrpages[i-1]+1 > endpage {
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

	uid := this.Ctx.GetCookie("uid")
	sign := this.Ctx.GetCookie("sign")
	signbytes, err := encrypt.DePwdCode(sign)
	var IsLogin bool
	fmt.Println(uid, string(signbytes))
	if err == nil && uid == string(signbytes) {
		IsLogin = true
	} else {
		this.Ctx.SetCookie("uid", "", -1)
		this.Ctx.SetCookie("sign", "", -1)
	}
	//msg := ProductList{
	//	ProductInfo: arr,
	//	PageInfo:    pageinfo,
	//	IsLogin:     IsLogin,
	//}
	//this.Data["json"] = msg
	//this.ServeJSON()
	this.Data["IsLogin"] = IsLogin
	this.Data["products"] = arr
	this.TplName = "product/listview.html"
}

func (this *ProductController) GetDetail() {
	idstring := this.GetString("id")
	id, _ := strconv.Atoi(idstring)
	product, err := this.ProductService.GetProductByID(int64(id))
	if err != nil {
		this.Abort("500")
	}
	uid := this.Ctx.GetCookie("uid")
	sign := this.Ctx.GetCookie("sign")
	signbytes, err := encrypt.DePwdCode(sign)
	var IsLogin bool
	if err == nil && uid == string(signbytes) {
		IsLogin = true
	} else {
		this.Ctx.SetCookie("uid", "", -1)
		this.Ctx.SetCookie("sign", "", -1)
	}
	this.Data["IsLogin"] = IsLogin
	this.Data["product"] = product
	//this.Layout = "shared/productLayout.html"
	this.TplName = "product/view.html"
}


func (this *ProductController) GetTestDetail() {
	idstring := this.GetString("id")
	id, _ := strconv.Atoi(idstring)
	product, err := this.ProductService.GetProductByID(int64(id))
	if err != nil {
		this.Abort("500")
	}
	uid := this.Ctx.GetCookie("uid")
	sign := this.Ctx.GetCookie("sign")
	signbytes, err := encrypt.DePwdCode(sign)
	var IsLogin bool
	if err == nil && uid == string(signbytes) {
		IsLogin = true
	} else {
		this.Ctx.SetCookie("uid", "", -1)
		this.Ctx.SetCookie("sign", "", -1)
	}
	this.Data["IsLogin"] = IsLogin
	this.Data["product"] = product
	//this.Layout = "shared/productLayout.html"
	this.TplName = "product/view_test.html"
}
