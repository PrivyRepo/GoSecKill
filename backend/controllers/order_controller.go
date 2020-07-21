package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"homework/models/datamodels"
	"homework/models/services"
	"strconv"
)

type OrderController struct {
	beego.Controller
	OrderService services.IOrderService
}

func (this *OrderController) GetList() {
	uidstr := this.Ctx.GetCookie("uid")
	uid, _ := strconv.Atoi(uidstr)
	s := this.GetString("pagenum")
	var indexpage int
	indexpage, err := strconv.Atoi(s)
	if err != nil || indexpage == 0 {
		indexpage = 1
	}
	arr, count, e := this.OrderService.GetOrderInfoByShop(uid, indexpage, 10)
	fmt.Println(arr)
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
	this.Data["orders"] = arr
	this.Layout = "shared/layout.html"
	this.TplName = "order/view.html"
}


func (this *OrderController) GetUpdate() {
	idstr := this.GetString("id")
	id, _ := strconv.Atoi(idstr)
	order, err := this.OrderService.GetOrderByID(id)
	if err != nil {
		this.Abort("501")
	}
	logs.Info(order)
	order.OrderDeliverStatus = datamodels.DeliverSuccess
	err = this.OrderService.UpdateOrder(order)
	if err != nil {
		logs.Error(err)
		this.Abort("501")
	}
	this.Redirect("/order/list", 301)
}

