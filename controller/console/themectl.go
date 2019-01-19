

package console

import (
	"net/http"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/theme"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// UpdateThemeAction updates theme.
func UpdateThemeAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	theme := c.Param("id")
	session := util.GetSession(c)

	settings := []*model.Setting{
		{
			Category: model.SettingCategoryTheme,
			Name:     model.SettingNameThemeName,
			Value:    theme,
			BlogID:   session.BID,
		},
	}
	if err := service.Setting.UpdateSettings(model.SettingCategoryTheme, settings, session.BID); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// GetThemesAction gets themes.
func GetThemesAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	currentID := theme.Themes[0]
	themeNameSetting := service.Setting.GetSetting(model.SettingCategoryTheme, model.SettingNameThemeName, session.BID)
	if nil == themeNameSetting {
		logger.Errorf("not found theme name setting")
	} else {
		currentID = themeNameSetting.Value
	}

	var themes []*ConsoleTheme
	for _, themeName := range theme.Themes {
		consoleTheme := &ConsoleTheme{
			Name:         themeName,
			ThumbnailURL: model.Conf.Server + "/theme/x/" + themeName + "/thumbnail.jpg",
		}

		themes = append(themes, consoleTheme)
	}

	result.Data = map[string]interface{}{
		"currentId": currentID,
		"themes":    themes,
	}
}
