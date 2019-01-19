

package console

import (
	"net/http"
	"strconv"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// GetNavigationsAction gets navigations.
func GetNavigationsAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	navigationModels, pagination := service.Navigation.ConsoleGetNavigations(util.GetPage(c), session.BID)

	var navigations []*ConsoleNavigation
	for _, navigationModel := range navigationModels {
		comment := &ConsoleNavigation{
			ID:         navigationModel.ID,
			Title:      navigationModel.Title,
			URL:        navigationModel.URL,
			IconURL:    navigationModel.IconURL,
			OpenMethod: navigationModel.OpenMethod,
			Number:     navigationModel.Number,
		}

		navigations = append(navigations, comment)
	}

	data := map[string]interface{}{}
	data["navigations"] = navigations
	data["pagination"] = pagination
	result.Data = data
}

// GetNavigationAction gets a navigation.
func GetNavigationAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = -1

		return
	}

	data := service.Navigation.ConsoleGetNavigation(uint64(id))
	if nil == data {
		result.Code = -1

		return
	}

	result.Data = data
}

// RemoveNavigationAction remove a navigation.
func RemoveNavigationAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	session := util.GetSession(c)
	blogID := session.BID

	if err := service.Navigation.RemoveNavigation(uint64(id), blogID); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// UpdateNavigationAction updates a navigation.
func UpdateNavigationAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	navigation := &model.Navigation{Model: model.Model{ID: uint64(id)}}
	if err := c.BindJSON(navigation); nil != err {
		result.Code = -1
		result.Msg = "parses update navigation request failed"

		return
	}

	session := util.GetSession(c)
	navigation.BlogID = session.BID

	if err := service.Navigation.UpdateNavigation(navigation); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// AddNavigationAction adds a navigation.
func AddNavigationAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	navigation := &model.Navigation{}
	if err := c.BindJSON(navigation); nil != err {
		result.Code = -1
		result.Msg = "parses add navigation request failed"

		return
	}

	navigation.BlogID = session.BID
	if err := service.Navigation.AddNavigation(navigation); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}
