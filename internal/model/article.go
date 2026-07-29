package model

import "time"

type Article struct {
	BaseModel
	//业务字段
	Title       string     `gorm:"type:varchar(200);not null;index" json:"title"`
	Summary     string     `gorm:"type:varchar(500)" json:"summary"`
	Content     string     `gorm:"type:longtext" json:"content"`
	Cover       string     `gorm:"type:varchar(255)" json:"cover"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	Status      int        `gorm:"type:tinyint;default:1;not nul" json:"status"`
	PublishedAt *time.Time `gorm:"type:datetime" json:"published_at"`

	//外键
	UserID     uint `json:"user_id"`
	CategoryID uint `json:"category_id"`

	//多对多
	Tags []Tag `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}

func (Article) TableName() string {
	return "articles"
}
