package main

import (
	"encoding/json"
	"fmt"
	"homework/common/rabbitmq"
	"homework/seckill/tool/consistent"
	"homework/seckill/tool/ip"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

//设置集群地址,最好内网ip
var (
	hostArray = []string{"192.168.19.101", "192.168.19.102", "192.168.19.103"}
	//hostArray = []string{"127.0.0.1", "127.0.0.1"}
	localHost = ""
	port      = "8083"
)

var hashConsistent *consistent.Consistent
var rabbitMqValidate *rabbitmq.RabbitMQ
var countControl = &CountControl{sourceArray: make(map[int]int)}

//用来存放商品数量
type CountControl struct {
	sourceArray map[int]int
	sync.RWMutex
}

func (m *CountControl) GetCount(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}
func (m *CountControl) SetCount(uid int, count int) bool {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = count
	fmt.Println("商品超卖控制：", m.sourceArray)
	return true
}
func (m *CountControl) IsEnough(uid int) (isOk bool) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	if m.sourceArray[uid] > 0 {
		m.sourceArray[uid]--
		return true
	}
	return false
}
func (m *CountControl) LocalOrder(uidstr string, productstr string) bool {
	userID, _ := strconv.ParseInt(uidstr, 10, 64)
	productID, _ := strconv.ParseInt(productstr, 10, 64)
	if m.IsEnough(int(productID)) {
		//生成订单
		//3.创建消息体
		message := rabbitmq.NewMessage(userID, productID)
		//类型转化
		byteMessage, err := json.Marshal(message)
		if err != nil {
			return false
		}
		//4.生产消息
		err = rabbitMqValidate.PublishSimple(string(byteMessage))
		fmt.Println("生成订单", message)
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}
}

func KillOrder(w http.ResponseWriter, req *http.Request) {
	//获取商品id
	uid, uerr := req.Cookie("uid")
	queryForm, perr := url.ParseQuery(req.URL.RawQuery)
	if uerr != nil || perr != nil || len(queryForm["productid"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	hostRequest, err := hashConsistent.Get(queryForm["productid"][0])
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	var right bool
	if hostRequest == localHost {
		productStr := queryForm["productid"][0]
		right = countControl.LocalOrder(uid.Value, productStr)
	} else {
		distUrl := "http://" + hostRequest + ":" + port + "/getOne?" + req.URL.RawQuery
		right = Proxy(hostRequest, req, distUrl)
	}
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}
func RegisterProduct(w http.ResponseWriter, req *http.Request) {
	//获取商品id
	queryForm, perr := url.ParseQuery(req.URL.RawQuery)
	if perr != nil || len(queryForm["productid"]) <= 0 || len(queryForm["count"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	hostRequest, err := hashConsistent.Get(queryForm["productid"][0])
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	var right bool
	if hostRequest == localHost {
		productidStr := queryForm["productid"][0]
		countStr := queryForm["count"][0]
		productid, e1 := strconv.Atoi(productidStr)
		count, e2 := strconv.Atoi(countStr)
		if e1 != nil || e2 != nil {
			w.Write([]byte("false"))
			return
		}
		right = countControl.SetCount(productid, count)
	} else {
		distUrl := "http://" + hostRequest + ":" + port + "/setCount?" + req.URL.RawQuery
		right = Proxy(hostRequest, req, distUrl)
	}
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return

}

//代理
func Proxy(host string, request *http.Request, url string) bool {
	fmt.Println("走代理", url)
	//获取uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	//获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	response, err := client.Do(req)
	fmt.Println("代理返回", response.Header)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}

	return false
}

func main() {
	//负载均衡器设置
	//采用一次性hash算法
	hashConsistent = consistent.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := ip.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println("获得本机Ip为：", localHost)

	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMqValidate.Destory()
	http.HandleFunc("/setCount", RegisterProduct)
	http.HandleFunc("/getOne", KillOrder)
	err = http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Print("失败")
	}
}
