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
	"container/list"
	"time"
	"os"
)

const (
	FIND_FILE_REGEX    = `<a href="([^\"\']+?)">`
	FIND_PATH_REGEX    = `<a href="([^\"\']+?)"">`
	FIND_CONTENT_REGEX = `"([^\"\']+?)"`
)

var (
	MAX_NO_TASK_COUNT     = 2000
	CURRENT_NO_TASK_COUNT = MAX_NO_TASK_COUNT
	QUEUE                 = list.New()
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
		go find(path)
	}
}

func download(path string) {
	if CURRENT_NO_TASK_COUNT < MAX_NO_TASK_COUNT {
		CURRENT_NO_TASK_COUNT = MAX_NO_TASK_COUNT
		log.Printf("Reset Loop Count [%d] \n", MAX_NO_TASK_COUNT)
	}
	if strings.Index(path, "SNAPSHOT") < 0 {
		runtime.Gosched()
			log.Printf("download start %s \n", path)
		cmd := exec.Command("wget.exe", "-r", path)
		err := cmd.Run()
		if err != nil {
			log.Printf("download [error] path=%s, ", path, err)
		}
		runtime.GC()
			log.Printf("download over %s \n", path)
	} else {
		log.Printf("this is SNAPSHOT \n")
	}
}

func main() {
	runtime.GOMAXPROCS(12)
	//	find("http://maven.open-ns.org/repo1/org/springframework/")
	//	find("http://maven.open-ns.org/clojars/")
	//	find("http://maven.open-ns.org/repo1/org/apache/cassandra/")
	//	find("http://maven.open-ns.org/repo1/org/apache/commons/")
	//	find("http://maven.open-ns.org/repo1-cache/org/apache/mina/")
	//	find("http://maven.open-ns.org/repo1/org/jboss/netty/")
	//	find("http://maven.open-ns.org/repo/org/apache/wicket/wicket-core/6.4.0/")
		go find("http://maven.open-ns.org/repo1/org/")
//	find("http://maven.open-ns.org/repo/io/netty/netty-all/")
	for {
		time.Sleep(time.Duration(3000)*time.Millisecond)

		if QUEUE.Len() > 0 {
			if CURRENT_NO_TASK_COUNT < MAX_NO_TASK_COUNT {
				CURRENT_NO_TASK_COUNT = MAX_NO_TASK_COUNT
				log.Printf("Reset Loop Count [%d] \n", MAX_NO_TASK_COUNT)
			}
			task := QUEUE.Back()
			QUEUE.Remove(task)
			go download(fmt.Sprintf("%s", task.Value))
		} else {
			log.Printf("Not Task Close By Loop [%d] count , Sleep 1s \n", CURRENT_NO_TASK_COUNT)
			CURRENT_NO_TASK_COUNT--
			time.Sleep(time.Duration(3)*time.Second)
			if CURRENT_NO_TASK_COUNT == 0 {
				log.Printf("Close Application Bye!!! \n")
				os.Exit(0)
			}
		}
	}
}
