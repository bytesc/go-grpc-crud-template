package mysql_db

import (
	"gorm.io/gorm"
	"time"
)

// CrudList gorm结构体
type CrudList struct {
	gorm.Model // 引入模板结构体
	//ID        uint `gorm:"primarykey"`
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
	Name     string `gorm:"type: varchar(50); not null" binding:"required" json:"name" `
	Level    string `gorm:"type: varchar(20);" json:"level"`
	Email    string `gorm:"type: varchar(50);" json:"email"`
	Phone    string `gorm:"type: varchar(20);" json:"phone"`
	Birthday string `gorm:"type: varchar(20);" json:"birthday"`
	Address  string `gorm:"type: varchar(200);" json:"address"`
	// 变量名要大写，才public，可以被gorm获取解析
}

type UserList struct {
	Name        string    `gorm:"type: varchar(50); not null; primarykey" binding:"required" json:"name" `
	AuthGroup   string    `gorm:"type: varchar(20);" json:"auth_level"`
	Email       string    `gorm:"type: varchar(50); not null;" json:"email" binding:"required"`
	Phone       string    `gorm:"type: varchar(20);" json:"phone"`
	Password    string    `gorm:"type: varchar(100); not null;" json:"password" binding:"required"`
	PasswordTry uint      `json:"password_try"`
	LockedUntil time.Time `json:"locked_until"`
	Status      string    `json:"status"`
}
