

package controller

import (
	"fmt"
	"net/http"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

var states = map[string]string{}

// redirectGitHubLoginAction redirects to GitHub auth page.
func redirectGitHubLoginAction(c *gin.Context) {
	state := util.RandString(16) + model.Conf.Server
	states[state] = state
	path := "https://github.com/login/oauth/authorize" + "?client_id=af7df3c80f26af88a8b3&state=" + state + "&scope=public_repo,user"

	c.Redirect(http.StatusSeeOther, path)
}

func githubCallbackAction(c *gin.Context) {
	state := c.Query("state")
	if _, exist := states[state]; !exist {
		c.Status(http.StatusForbidden)

		return
	}
	delete(states, state)

	accessToken := c.Query("ak")
	githubUser := util.GitHubUserInfo(accessToken)
	if nil == githubUser {
		c.Status(http.StatusForbidden)

		return
	}

	githubId := fmt.Sprintf("%v", githubUser["id"])
	userName := githubUser["login"].(string)
	user := service.User.GetUserByGitHubId(githubId)
	if nil == user {
		if !service.Init.Inited() {
			user = &model.User{
				Name:      userName,
				Password:  util.RandString(8),
				AvatarURL: githubUser["avatar_url"].(string),
				GithubId:  githubId,
			}

			if err := service.Init.InitPlatform(user); nil != err {
				logger.Errorf("init platform via github login failed: " + err.Error())
				c.Status(http.StatusInternalServerError)

				return
			}
		} else {
			if !model.Conf.OpenRegister {
				c.Status(http.StatusForbidden)

				return
			}

			user = service.User.GetUserByName(userName)
			if nil == user {
				user = &model.User{
					Name:      userName,
					Password:  util.RandString(8),
					AvatarURL: githubUser["avatar_url"].(string),
					GithubId:  githubId,
				}

				if err := service.Init.InitBlog(user); nil != err {
					logger.Errorf("init blog via github login failed: " + err.Error())
					c.Status(http.StatusInternalServerError)

					return
				}
			} else {
				user.GithubId = githubId
				service.User.UpdateUser(user)
			}
		}
	}

	ownBlog := service.User.GetOwnBlog(user.ID)
	session := &util.SessionData{
		UID:     user.ID,
		UName:   user.Name,
		UB3Key:  user.B3Key,
		UAvatar: user.AvatarURL,
		URole:   ownBlog.UserRole,
		BID:     ownBlog.ID,
		BURL:    ownBlog.URL,
	}
	if err := session.Save(c); nil != err {
		logger.Errorf("saves session failed: " + err.Error())
		c.Status(http.StatusInternalServerError)
	}

	c.Redirect(http.StatusSeeOther, model.Conf.Server+util.PathAdmin)
}
