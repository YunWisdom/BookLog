

package controller

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/parnurzeal/gorequest"
)

// DataModel represents data model.
type DataModel map[string]interface{}

const nilB3id = "H9oxzSym"

func fillUser(c *gin.Context) {
	inited := service.Init.Inited()
	if !inited && util.PathInit != c.Request.URL.Path {
		c.Redirect(http.StatusSeeOther, model.Conf.Server+util.PathInit)
		c.Abort()

		return
	}

	dataModel := &DataModel{}
	c.Set("dataModel", dataModel)
	session := util.GetSession(c)
	(*dataModel)["User"] = session
	if 0 != session.UID {
		c.Next()

		return
	}

	uaStr := c.Request.UserAgent()
	if isBot(uaStr) {
		logger.Tracef("Bot User-Agent [%s]", uaStr)
		c.Next()

		return
	}

	if strings.Contains(c.Request.URL.Path, util.PathAPIsSymArticle) || strings.Contains(c.Request.URL.Path, util.PathAPIsSymComment) { // https://github.com/YunWisdom/BookLog/issues/138
		c.Next()

		return
	}

	b3id := c.Request.URL.Query().Get("b3id")
	switch b3id {
	case nilB3id:
		c.Next()

		return
	case "":
		redirectURL := model.Conf.Server + c.Request.URL.Path
		if strings.Contains(c.Request.URL.Path, util.PathAdmin) {
			redirectURL = model.Conf.Server + c.Request.URL.Path
			if "" != c.Request.URL.RawQuery {
				redirectURL += "?" + c.Request.URL.RawQuery
			}
		}
		if strings.HasPrefix(c.Request.URL.Path, util.PathBlogs) {
			name := c.Request.URL.Path[len(util.PathBlogs)+1:]
			name = strings.Split(name, "?")[0]
			name = strings.Split(name, "/")[0]
			if "" != name {
				user := service.User.GetUserByName(name)
				if nil != user {
					userBlog := service.User.GetOwnBlog(user.ID)
					blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, userBlog.ID)
					redirectURL = blogURLSetting.Value + strings.Split(c.Request.URL.Path, util.PathBlogs+"/"+name)[1]
					if "" != c.Request.URL.RawQuery {
						redirectURL += "?" + c.Request.URL.RawQuery
					}
				}
			}
		} else {
			if !strings.HasPrefix(redirectURL, model.Conf.Server) {
				redirectURL = model.Conf.Server + c.Request.URL.Path
			}
		}
		redirectURL = strings.TrimSpace(redirectURL)
		if "" == redirectURL {
			redirectURL = model.Conf.Server + c.Request.URL.Path
		}
		redirectURL = url.QueryEscape(redirectURL)
		c.Redirect(http.StatusSeeOther, util.HacPaiURL+"/apis/b3-identity?goto="+redirectURL)
		c.Abort()

		return
	default:
		result := util.NewResult()
		_, _, errs := gorequest.New().Get(util.HacPaiURL+"/apis/check-b3-identity?b3id="+b3id).
			Set("user-agent", model.UserAgent).Timeout(5*time.Second).
			Retry(3, 2*time.Second).EndStruct(result)
		if nil != errs {
			logger.Errorf("check b3 identity failed: %s", errs)
			c.Next()

			return
		}

		if 0 != result.Code {
			c.Next()

			return
		}

		data := result.Data.(map[string]interface{})
		username := data["userName"].(string)
		b3Key := data["userB3Key"].(string)
		userAvatar := data["userAvatarURL"].(string)

		session = &util.SessionData{
			UName:   username,
			UB3Key:  b3Key,
			UAvatar: userAvatar,
			URole:   model.UserRoleBlogAdmin,
		}

		user := &model.User{
			Name:      session.UName,
			B3Key:     b3Key,
			AvatarURL: session.UAvatar,
		}

		if service.Init.Inited() {
			if err := service.Init.InitBlog(user); nil != err {
				logger.Errorf("init user [name=%s] blog failed: %s", username, err.Error())
			}
		}

		if existUser := service.User.GetUserByName(username); nil != existUser {
			existUser.AvatarURL = session.UAvatar
			ownBlog := service.User.GetOwnBlog(existUser.ID)
			if nil != ownBlog {
				session.BID = ownBlog.ID
				session.BURL = ownBlog.URL
				session.URole = ownBlog.UserRole
			}
			session.UID = existUser.ID
			session.UB3Key = existUser.B3Key

			service.User.UpdateUser(existUser)
		} else {
			if err := service.User.AddUser(user); nil != err {
				logger.Errorf("add user [name=%s] failed: %s", username, err.Error())
			}

			session.UID = user.ID
		}

		if err := session.Save(c); nil != err {
			result.Code = -1
			result.Msg = "saves session failed: " + err.Error()
		}

		(*dataModel)["User"] = session

		if util.PathLogin == c.Request.URL.Path || util.PathRegister == c.Request.URL.Path {
			c.Redirect(http.StatusSeeOther, model.Conf.Server+util.PathAdmin)
			c.Abort()

			return
		}

		c.Next()
	}
}

func isBot(uaStr string) bool {
	var ua = user_agent.New(uaStr)

	return ua.Bot() || strings.HasPrefix(uaStr, "Sym")
}
