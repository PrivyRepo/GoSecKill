package config

import (
	"bufio"
	"os"
)

func ReadLineFile(fileName string) []string{
	res := make([]string,0)
	if file, err := os.Open(fileName);err !=nil{
		panic(err)
	}else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan(){
			res = append(res,scanner.Text())
		}
	}
	return res
}
