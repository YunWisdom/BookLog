

package service

import (
	"fmt"
	"sync"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/util"
)

// Navigation service.
var Navigation = &navigationService{
	mutex: &sync.Mutex{},
}

type navigationService struct {
	mutex *sync.Mutex
}

// Navigation pagination arguments of admin console.
const (
	adminConsoleNavigationListPageSize   = 15
	adminConsoleNavigationListWindowSize = 20
)

func (srv *navigationService) AddNavigation(navigation *model.Navigation) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	tx := db.Begin()
	if err := tx.Create(navigation).Error; nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *navigationService) RemoveNavigation(id, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	navigation := &model.Navigation{}

	tx := db.Begin()
	if err := tx.Where("`id` = ? AND `blog_id` = ?", id, blogID).Find(navigation).Error; nil != err {
		tx.Rollback()

		return err
	}
	if err := db.Delete(navigation).Error; nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *navigationService) UpdateNavigation(navigation *model.Navigation) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	count := 0
	if db.Model(&model.Navigation{}).Where("`id` = ? AND `blog_id` = ?", navigation.ID, navigation.BlogID).
		Count(&count); 1 > count {
		return fmt.Errorf("not found navigation [id=%d] to update", navigation.ID)
	}

	tx := db.Begin()
	if err := tx.Model(navigation).Updates(navigation).Error; nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *navigationService) ConsoleGetNavigations(page int, blogID uint64) (ret []*model.Navigation, pagination *util.Pagination) {
	offset := (page - 1) * adminConsoleNavigationListPageSize
	count := 0
	if err := db.Model(&model.Navigation{}).Order("`number` ASC, `id` DESC").
		Where("`blog_id` = ?", blogID).
		Count(&count).Offset(offset).Limit(adminConsoleNavigationListPageSize).Find(&ret).Error; nil != err {
		logger.Errorf("get navigations failed: " + err.Error())
	}

	pagination = util.NewPagination(page, adminConsoleNavigationListPageSize, adminConsoleNavigationListWindowSize, count)

	return
}

func (srv *navigationService) GetNavigations(blogID uint64) (ret []*model.Navigation) {
	if err := db.Model(&model.Navigation{}).Order("`number` ASC, `id` DESC").
		Where("`blog_id` = ?", blogID).Find(&ret).Error; nil != err {
		logger.Errorf("get navigations failed: " + err.Error())
	}

	return
}

func (srv *navigationService) ConsoleGetNavigation(id uint64) *model.Navigation {
	ret := &model.Navigation{}
	if err := db.First(ret, id).Error; nil != err {
		return nil
	}

	return ret
}
