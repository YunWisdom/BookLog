

package console

import (
	"math"
	"net/http"

	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// GetTagsAction gets tags.
func GetTagsAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	var tags []*ConsoleTag
	tagModels := service.Tag.GetTags(math.MaxInt64, session.BID)
	for _, tagModel := range tagModels {
		tags = append(tags, &ConsoleTag{Title: tagModel.Title})
	}

	result.Data = tags
}
