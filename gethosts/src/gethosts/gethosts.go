package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"runtime"
)

const (
	HOST_SERVER = "https://raw.githubusercontent.com/racaljk/hosts/master/hosts"
	//AUTHOR  = "http://www.wzjg520.com"
	LINUX_HOST   = "/etc/hosts"
	WINDOWS_HOST = "c:/Windows/System32/drivers/etc/hosts"
)

var runtimePlatform string;

func getHost() (hostByte []byte, err error) {
	log.Println("connect to ", HOST_SERVER)
	resp, err := http.Get(HOST_SERVER)
	checkError(err)
	hostByte, err = ioutil.ReadAll(resp.Body)
	checkError(err)

	defer func() {
		resp.Body.Close()
		log.Println("get ", HOST_SERVER, "complete")
	}()

	return hostByte, err
}

func getLocalHosts() (host []byte) {
	file, err := os.Open(runtimePlatform)
	defer file.Close()
	checkError(err)
	buf := make([]byte, 1)
	localHost := make([]byte, 1)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		localHost = append(localHost, buf...)
	}

	if m, _ := regexp.Match(`#\sgo_hosts\sstart[\S\s]+?#\sgo_hosts\send`, localHost); m {
		re, _ := regexp.Compile(`#\sgo_hosts\sstart[\S\s]+?#\sgo_hosts\send`)
		localHost = re.ReplaceAll(localHost, []byte(""))
		log.Println("start replace hosts")
	}
	return localHost

}
func replaceOsHosts(host, localHost []byte) {
	file, err := os.OpenFile(runtimePlatform, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	checkError(err)
	host = append([]byte("\n# go_hosts start\n"), host...)
	host = append(host, []byte("\n# go_hosts end\n")...)
	localHost = append(localHost, host...)
	localHost = bytes.Trim(localHost, "\x00")
	_, err = file.Write(localHost)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("replace hosts success")
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func init() {
	platFrom := runtime.GOOS
	log.Println("platform: ",platFrom)
	if platFrom == "linux" {
		runtimePlatform = LINUX_HOST
	} else if platFrom == "windows" {
		runtimePlatform = WINDOWS_HOST
	}
}

func main() {
	start := time.Now().Unix()
	host, err := getHost()
	if err != nil {
		log.Fatalln(err.Error())
	}
	localHosts := getLocalHosts()
	replaceOsHosts(host, localHosts)
	end := time.Now().Unix()
	time.Sleep(time.Second * 3)
	runtime := end - start
	log.Println("run time ", runtime)
}
