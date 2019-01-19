

package service

import (
	"sync"

	"github.com/YunWisdom/BookLog/model"
)

// Tag service.
var Tag = &tagService{
	mutex: &sync.Mutex{},
}

type tagService struct {
	mutex *sync.Mutex
}

func (srv *tagService) GetTags(size int, blogID uint64) (ret []*model.Tag) {
	if err := db.Where("`blog_id` = ?", blogID).Order("`article_count` DESC, `id` DESC").Limit(size).Find(&ret).Error; nil != err {
		logger.Errorf("get tags failed: " + err.Error())
	}

	return
}

func (srv *tagService) GetTagByTitle(title string, blogID uint64) *model.Tag {
	ret := &model.Tag{}
	if err := db.Where("`title` = ? AND `blog_id` = ?", title, blogID).First(ret).Error; nil != err {
		return nil
	}

	return ret
}
