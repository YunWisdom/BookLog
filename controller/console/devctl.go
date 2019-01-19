

package console

import (
	"net/http"
	"strconv"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
	"math/rand"
)

// GenArticlesAction generates articles for testing.
func GenArticlesAction(c *gin.Context) {
	session := util.GetSession(c)

	for i := 0; i < 100; i++ {
		article := &model.Article{
			AuthorID: session.UID,
			Title:    "title " + strconv.Itoa(i) + "_" + strconv.Itoa(rand.Int()),
			Tags:     "开发生成",
			Content:  "开发生成",
			BlogID:   session.BID,
		}
		if err := service.Article.AddArticle(article); nil != err {
			logger.Errorf("generate article failed: " + err.Error())
		}
	}

	c.Redirect(http.StatusTemporaryRedirect, model.Conf.Server)
}
