package golang

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	Id      int `gorm:"primaryKey"`
	Code    string
	Checked bool
}

type Tag struct {
	gorm.Model
	Id   int `gorm:"primaryKey"`
	Name string
}

type TopicTag struct {
	gorm.Model
	Id      int `gorm:"primaryKey"`
	TopicId int
	TagId   int
}
