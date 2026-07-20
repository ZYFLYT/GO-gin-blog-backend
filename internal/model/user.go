package model

type User struct {
	BaseModel        //包含了ID,CreatedAt,UpdatedAt,DeletedAt
	Username  string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string `gorm:"type:varchar(255);not null" json:"password"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Nickname  string `gorm:"type:varchar(50);" json:"nickname"`
	Status    int    `gorm:"default:1" json:"status"`
}
