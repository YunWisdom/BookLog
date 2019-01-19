

package model

import (
	"github.com/YunWisdom/BookLog/util"
)

// User model.
type User struct {
	Model

	Name              string `gorm:"size:32" json:"name"`
	Password          string `gorm:"size:255" json:"password"` // https://github.com/YunWisdom/BookLog/issues/130
	Nickname          string `gorm:"size:32" json:"nickname"`
	AvatarURL         string `gorm:"size:255" json:"avatarURL"`
	B3Key             string `gorm:"size:32" json:"b3Key"`
	Locale            string `gorm:"size:32" json:"locale"`
	TotalArticleCount int    `json:"totalArticleCount"`
	GithubId          string `gorm:"255" json:"githubId"` // 支持 GitHub 登录 https://github.com/YunWisdom/BookLog/issues/150
}

// User roles.
const (
	UserRoleNoLogin = iota
	UserRolePlatformAdmin
	UserRoleBlogAdmin
	UserRoleBlogUser
)

// AvatarURLWithSize returns avatar URL with the specified size.
func (u *User) AvatarURLWithSize(size int) string {
	return util.ImageSize(u.AvatarURL, size, size)
}
