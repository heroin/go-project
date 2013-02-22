package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	SERVER    = "GO-UPLOAD-SERVER"
	PORT      = 12322
	DOWNLOAD  = "http://127.0.0.1/%s"
	LOCAL_URL = "/www/%s"
	INDEX = `<!DOCTYPE html><html><head><title>UPLOAD</title></head><body><form action="/upload" method="post" enctype ="multipart/form-data"><input name="file" type="file"/><button>upload</button></form></body></html>`
	UPLOAD = `<!DOCTYPE html><html><head><title>UPLOAD</title></head><body><a href="%s">%s</a><br/><a href="http://192.168.192.81">Home</a></body></html>`
)

func index(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	fmt.Fprintf(out, INDEX)
}

func upload(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)

	upload, head, _ := request.FormFile("file")
	defer upload.Close()

	tmp, _ := os.Create(fmt.Sprintf(LOCAL_URL, head.Filename))
	defer tmp.Close()
	_, _ = io.Copy(tmp, upload)
	download := fmt.Sprintf(DOWNLOAD, head.Filename)
	fmt.Fprintf(out, fmt.Sprintf(UPLOAD, download, download))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
