

package service

import (
	"testing"

	"github.com/YunWisdom/BookLog/model"
)

func TestGetSetting(t *testing.T) {
	setting := Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogTitle, 1)
	if nil == setting {
		t.Errorf("setting is nil")

		return
	}

	if "pipe 的博客" != setting.Value {
		t.Errorf("expected is [%s], actual is [%s]", "pipe 的博客", setting.Value)
	}
}

func TestGetAllSettings(t *testing.T) {
	settings := Setting.GetAllSettings(1)
	settingsCount := 26
	if settingsCount != len(settings) {
		t.Errorf("expected is [%d], actual is [%d]", settingsCount, len(settings))
	}
}

func TestGetCategorySettings(t *testing.T) {
	basicSettings := Setting.GetCategorySettings(model.SettingCategoryBasic, 1)
	if 11 != len(basicSettings) {
		t.Errorf("expected is [%d], actual is [%d]", 10, len(basicSettings))
	}
}

func TestGetSettings(t *testing.T) {
	settings := Setting.GetSettings(model.SettingCategoryBasic,
		[]string{model.SettingNameBasicBlogTitle, model.SettingNameBasicBlogSubtitle}, 1)
	if nil == settings {
		t.Errorf("settings is nil")

		return
	}
	if 1 > len(settings) {
		t.Errorf("settings is empty")

		return
	}

	if "pipe 的博客" != settings[model.SettingNameBasicBlogTitle].Value {
		t.Errorf("expected is [%s], actual is [%s]", "pipe 的博客", settings[model.SettingNameBasicBlogTitle].Value)
	}
}

func TestUpdateSettings(t *testing.T) {
	settings := Setting.GetSettings(model.SettingCategoryBasic,
		[]string{model.SettingNameBasicBlogTitle, model.SettingNameBasicBlogSubtitle}, 1)
	settings[model.SettingNameBasicBlogTitle].Value = "更新后的标题"
	var basics []*model.Setting
	for _, setting := range settings {
		basics = append(basics, setting)
	}
	if err := Setting.UpdateSettings(model.SettingCategoryBasic, basics, 1); nil != err {
		t.Errorf("updates settings failed: " + err.Error())

		return
	}

	settings = Setting.GetSettings(model.SettingCategoryBasic,
		[]string{model.SettingNameBasicBlogTitle, model.SettingNameBasicBlogSubtitle}, 1)
	if "更新后的标题" != settings[model.SettingNameBasicBlogTitle].Value {
		t.Errorf("expected is [%s], actual is [%s]", "更新后的标题", settings[model.SettingNameBasicBlogTitle].Value)
	}
}
