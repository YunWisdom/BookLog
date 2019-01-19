

package console

import (
	"net/http"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// AddUserAction adds a user.
func AddUserAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = "parses add user request failed"

		return
	}

	name := arg["name"].(string)
	user := service.User.GetUserByName(name)
	if nil == user {
		result.Code = -1
		result.Msg = "the user should login first"

		return
	}

	session := util.GetSession(c)
	if err := service.User.AddUserToBlog(user.ID, session.BID); nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}
}

// GetUsersAction gets users.
func GetUsersAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var users []*ConsoleUser
	userModels, pagination := service.User.GetBlogUsers(util.GetPage(c), session.BID)
	for _, userModel := range userModels {
		userBlog := service.User.GetUserBlog(userModel.ID, session.BID)
		users = append(users, &ConsoleUser{
			ID:           userModel.ID,
			Name:         userModel.Name,
			Nickname:     userModel.Nickname,
			Role:         userBlog.UserRole,
			URL:          blogURLSetting.Value + util.PathAuthors + "/" + userModel.Name,
			AvatarURL:    userModel.AvatarURL,
			ArticleCount: userBlog.UserArticleCount,
		})
	}

	result.Data = map[string]interface{}{
		"users":      users,
		"pagination": pagination,
	}
}
