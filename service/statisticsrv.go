

package service

import (
	"strconv"
	"sync"

	"github.com/YunWisdom/BookLog/model"
	"github.com/jinzhu/gorm"
)

// Statistic service.
var Statistic = &statisticService{
	mutex: &sync.Mutex{},
}

type statisticService struct {
	mutex *sync.Mutex
}

func (srv *statisticService) GetAllStatistics(blogID uint64) []*model.Setting {
	var ret []*model.Setting

	if err := db.Where("`category` = ? AND `blog_id` = ?", model.SettingCategoryStatistic, blogID).Find(&ret).Error; nil != err {
		logger.Errorf("get all statistics failed: " + err.Error())

		return nil
	}

	return ret
}

func (srv *statisticService) GetStatistic(statisticName string, blogID uint) *model.Setting {
	ret := &model.Setting{}
	if err := db.Where("`name` = ? AND `category` = ? AND `blog_id` = ?", statisticName, model.SettingCategoryStatistic, blogID).Find(ret).Error; nil != err {
		logger.Errorf("get statistic failed: " + err.Error())

		return nil
	}

	return ret
}

func (srv *statisticService) GetStatistics(blogID uint, statisticNames ...string) map[string]*model.Setting {
	ret := map[string]*model.Setting{}
	var settings []*model.Setting
	if err := db.Where("`name` IN (?) AND `category` = ? AND `blog_id` = ?", statisticNames, model.SettingCategoryStatistic, blogID).Find(&settings).Error; nil != err {
		logger.Errorf("get statistics failed: " + err.Error())

		return nil
	}

	for _, setting := range settings {
		ret[setting.Name] = setting
	}

	return ret
}

func (srv *statisticService) IncViewCount(blogID uint64) error {
	tx := db.Begin()
	if err := srv.IncViewCountWithoutTx(tx, blogID); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *statisticService) IncViewCountWithoutTx(tx *gorm.DB, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	setting := &model.Setting{}
	if err := tx.Where("`name` = ? AND `category` = ? AND `blog_id` = ?", model.SettingNameStatisticViewCount, model.SettingCategoryStatistic, blogID).Find(setting).Error; nil != err {
		return err
	}

	count, err := strconv.Atoi(setting.Value)
	if nil != err {
		return err
	}

	setting.Value = strconv.Itoa(count + 1)
	if err := tx.Model(setting).Updates(setting).Error; nil != err {
		return err
	}

	return nil
}

func (srv *statisticService) IncArticleCount(blogID uint64) error {
	tx := db.Begin()
	if err := srv.IncArticleCountWithoutTx(tx, blogID); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *statisticService) IncArticleCountWithoutTx(tx *gorm.DB, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	setting := &model.Setting{}
	if err := tx.Where("`name` = ? AND `category` = ? AND `blog_id` = ?", model.SettingNameStatisticArticleCount, model.SettingCategoryStatistic, blogID).Find(setting).Error; nil != err {
		return err
	}

	count, err := strconv.Atoi(setting.Value)
	if nil != err {
		return err
	}

	setting.Value = strconv.Itoa(count + 1)
	if err := tx.Model(setting).Updates(setting).Error; nil != err {
		return err
	}

	return nil
}

func (srv *statisticService) DecArticleCount(blogID uint64) error {
	tx := db.Begin()
	if err := srv.DecArticleCountWithoutTx(tx, blogID); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *statisticService) DecArticleCountWithoutTx(tx *gorm.DB, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	setting := &model.Setting{}
	if err := tx.Where("`name` = ? AND `category` = ? AND `blog_id` = ?", model.SettingNameStatisticArticleCount, model.SettingCategoryStatistic, blogID).Find(setting).Error; nil != err {
		return err
	}

	count, err := strconv.Atoi(setting.Value)
	if nil != err {
		return err
	}

	setting.Value = strconv.Itoa(count - 1)
	if err := tx.Model(setting).Updates(setting).Error; nil != err {
		return err
	}

	return nil
}

func (srv *statisticService) IncCommentCount(blogID uint64) error {
	tx := db.Begin()
	if err := srv.IncCommentCountWithoutTx(tx, blogID); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *statisticService) IncCommentCountWithoutTx(tx *gorm.DB, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	setting := &model.Setting{}
	if err := tx.Where("`name` = ? AND `category` = ? AND `blog_id` = ?", model.SettingNameStatisticCommentCount, model.SettingCategoryStatistic, blogID).Find(setting).Error; nil != err {
		return err
	}

	count, err := strconv.Atoi(setting.Value)
	if nil != err {
		return err
	}

	setting.Value = strconv.Itoa(count + 1)
	if err := tx.Model(setting).Updates(setting).Error; nil != err {
		return err
	}

	return nil
}

func (srv *statisticService) DecCommentCount(blogID uint64) error {
	tx := db.Begin()
	if err := srv.DecCommentCountWithoutTx(tx, blogID); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func (srv *statisticService) DecCommentCountWithoutTx(tx *gorm.DB, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	setting := &model.Setting{}
	if err := tx.Where("`name` = ? AND `category` = ? AND `blog_id` = ?", model.SettingNameStatisticCommentCount, model.SettingCategoryStatistic, blogID).Find(setting).Error; nil != err {
		return err
	}

	count, err := strconv.Atoi(setting.Value)
	if nil != err {
		return err
	}
	setting.Value = strconv.Itoa(count - 1)
	if err := tx.Model(setting).Updates(setting).Error; nil != err {
		return err
	}

	return nil
}
