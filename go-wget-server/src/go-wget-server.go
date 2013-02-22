package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

const (
	DOWNLOAD_PATH = "/www/domain/download.heroin.so/"
	SERVER        = "GO-WGET-SERVER"
	PORT          = 12321
)

func download(path string) {
	runtime.Gosched()
	log.Printf("download start %s \n", path)
	cmd := exec.Command("wget.exe", "-P", DOWNLOAD_PATH, path)
	err := cmd.Run()
	if err != nil {
		log.Printf("download [error] path=%s, ", path, err)
	}
	runtime.GC()
	log.Printf("download over %s \n", path)
}

func index(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)

	if request.URL.Path != "/" && request.URL.Path != "/favicon.ico" {
		go download(fmt.Sprintf("http:/%s", request.URL.Path))
		fmt.Fprintf(out, "download http:/"+request.URL.Path+"\r\n")
	} else {
		fmt.Fprintf(out, "error\r\n")
	}
}

func remove(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)

	request.ParseForm()
	fmt.Println(request.Form["file"])
	fmt.Println(request.Form["dir"])
	fmt.Fprintf(out, "rm\r\n")
}

func main() {
	runtime.GOMAXPROCS(5)
	http.HandleFunc("/", index)
	http.HandleFunc("/rm", remove)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
