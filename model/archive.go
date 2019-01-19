

// Package model is the "model" layer which defines entity structures with ORM and controller.
package model

// Archive model.
type Archive struct {
	Model

	Year         string `gorm:"size:4" json:"year"`
	Month        string `gorm:"size:2" json:"month"`
	ArticleCount int    `json:"articleCount"`

	BlogID uint64 `sql:"index" json:"blogID"`
}
