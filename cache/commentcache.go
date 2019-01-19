

package cache

import (
	"github.com/YunWisdom/BookLog/model"
	"github.com/bluele/gcache"
)

// Comment service.
var Comment = &commentCache{
	idHolder: gcache.New(1024 * 10 * 10).LRU().Build(),
}

type commentCache struct {
	idHolder gcache.Cache
}

func (cache *commentCache) Put(comment *model.Comment) {
	if err := cache.idHolder.Set(comment.ID, comment); nil != err {
		logger.Errorf("put comment [id=%d] into cache failed: %s", comment.ID, err)
	}
}

func (cache *commentCache) Get(id uint) *model.Comment {
	ret, err := cache.idHolder.Get(id)
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get comment [id=%d] from cache failed: %s", id, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.Comment)
}
