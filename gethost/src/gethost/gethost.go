package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	HOST_SERVER  = "https://raw.githubusercontent.com/racaljk/hosts/master/hosts"
	//AUTHOR  = "http://www.wzjg520.com"
	LINUX_HOST   = "/etc/hosts"
	WINDOWS_HOST = "c:/Windows/System32/drivers/etc/hosts"
)

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

func getLocalHosts()(host []byte) {
	file, err := os.Open(LINUX_HOST)
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
	file, err := os.OpenFile(LINUX_HOST, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	checkError(err)
	host = append([]byte("\n# go_hosts start\n"), host...)
	host = append(host, []byte("\n# go_hosts end\n")...)
	localHost = append(localHost, host...)
	_, err = file.Write(localHost)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("replace host complete")
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
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
	runtime := end - start
	log.Println("run time ", runtime)
}
