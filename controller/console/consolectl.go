

package console

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/YunWisdom/BookLog/model"
	"github.com/gin-gonic/gin"
)

// ShowAdminPagesAction shows admin pages.
func ShowAdminPagesAction(c *gin.Context) {
	t, err := template.ParseFiles(filepath.ToSlash(filepath.Join(model.Conf.StaticRoot, "console/dist/admin"+c.Param("path")+"/index.html")))
	if nil != err {
		logger.Errorf("load console page [" + c.Param("path") + "] failed: " + err.Error())
		c.String(http.StatusNotFound, "load console page failed")

		return
	}

	t.Execute(c.Writer, nil)
}
