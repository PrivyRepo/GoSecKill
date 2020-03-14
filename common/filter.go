package common

import "net/http"

//声明一个新的数据类型(函数类型)
type FilterHandle func(rw http.ResponseWriter, r *http.Request) error

//拦截器结构体
type Filter struct {
	//用来存储需要拦截的URI
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{make(map[string]FilterHandle)}
}

//注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

//根据uri获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

type WebHandle func(rw http.ResponseWriter, r *http.Request)

func (f *Filter) Handle(webhandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				//执行拦截业务逻辑
				e := handle(rw, r)
				if e != nil {
					rw.Write([]byte(e.Error()))
					return
				}
				break
			}
		}
		//执行正常注册的函数
		webhandle(rw, r)
	}

}
