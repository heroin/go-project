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
	"time"
)

const (
	FIND_PROJECT_REGEX = `<a href="([^\"\']+?)">`
	FIND_CONTENT_REGEX = `"([^\"\']+?)"`
	GET_URL            = "http://192.168.192.81/"
	OPT_URL            = "http://ci.las.360buy.net/job/%s/build?delay=0sec"
)

func auto(url string) {
	re_project, _ := regexp.Compile(FIND_PROJECT_REGEX)
	re_content, _ := regexp.Compile(FIND_CONTENT_REGEX)
	result := re_project.FindAll(get(url), -1)

	for i := range result {
		project := string(re_content.Find(result[i]))
		project = strings.Replace(project, "\"", "", -1)
		if project != "../" {
			get_url := fmt.Sprintf("%s%s", GET_URL, project)
			if strings.Replace(string(get(get_url)), "\n", "", -1) == "1" {
				log.Printf("build [info] project=%s, ", project)
				opt_url := fmt.Sprintf(OPT_URL, project)
				get(opt_url)
				go reset(project)
			} else {
				log.Printf("auto [info] project=%s, ", project)
			}
		}
	}
}

func get(url string) []byte {
	response, _ := http.Get(url)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return body
}

func reset(project string) {
	runtime.Gosched()
	cmd := exec.Command("reset.bat", project)
	err := cmd.Run()
	if err != nil {
		log.Printf("reset [error] project=%s, ", project, err)
	}
	runtime.GC()
}

func main() {
	for {
		auto(GET_URL)
		time.Sleep(time.Second * 15)
	}
}
