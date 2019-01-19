

package controller

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
	model "github.com/YunWisdom/BookLog/model"
)

func showIndexAction(c *gin.Context) {
	t, err := template.ParseFiles(filepath.ToSlash(filepath.Join(model.Conf.StaticRoot, "console/dist/index.html")))
	if nil != err {
		logger.Errorf("load index page failed: " + err.Error())
		c.String(http.StatusNotFound, "load index page failed")

		return
	}

	t.Execute(c.Writer, nil)
}

func showPlatInfoAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	data := map[string]interface{}{}
	data["version"] = model.Version
	data["database"] = service.Database()
	data["mode"] = model.Conf.RuntimeMode
	data["server"] = model.Conf.Server
	data["staticServer"] = model.Conf.StaticServer
	data["staticResourceVer"] = model.Conf.StaticResourceVersion

	result.Data = data
}

func showTopBlogsAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	blogs := service.User.GetTopBlogs(10)
	for _, blog := range blogs {
		blog.ID = 0
		blog.UserID = 0
		blog.UserRole = 0
	}

	result.Data = blogs
}
