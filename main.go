package main

import (
	database "confluence-proxy-attachment/config"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ConfluenceContent struct {
	CONTENTID string
	PAGEID    string
	SPACEID   string
	TITLE     string
}

type Config struct {
	URLPatterns   []Info   `yaml:"urlPatterns"`
	SourceBegin   string   `yaml:"sourceBegin"`
	SecurityChain []string `yaml:"securityChain"`
}

type Info struct {
	Info URLPattern `yaml:"info"`
}

type URLPattern struct {
	Pattern string `yaml:"pattern"`
	Fields  Field  `yaml:"fields"`
}

type Field struct {
	Type      string `yaml:"type"`
	PageId    int    `yaml:"pageId"`
	ContentId int    `yaml:"contentId"`
}

var configPath *string

func main() {
	port := flag.Int("port", 8080, "Port number to listen on")
	configPath = flag.String("config", "config.yaml", "配置文件地址")
	flag.Parse()

	// 初始化数据库连接
	err := database.InitDB(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	// 定义处理函数
	http.HandleFunc("/", httpHandleFunc)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	http.HandleFunc("/faros", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	addr := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func httpHandleFunc(w http.ResponseWriter, r *http.Request) {
	// 获取请求的 URL
	url := r.URL
	urlPath := url.String()

	// 读取 YAML 文件
	yamlFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// 解析 YAML 数据
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	// 获取请求的 Referer 头
	referer := r.Header.Get("Referer")
	// 判断是否是直接访问，如果直接访问的话，那么就判断是否登录
	if referer == "" {
		if hasCookie(r, "pubinternalsso", "qunheinternalsso") {
			pushImage(w, r, config, urlPath)
			return
		} else {
			imageRender(w, r, "image/noLogin.png")
			return
		}
	}
	// 判断是否防盗链
	for _, value := range config.SecurityChain {
		if strings.Contains(referer, value) {
			pushImage(w, r, config, urlPath)
			return
		} else {
			imageRender(w, r, "image/security.png")
			return
		}
	}
}
func hasCookie(r *http.Request, cookieNames ...string) bool {
	for _, name := range cookieNames {
		_, err := r.Cookie(name)
		if err == nil {
			return true
		}
	}
	return false
}

func pushImage(w http.ResponseWriter, r *http.Request, config Config, urlPath string) {
	for _, value := range config.URLPatterns {
		// 编译正则表达式
		re := regexp.MustCompile(value.Info.Pattern)
		// 匹配正则表达式
		matches := re.FindStringSubmatch(urlPath)
		var confluenceContent ConfluenceContent
		if len(matches) > value.Info.Fields.ContentId && len(matches) > value.Info.Fields.PageId {
			confluenceContent = query(value.Info.Fields.Type, matches[value.Info.Fields.ContentId], matches[value.Info.Fields.PageId])
			confluenceAttachmentPath := confluencePath(confluenceContent.PAGEID, confluenceContent.SPACEID, confluenceContent.CONTENTID)
			if len(confluenceAttachmentPath) < 20 {
				http.Error(w, "地址不正确", http.StatusBadRequest)
				return
			} else {
				imageRender(w, r, config.SourceBegin+confluenceAttachmentPath)
				return
			}
		}
	}
}

func imageRender(w http.ResponseWriter, r *http.Request, imagePath string) {
	image, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Error loading image:"+imagePath, http.StatusInternalServerError)
		return
	}
	defer image.Close()

	// 设置图像的 Content-Type
	w.Header().Set("Content-Type", "image/jpeg")

	// 将图像写入 ResponseWriter
	_, err = io.Copy(w, image)
	if err != nil {
		http.Error(w, "Error writing image", http.StatusInternalServerError)
		return
	}
}
