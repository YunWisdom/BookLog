

package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// addSymCommentAction adds a comment come from Sym. Sees https://hacpai.com/article/1457158841475 for more details.
func addSymCommentAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = "parses add comment request failed"

		return
	}

	requestCmt := arg["comment"].(map[string]interface{})
	articleId, err := strconv.ParseUint(requestCmt["articleId"].(string), 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = "parses add comment request failed"

		return
	}

	blogID := getBlogID(c)
	comment := &model.Comment{
		BlogID:          blogID,
		ArticleID:       articleId,
		AuthorID:        model.SyncCommentAuthorID,
		Content:         requestCmt["content"].(string),
		IP:              requestCmt["ip"].(string),
		UserAgent:       requestCmt["ua"].(string),
		AuthorName:      requestCmt["authorName"].(string),
		AuthorURL:       requestCmt["authorURL"].(string),
		AuthorAvatarURL: requestCmt["authorAvatarURL"].(string),
	}

	if err := service.Comment.AddComment(comment); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}

	if err := service.Comment.UpdatePushedAt(comment); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// addSymArticleAction adds an article come from Sym. Sees https://hacpai.com/article/1457158841475 for more details.
func addSymArticleAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = "parses add article request failed"

		return
	}

	requestArticle := arg["article"].(map[string]interface{})
	articleId, err := strconv.ParseUint(requestArticle["oId"].(string), 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = "parses add article request failed"

		return
	}

	blogID := getBlogID(c)
	b3Key := requestArticle["userB3Key"].(string)
	blogAdmin := service.User.GetBlogAdmin(blogID)
	if b3Key != blogAdmin.B3Key {
		result.Code = -1
		result.Msg = "B3 key not match, ignored add article"

		return
	}

	article := &model.Article{
		BlogID:   blogID,
		AuthorID: blogAdmin.ID,
		Title:    requestArticle["articleTitle"].(string),
		Tags:     requestArticle["articleTags"].(string),
		Content:  requestArticle["articleContent"].(string),
		Commentable: true,
	}
	article.ID = articleId

	if err := service.Article.AddArticle(article); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}

	if err := service.Article.UpdatePushedAt(article); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

// updateSymArticleAction updates an article come from Sym. Sees https://hacpai.com/article/1457158841475 for more details.
func updateSymArticleAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = "parses update article request failed"

		return
	}

	requestArticle := arg["article"].(map[string]interface{})
	articleId, err := strconv.ParseUint(requestArticle["oId"].(string), 10, 64)
	if nil != err {
		result.Code = -1
		result.Msg = "parses update article request failed"

		return
	}

	blogID := getBlogID(c)
	b3Key := requestArticle["userB3Key"].(string)
	blogAdmin := service.User.GetBlogAdmin(blogID)
	if b3Key != blogAdmin.B3Key {
		result.Code = -1
		result.Msg = "B3 key not match, ignored add article"

		return
	}

	article := service.Article.ConsoleGetArticle(articleId)
	if nil == article {
		result.Code = -1
		result.Msg = "not found article [ID=" + fmt.Sprintf("%d", articleId) + "] to update"

		return
	}

	article.Title = requestArticle["articleTitle"].(string)
	article.Tags = requestArticle["articleTags"].(string)
	article.Content = requestArticle["articleContent"].(string)

	if err := service.Article.UpdateArticle(article); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}

	if err := service.Article.UpdatePushedAt(article); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}
