

package model

import (
	"time"
)

// Article model.
type Article struct {
	Model

	AuthorID     uint64    `json:"authorID" structs:"authorID"`
	Title        string    `gorm:"size:128" json:"title" structs:"title"`
	Abstract     string    `gorm:"type:mediumtext" json:"abstract" structs:"abstract"`
	Tags         string    `gorm:"type:text" json:"tags" structs:"tags"`
	Content      string    `gorm:"type:mediumtext" json:"content" structs:"content"`
	Path         string    `sql:"index" gorm:"size:255" json:"path" structs:"path"`
	Status       int       `sql:"index" json:"status" structs:"status"`
	Topped       bool      `json:"topped" structs:"topped"`
	Commentable  bool      `json:"commentable" structs:"commentable"`
	ViewCount    int       `json:"viewCount" structs:"viewCount"`
	CommentCount int       `json:"commentCount" structs:"commentCount"`
	IP           string    `gorm:"size:128" json:"ip" structs:"ip"`
	UserAgent    string    `gorm:"size:255" json:"userAgent" structs:"userAgent"`
	PushedAt     time.Time `json:"pushedAt" structs:"pushedAt"`

	BlogID uint64 `sql:"index" json:"blogID" structs:"blogID"`
}

// Article statuses.
const (
	ArticleStatusOK = iota
)
