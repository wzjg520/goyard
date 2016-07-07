package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	testString := "Hi guihua"
	//	md5Int := md5.New()
	//	md5Int.Write([]byte(testString))
	//	result := md5Int.Sum([]byte(""))
	//	fmt.Println(result)

	fmt.Println(md5.Sum([]byte(testString)))
}
