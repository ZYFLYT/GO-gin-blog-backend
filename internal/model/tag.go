package model

type Tag struct {
	BaseModel
	TagName string `gorm:"type:varchar(50);uniqueIndex;not null" json:"tag_name"`
	Status  int    `gorm:"default:1" json:"status"`
}

func (Tag) TableName() string {
	return "tags"
}
