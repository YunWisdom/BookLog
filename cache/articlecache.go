

// Package cache includes caches.
package cache

import (
	"os"

	"github.com/YunWisdom/BookLog/log"
	"github.com/YunWisdom/BookLog/model"
	"github.com/bluele/gcache"
)

// Logger
var logger = log.NewLogger(os.Stdout)

// Article cache.
var Article = &articleCache{
	idHolder: gcache.New(1024 * 10).LRU().Build(),
}

type articleCache struct {
	idHolder gcache.Cache
}

func (cache *articleCache) Put(article *model.Article) {
	if err := cache.idHolder.Set(article.ID, article); nil != err {
		logger.Errorf("put article [id=%d] into cache failed: %s", article.ID, err)
	}
}

func (cache *articleCache) Get(id uint) *model.Article {
	ret, err := cache.idHolder.Get(id)
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get article [id=%d] from cache failed: %s", id, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.Article)
}
