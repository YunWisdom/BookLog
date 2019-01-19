

package console

import (
	"net/http"
	"strconv"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// UpdateCategoryAction updates a category.
func UpdateCategoryAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	category := &model.Category{Model: model.Model{ID: uint64(id)}}
	if err := c.BindJSON(category); nil != err {
		result.Code = -1
		result.Msg = "parses update category request failed"

		return
	}

	session := util.GetSession(c)
	category.BlogID = session.BID

	if err := service.Category.UpdateCategory(category); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// GetCategoryAction gets a category.
func GetCategoryAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	data := service.Category.ConsoleGetCategory(id)
	if nil == data {
		result.Code = -1

		return
	}

	result.Data = data
}

// GetCategoriesAction gets categories.
func GetCategoriesAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	categoryModels, pagination := service.Category.ConsoleGetCategories(util.GetPage(c), session.BID)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var categories []*ConsoleCategory
	for _, categoryModel := range categoryModels {
		categories = append(categories, &ConsoleCategory{
			ID:          categoryModel.ID,
			Title:       categoryModel.Title,
			URL:         blogURLSetting.Value + util.PathCategories + categoryModel.Path,
			Description: categoryModel.Description,
			Number:      categoryModel.Number,
			Tags:        categoryModel.Tags,
		})
	}

	data := map[string]interface{}{}
	data["categories"] = categories
	data["pagination"] = pagination
	result.Data = data
}

// AddCategoryAction adds a category.
func AddCategoryAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	category := &model.Category{}
	if err := c.BindJSON(category); nil != err {
		result.Code = -1
		result.Msg = "parses add category request failed"

		return
	}

	category.BlogID = session.BID
	if err := service.Category.AddCategory(category); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// RemoveCategoryAction removes a category.
func RemoveCategoryAction(c *gin.Context) {
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
	if err := service.Category.RemoveCategory(id, blogID); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}
