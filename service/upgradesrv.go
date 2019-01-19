

package service

import (
	"sync"

	"github.com/YunWisdom/BookLog/model"
)

// Upgrade service.
var Upgrade = &upgradeService{
	mutex: &sync.Mutex{},
}

type upgradeService struct {
	mutex *sync.Mutex
}

const (
	fromVer = "1.8.4"
	toVer   = model.Version
)

func (srv *upgradeService) Perform() {
	if !Init.Inited() {
		return
	}
	sysVerSetting := Setting.GetSetting(model.SettingCategorySystem, model.SettingNameSystemVer, 1)
	if nil == sysVerSetting {
		logger.Fatalf("system state is error, please contact developer: https://github.com/YunWisdom/BookLog/issues/new")
	}

	currentVer := sysVerSetting.Value
	if model.Version == currentVer {
		return
	}

	if fromVer == currentVer {
		perform()

		return
	}

	logger.Fatalf("attempt to skip more than one version to upgrade. Expected: %s, Actually: %s", fromVer, currentVer)
}

func perform() {
	logger.Infof("upgrading from version [%s] to version [%s]....", fromVer, toVer)

	var allSettings []model.Setting
	if err := db.Find(&allSettings).Error; nil != err {
		logger.Fatalf("load settings failed: %s", err)
	}

	var updateSettings []model.Setting
	for _, setting := range allSettings {
		if model.SettingNameSystemVer == setting.Name {
			setting.Value = model.Version
			updateSettings = append(updateSettings, setting)

			continue
		}
	}

	tx := db.Begin()
	for _, setting := range updateSettings {
		if err := tx.Save(setting).Error; nil != err {
			tx.Rollback()

			logger.Fatalf("update setting [%+v] failed: %s", setting, err.Error())
		}
	}

	rows, err := tx.Model(&model.Setting{}).Select("`blog_id`").Group("`blog_id`").Rows()
	defer rows.Close()
	if nil != err {
		tx.Rollback()

		logger.Fatalf("update settings failed: %s", err.Error())
	}
	for rows.Next() {
		var blogID uint64
		err := rows.Scan(&blogID)
		if nil != err {
			tx.Rollback()

			logger.Fatalf("update settings failed: %s", err.Error())
		}

		googleAdSenseArticleEmbedSetting := &model.Setting{
			Category: model.SettingCategoryAd,
			Name:     model.SettingNameAdGoogleAdSenseArticleEmbed,
			Value:    "",
			BlogID:   blogID}
		if err := Setting.AddSetting(googleAdSenseArticleEmbedSetting); nil != err {
			logger.Error("create Google AdSense setting failed: " + err.Error())
		}
	}

	tx.Commit()

	logger.Infof("upgraded from version [%s] to version [%s] successfully :-)", fromVer, toVer)
}
