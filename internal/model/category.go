package model

type Category struct {
	BaseModel
	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	//Slug        string `gorm:"type:varchar(50);uniqueIndex" json:"slug"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int    `gorm:"default:1" json:"status"`
}
