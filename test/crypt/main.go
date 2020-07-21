package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/crc32"
)

func MD5(str string)string{
	c := md5.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

func SHA1(str string)string{
	c:=sha1.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

func CRC32(str string)uint32{
	return crc32.ChecksumIEEE([]byte(str))
}

func main() {
	fmt.Println(CRC32("12345"))
	fmt.Println(MD5("12345"))
	fmt.Println(SHA1("12345"))
}