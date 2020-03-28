package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"homework/models/services"
	"strconv"
)

type OrderController struct {
	beego.Controller
	OrderService services.IOrderService
}

func (this *OrderController) OrderList() {
	uidstr := this.Ctx.GetCookie("uid")
	uid, _ := strconv.Atoi(uidstr)
	s := this.GetString("pagenum")
	var indexpage int
	indexpage, err := strconv.Atoi(s)
	if err != nil || indexpage == 0 {
		indexpage = 1
	}
	arr, count, e := this.OrderService.GetOrderInfoByUser(uid, indexpage, 10)
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
	this.Data["json"] = arr
	this.ServeJSON()
}
