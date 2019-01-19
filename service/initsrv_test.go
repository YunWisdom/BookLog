

package service

import (
	"log"
	"os"
	"testing"

	"github.com/YunWisdom/BookLog/model"
	"github.com/YunWisdom/BookLog/util"
)

const (
	testPlatformAdminName = "pipe"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	home, err := util.UserHome()
	if nil != err {
		logger.Fatal(err)
	}

	model.Conf = &model.Configuration{}
	model.Conf.SQLite = home + "/pipe.test.db"

	if util.File.IsExist(model.Conf.SQLite) {
		os.Remove(model.Conf.SQLite)
	}

	ConnectDB()

	Init.InitPlatform(&model.User{
		Name:   testPlatformAdminName,
		B3Key:  "beyond",
		Locale: "zh_CN",
	})

	log.Println("setup tests")
}

func teardown() {
	DisconnectDB()

	log.Println("teardown tests")
}
