package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"github.com/YunWisdom/BookLog/controller"
	"github.com/YunWisdom/BookLog/cron"
	"github.com/YunWisdom/BookLog/i18n"
	"github.com/YunWisdom/BookLog/log"
	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/theme"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// Logger
var logger *log.Logger

// The only one init function in pipe.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	log.SetLevel("warn")
	logger = log.NewLogger(os.Stdout)

	model.LoadConf()
	util.LoadMarkdown()
	i18n.Load()
	theme.Load()
	replaceServerConf()

	if "dev" == model.Conf.RuntimeMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
}

// Entry point.
func main() {
	service.ConnectDB()
	service.Upgrade.Perform()
	cron.Start()

	router := controller.MapRoutes()
	server := &http.Server{
		Addr:    "0.0.0.0:" + model.Conf.Port,
		Handler: router,
	}

	handleSignal(server)

	logger.Infof("Pipe (v%s) is running [%s]", model.Version, model.Conf.Server)
	server.ListenAndServe()
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Infof("got signal [%s], exiting pipe now", s)
		if err := server.Close(); nil != err {
			logger.Errorf("server close failed: " + err.Error())
		}

		service.DisconnectDB()

		logger.Infof("Pipe exited")
		os.Exit(0)
	}()
}

func replaceServerConf() {
	err := filepath.Walk(filepath.ToSlash(filepath.Join(model.Conf.StaticRoot, "theme")), func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".min.js") {
			data, e := ioutil.ReadFile(path)
			if nil != e {
				logger.Fatal("read file [" + path + "] failed: " + err.Error())
			}
			content := string(data)
			if !strings.Contains(content, "exports={Server:") {
				return err
			}

			json := "{Server:" + strings.Split(content, "{Server:")[1]
			json = strings.Split(json, "}}")[0] + "}"
			newJSON := "{Server:\"" + model.Conf.Server + "\",StaticServer:\"" + model.Conf.StaticServer + "\",StaticResourceVersion:\"" +
				model.Conf.StaticResourceVersion + "\",RuntimeMode:\"" + model.Conf.RuntimeMode + "\",AxiosBaseURL:\"" + model.Conf.AxiosBaseURL +
				"\",MockServer:\"" + model.Conf.MockServer + "\"}"
			content = strings.Replace(content, json, newJSON, -1)
			if e = ioutil.WriteFile(path, []byte(content), 0644); nil != e {
				logger.Fatal("replace server conf in [" + path + "] failed: " + err.Error())
			}
		}

		return err
	})
	if nil != err {
		logger.Fatal("replace server conf in [theme] failed: " + err.Error())
	}

	paths, err := filepath.Glob(filepath.ToSlash(filepath.Join(model.Conf.StaticRoot, "console/dist/*.js")))
	if 0 < len(paths) {
		for _, path := range paths {
			data, e := ioutil.ReadFile(path)
			if nil != e {
				logger.Fatal("read file [" + path + "] failed: " + err.Error())
			}
			content := string(data)
			if strings.Contains(content, "{rel:\"manifest") {
				json := "{rel:\"manifest\",href:\"" + strings.Split(content, "{rel:\"manifest\",href:\"")[1]
				json = strings.Split(json, "}]")[0] + "}"
				newJSON := "{rel:\"manifest\",href:\"" + model.Conf.StaticServer + "/theme/js/manifest.json\"}"
				content = strings.Replace(content, json, newJSON, -1)
			}
			if strings.Contains(content, "env:{Server:") {
				json := "env:{Server:" + strings.Split(content, "env:{Server:")[1]
				json = strings.Split(json, "}}")[0] + "}"
				newJSON := "env:{Server:\"" + model.Conf.Server + "\",StaticServer:\"" + model.Conf.StaticServer + "\",StaticResourceVersion:\"" +
					model.Conf.StaticResourceVersion + "\",RuntimeMode:\"" + model.Conf.RuntimeMode + "\",AxiosBaseURL:\"" + model.Conf.AxiosBaseURL +
					"\",MockServer:\"" + model.Conf.MockServer + "\"}"
				content = strings.Replace(content, json, newJSON, -1)
			}
			if strings.Contains(content, "/console/dist/") {
				part := strings.Split(content, "/console/dist/")[0]
				part = part[strings.LastIndex(part, "\"")+1:]
				content = strings.Replace(content, part, model.Conf.StaticServer, -1)
			}
			if e = ioutil.WriteFile(path, []byte(content), 0644); nil != e {
				logger.Fatal("replace server conf in [" + path + "] failed: " + err.Error())
			}
		}
	}

	if util.File.IsExist("console/dist/") { // dose not exist if npm run dev
		err = filepath.Walk(filepath.ToSlash(filepath.Join(model.Conf.StaticRoot, "console/dist/")), func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".html") {
				data, e := ioutil.ReadFile(path)
				if nil != e {
					logger.Fatal("read file [" + path + "] failed: " + err.Error())
				}
				content := string(data)
				if strings.Contains(content, "rel=\"manifest\" href=\"") {
					rel := "rel=\"manifest\" href=\"" + strings.Split(content, "rel=\"manifest\" href=\"")[1]
					rel = strings.Split(rel, "/>")[0] + "/>"
					newRel := "rel=\"manifest\" href=\"" + model.Conf.StaticServer + "/theme/js/manifest.json\"/>"
					content = strings.Replace(content, rel, newRel, -1)
				}
				if strings.Contains(content, "/console/dist/") {
					part := strings.Split(content, "/console/dist/")[0]
					part = part[strings.LastIndex(part, "\"")+1:]
					content = strings.Replace(content, part, model.Conf.StaticServer, -1)
				}
				v := fmt.Sprintf("%d", time.Now().Unix())
				content = strings.Replace(content, ".js\"", ".js?"+v+"\"", -1)
				content = strings.Replace(content, ".json\"", ".json?"+v+"\"", -1)
				if e = ioutil.WriteFile(path, []byte(content), 0644); nil != e {
					logger.Fatal("replace server conf in [" + path + "] failed: " + err.Error())
				}
			}

			return err
		})
		if nil != err {
			logger.Fatal("replace server conf in [theme] failed: " + err.Error())
		}
	}
}
