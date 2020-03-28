package main

import (
	"errors"
	"fmt"
	"homework/common/encrypt"
	"homework/seckill/tool/consistent"
	"homework/seckill/tool/filter"
	"homework/seckill/tool/ip"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	hostArray = []string{"192.168.19.101", "192.168.19.102", "192.168.19.103"}
	//hostArray = []string{"127.0.0.1", "127.0.0.1"}
	localHost = ""
	port      = "8081"
	killport  = "8083"
)

var hashConsistent *consistent.Consistent

//同一验证拦截器
func Auth(w http.ResponseWriter, r *http.Request) error {
	//添加基于cookie的权限验证
	return CheckUserInfo(r)
}
func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("uid cookie 获取失败")
	}
	signcookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("加密串已被篡改")
	}
	signbyte, err := encrypt.DePwdCode(signcookie.Value)
	if err != nil {
		return errors.New("加密串Cookie获取失败")
	}
	fmt.Println("uid= ", uidCookie.Value, " sign= ", signcookie.Value, " 执行验证")
	if checkInfo(uidCookie.Value, string(signbyte)) {
		fmt.Println("本机身份校验成功")
		return nil
	}
	fmt.Println("本机身份校验失败")
	return errors.New("身份校验失败")
}
func checkInfo(checkStr string, signStr string) bool {

	if checkStr == signStr {
		return true
	}
	return false
}

//用来存放控制信息,创建全局变量,线程安全
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourceArray map[int]interface{}
	sync.RWMutex
}

var accessControl = &AccessControl{sourceArray: make(map[int]interface{})}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	//获取用户uid
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	hostRequest, err := hashConsistent.Get(uid.Value)
	fmt.Println("目标主机:", hostRequest)
	if err != nil {
		return false
	}
	//判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校验
		if m.GetDataFromMap(uid.Value) {
			fmt.Println("筛选成功，尝试生成订单")
			distUrl := "http://" + localHost + ":" + killport + "/getOne?" + req.URL.RawQuery
			return Proxy(hostRequest, req, distUrl)
		} else {
			fmt.Println("筛选成功，此次秒杀无效")
			return false
		}
	} else {
		//不是本机充当代理访问数据放回结果
		distUrl := "http://" + hostRequest + ":" + port + "/access?" + req.URL.RawQuery
		return Proxy(hostRequest, req, distUrl)
	}
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

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = "hello world"
}
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	fmt.Println("本机执行访问")
	//TODO:本机过滤及筛选
	/*
		1. 权限验证
		2. 黑名单机制
		3. 限购一件，布隆过滤器
	*/
	return true
}

func main() {
	//负载均衡器设置
	//采用一次性hash算法
	hashConsistent = consistent.NewConsistent()
	//采用一致性hash算法,添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := ip.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println("获得当机Ip为：", localHost)

	//1. 过滤器
	filter := filter.NewFilter()
	//2. 注册拦截器
	filter.RegisterFilterUri("/access", Auth)
	//3. 启动服务
	http.HandleFunc("/access", filter.Handle(CheckRight))
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Print("shibai")
	}
}
