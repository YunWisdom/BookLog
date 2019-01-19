

package controller

import (
	"github.com/YunWisdom/BookLog/service"
	"github.com/gin-gonic/gin"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

func outputSitemapAction(c *gin.Context) {
	sm := stm.NewSitemap()
	sm.Create()

	blogs := service.User.GetTopBlogs(10)
	for _, blog := range blogs {
		sm.Add(stm.URL{"loc": blog.URL})
	}

	c.Writer.Write(sm.XMLContent())
}
