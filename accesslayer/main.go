package main

import (
	"errors"
	"fmt"
	"homework/common"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//设置集群地址,最好内网ip
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	//获取用户uid
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	//判断是否为本机
	if hostRequest == localHost {
		//执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	} else {
		//不是本机充当代理访问数据放回结果
		return GetDataFromOtherMap(hostRequest, req)
	}
}

//获取其他节点处理结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
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
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/access", nil)
	if err != nil {
		return false
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	response, err := client.Do(req)
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

//获取本机map,并且处理业务逻辑,返回的结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)
	//执行逻辑判断
	if data != nil {
		return true
	}
	return
}

//用来存放控制信息
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourceArray map[int]interface{}
	sync.RWMutex
}

//创建全局变量
var accessControl = &AccessControl{sourceArray: make(map[int]interface{})}

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

//同一验证拦截器
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证")
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
	signbyte, err := common.DePwdCode(signcookie.Value)
	if err != nil {
		return errors.New("加密串Cookie获取失败")
	}
	log.Println(uidCookie.Value, string(signbyte))
	if checkInfo(uidCookie.Value, string(signbyte)) {
		return nil
	}
	return errors.New("身份校验失败")

}

//自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool {

	if checkStr == signStr {
		fmt.Println(checkStr, signStr)
		return true
	}
	return false
}

func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	w.Write([]byte("验证成功"))
}

func main() {
	//负载均衡器设置
	//采用一次性hash算法
	hashConsistent = common.NewConsistent()
	//采用一致性hash算法,添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	//1. 过滤器
	filter := common.NewFilter()
	//2. 注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	//3. 启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Print("shibai")
	}
}
