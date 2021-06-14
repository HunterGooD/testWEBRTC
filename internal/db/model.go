package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"size:30"`
	Surname  string `gorm:"size:40"`
	Lastname string `gorm:"size:40"`
	Login    string `gorm:"size:20"`
	Password string `gorm:"size:70"`
	Avatar   string `gorm:"type:text"`
	Rooms    []Room `gorm:"many2many:user_rooms;"`
}

type Room struct {
	gorm.Model
	Name     string `gorm:"size:50"`
	Password string `gorm:"size:70"`
}
