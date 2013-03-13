package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

const (
	FIND_FILE_REGEX    = `<a href="([^\"\']+?)">`
	FIND_PATH_REGEX    = `<a href="([^\"\']+?)"">`
	FIND_CONTENT_REGEX = `"([^\"\']+?)"`
)

func find(url string) {
	response, _ := http.Get(url)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	re_file, _ := regexp.Compile(FIND_FILE_REGEX)
	re_path, _ := regexp.Compile(FIND_PATH_REGEX)
	re_content, _ := regexp.Compile(FIND_CONTENT_REGEX)

	result_file := re_file.FindAll(body, -1)
	result_path := re_path.FindAll(body, -1)

	for i := range result_file {
		file := string(re_content.Find(result_file[i]))
		file = strings.Replace(file, "\"", "", -1)
		if file != "../" {
			file = fmt.Sprintf("%s%s", url, file)
			go download(file)
		}
	}

	for i := range result_path {
		path := string(re_content.Find(result_path[i]))
		path = strings.Replace(path, "\"", "", -1)
		path = fmt.Sprintf("%s%s", url, path)
		find(path)
	}
}

func download(path string) {
	runtime.Gosched()
	log.Printf("download start %s \n", path)
	cmd := exec.Command("wget.exe", "-r", path)
	err := cmd.Run()
	if err != nil {
		log.Printf("download [error] path=%s, ", path, err)
	}
	runtime.GC()
	log.Printf("download over %s \n", path)
}

func main() {
	runtime.GOMAXPROCS(5)
	find("path")
}
