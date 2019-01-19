

package controller

import (
	"net/http"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Status represents platform status and blog info.
type Status struct {
	*service.PlatformStatus

	Name      string              `json:"name"`
	Nickname  string              `json:"nickname"`
	AvatarURL string              `json:"avatarURL"`
	BlogTitle string              `json:"blogTitle"`
	BlogURL   string              `json:"blogURL"`
	Role      int                 `json:"role"`
	Blogs     []*service.UserBlog `json:"blogs"`
}

func getStatusAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	platformStatus, err := service.Init.Status()
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	data := &Status{
		PlatformStatus: platformStatus,
	}

	session := util.GetSession(c)
	if 0 != session.UID {
		user := service.User.GetUser(session.UID)
		if nil == user {
			session := sessions.Default(c)
			session.Options(sessions.Options{
				Path:   "/",
				MaxAge: -1,
			})
			session.Clear()
			if err := session.Save(); nil != err {
				logger.Errorf("saves session failed: " + err.Error())
			}

			return
		}

		data.Name = user.Name
		data.Nickname = user.Nickname
		data.AvatarURL = user.AvatarURL
		data.Role = model.UserRoleBlogAdmin

		if model.UserRoleNoLogin != session.URole && platformStatus.Inited {
			ownBlog := service.User.GetOwnBlog(user.ID)
			if nil != ownBlog {
				data.BlogTitle = ownBlog.Title
				data.BlogURL = ownBlog.URL
			}
			data.Blogs = service.User.GetUserBlogs(user.ID)
		}
	}

	result.Data = data
}
