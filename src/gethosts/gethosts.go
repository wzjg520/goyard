package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"
)

const (
	HOST_SERVER = "http://www.wzjg520.com/hosts.php"
	//AUTHOR  = "http://www.wzjg520.com"
	LINUX_HOST   = "/etc/hosts"
	WINDOWS_HOST = "c:/Windows/System32/drivers/etc/hosts"
)

var (
	runtimePlatform string
	linuxNewLine    = []byte("\n")
	windosNewLine   = []byte("\r\n")
	currNewLine     []byte
	startHosts      = "#----------MODIFY HOSTS START----------#"
	endHosts        = "#----------MODIFY HOSTS END----------#"
)

func getHost() (hostByte []byte, err error) {
	log.Println("connect to ", HOST_SERVER)
	resp, err := http.Get(HOST_SERVER)

	if resp.StatusCode != http.StatusOK {
		log.Fatal("cant not download", HOST_SERVER)
	}

	checkError(err)
	hostByte, err = ioutil.ReadAll(resp.Body)
	checkError(err)

	re, err := regexp.Compile(`\n$`)
	checkError(err)
	hostByte = re.ReplaceAll(hostByte, currNewLine)

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

	if m, _ := regexp.Match(startHosts+`[\S\s]+?`+endHosts, localHost); m {
		re, _ := regexp.Compile(startHosts + `[\S\s]+?` + endHosts)
		localHost = re.ReplaceAll(localHost, []byte(""))
		log.Println("start replace hosts")
	}
	return localHost

}
func replaceOsHosts(host, localHost []byte) {
	file, err := os.OpenFile(runtimePlatform, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	checkError(err)
	startHosts := append([]byte(startHosts), currNewLine...)
	host = append(startHosts, host...)
	host = append(host, []byte([]byte(endHosts))...)
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
	log.Println("platform: ", platFrom)
	if platFrom == "linux" {
		runtimePlatform = LINUX_HOST
		currNewLine = linuxNewLine
	} else if platFrom == "windows" {
		runtimePlatform = WINDOWS_HOST
		currNewLine = windosNewLine
	}
}

func main() {
	start := time.Now().Unix()
	localHosts := getLocalHosts()
	host, err := getHost()
	if err != nil {
		log.Fatalln(err.Error())
	}
	replaceOsHosts(host, localHosts)
	end := time.Now().Unix()
	time.Sleep(time.Second * 3)
	runtime := end - start
	log.Println("run time ", runtime)
}
