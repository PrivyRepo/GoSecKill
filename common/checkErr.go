package common

import "github.com/astaxie/beego/logs"

func CheckErr(err error) bool {
	if err != nil {
		logs.Error(err)
		return true
	}
	return false
}
