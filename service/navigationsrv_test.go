

package service

import (
	"testing"

	"github.com/YunWisdom/BookLog/model"
)

func TestConsoleGetNavigations(t *testing.T) {
	navigations, pagination := Navigation.ConsoleGetNavigations(1, 1)

	if 1 != len(navigations) {
		t.Errorf("expected is [%d], actual is [%d]", 1, len(navigations))
	}
	if 1 != pagination.RecordCount {
		t.Errorf("expected is [%d], actual is [%d]", 1, pagination.RecordCount)
	}
}

func TestGetNavigations(t *testing.T) {
	navigations := Navigation.GetNavigations(1)

	if 1 != len(navigations) {
		t.Errorf("expected is [%d], actual is [%d]", 1, len(navigations))
	}
}

func TestConsoleGetNavigation(t *testing.T) {
	navigation := Navigation.ConsoleGetNavigation(1)
	if nil == navigation {
		t.Errorf("navigation is nil")

		return
	}

	if 1 != navigation.ID {
		t.Errorf("id is not [1]")
	}
}

func TestConsoleAddNavigation(t *testing.T) {
	navigation := &model.Navigation{
		Title:      "测试添加的导航",
		URL:        "https://b3log.org",
		IconURL:    "图标 URL",
		OpenMethod: model.NavigationOpenMethodBlank,
		Number:     3,
		BlogID:     1,
	}

	if err := Navigation.AddNavigation(navigation); nil != err {
		t.Errorf("add navigation failed: " + err.Error())

		return
	}

	navigation = Navigation.ConsoleGetNavigation(navigation.ID)
	if nil == navigation {
		t.Errorf("navigation is nil")

		return
	}

	if 2 != navigation.ID {
		t.Errorf("id is not [2]")
	}
}

func TestConsoleUpdateNavigation(t *testing.T) {
	navigation := Navigation.ConsoleGetNavigation(2)
	if nil == navigation {
		t.Errorf("navigation is nil")

		return
	}

	navigation.Title = "更新后的导航标题"
	if err := Navigation.UpdateNavigation(navigation); nil != err {
		t.Errorf("update navigation failed: " + err.Error())

		return
	}

	navigation = Navigation.ConsoleGetNavigation(2)
	if nil == navigation {
		t.Errorf("navigation is nil")

		return
	}

	if "更新后的导航标题" != navigation.Title {
		t.Errorf("expected is [%s], actual is [%s]", "更新后的导航标题", navigation.Title)
	}
}
