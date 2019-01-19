

package cron

import (
	"net/url"
	"strings"
	"time"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/parnurzeal/gorequest"
)

func pushArticlesPeriodically() {
	go pushArticles()

	go func() {
		for range time.Tick(time.Second * 30) {
			pushArticles()
		}
	}()
}

func pushArticles() {
	defer util.Recover()

	server, _ := url.Parse(model.Conf.Server)
	if !util.IsDomain(server.Hostname()) {
		return
	}

	articles := service.Article.GetUnpushedArticles()
	for _, article := range articles {
		author := service.User.GetUser(article.AuthorID)
		b3Key := author.B3Key
		b3Name := author.Name
		if "" == b3Key && !strings.Contains(model.Conf.Server, "pipe.b3log.org") {
			pa := service.User.GetPlatformAdmin()
			b3Key = pa.B3Key
			b3Name = pa.Name
		}
		if "" == b3Key {
			continue
		}

		blogTitleSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogTitle, article.BlogID)
		blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, article.BlogID)
		requestJSON := map[string]interface{}{
			"article": map[string]interface{}{
				"id":        article.ID,
				"title":     article.Title,
				"permalink": article.Path,
				"tags":      article.Tags,
				"content":   article.Content,
			},
			"client": map[string]interface{}{
				"name":  "Pipe",
				"ver":   model.Version,
				"title": blogTitleSetting.Value,
				"host":  blogURLSetting.Value,
				"email": b3Name,
				"key":   b3Key,
			},
		}
		result := &map[string]interface{}{}
		_, _, errs := gorequest.New().Post("https://rhythm.b3log.org/api/article").SendMap(requestJSON).
			Set("user-agent", model.UserAgent).Timeout(30*time.Second).
			Retry(3, 5*time.Second).EndStruct(result)
		if nil != errs {
			logger.Errorf("push article to Rhythm failed: " + errs[0].Error())
		}

		article.PushedAt = article.UpdatedAt
		service.Article.UpdatePushedAt(article)
	}
}

func pushCommentsPeriodically() {
	go pushComments()

	go func() {
		for range time.Tick(time.Second * 30) {
			pushComments()
		}
	}()
}

func pushComments() {
	defer util.Recover()

	server, _ := url.Parse(model.Conf.Server)
	if !util.IsDomain(server.Hostname()) {
		return
	}

	comments := service.Comment.GetUnpushedComments()
	for _, comment := range comments {
		author := service.User.GetUser(comment.AuthorID)
		article := service.Article.ConsoleGetArticle(comment.ArticleID)
		articleAuthor := service.User.GetUser(article.AuthorID)
		b3Key := articleAuthor.B3Key
		b3Name := articleAuthor.Name
		if "" == b3Key {
			pa := service.User.GetPlatformAdmin()
			b3Key = pa.B3Key
			b3Name = pa.Name
		}
		if "" == b3Key {
			continue
		}

		blogTitleSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogTitle, comment.BlogID)
		requestJSON := map[string]interface{}{
			"comment": map[string]interface{}{
				"id":          comment.ID,
				"articleId":   comment.ArticleID,
				"content":     comment.Content,
				"authorName":  author.Name,
				"authorEmail": "",
			},
			"client": map[string]interface{}{
				"name":  "Pipe",
				"ver":   model.Version,
				"title": blogTitleSetting.Value,
				"host":  model.Conf.Server,
				"email": b3Name,
				"key":   b3Key,
			},
		}
		result := &map[string]interface{}{}
		_, _, errs := gorequest.New().Post("https://rhythm.b3log.org/api/comment").SendMap(requestJSON).
			Set("user-agent", model.UserAgent).Timeout(30*time.Second).
			Retry(3, 5*time.Second).EndStruct(result)
		if nil != errs {
			logger.Errorf("push comment to Rhythm failed: " + errs[0].Error())
		}

		comment.PushedAt = comment.UpdatedAt
		service.Comment.UpdatePushedAt(comment)
	}
}
