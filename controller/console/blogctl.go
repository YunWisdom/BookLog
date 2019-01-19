

package console

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// BlogSwitchAction switches blog.
func BlogSwitchAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	blogID, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = -1

		return
	}

	session := util.GetSession(c)
	userID := session.UID

	userBlogs := service.User.GetUserBlogs(userID)
	if 1 > len(userBlogs) {
		result.Code = -1
		result.Msg = "switch blog failed"

		return
	}

	role := -1
	for _, userBlog := range userBlogs {
		if userBlog.ID == uint64(blogID) {
			role = userBlog.UserRole

			break
		}
	}

	if -1 == role {
		result.Code = -1
		result.Msg = "switch blog failed"

		return
	}

	result.Data = role

	session.URole = role
	session.BID = uint64(blogID)
	session.Save(c)
}

// CheckVersionAction checks version.
func CheckVersionAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	rhyResult := map[string]interface{}{}
	request := gorequest.New()
	_, _, errs := request.Get("https://rhythm.b3log.org/version/pipe/latest/"+model.Version).
		Set("User-Agent", model.UserAgent).Timeout(30*time.Second).
		Retry(3, 5*time.Second).EndStruct(&rhyResult)
	if nil != errs {
		result.Code = -1
		result.Msg = errs[0].Error()

		return
	}

	data := map[string]interface{}{}
	data["version"] = rhyResult["pipeVersion"]
	data["download"] = rhyResult["pipeDownload"]
	result.Data = data
}
