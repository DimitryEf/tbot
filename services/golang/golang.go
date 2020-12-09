package golang

import (
	"gorm.io/gorm"
	"tbot/config"
	"tbot/internal/errors"
)

type Golang struct {
	golangStg *config.GolangSettings
	db        *gorm.DB
	ready     bool
}

func NewGolang(golangStg *config.GolangSettings, db *gorm.DB) *Golang {
	err := db.AutoMigrate(&Topic{}, &Tag{}, &TopicTag{})
	errors.PanicIfErr(err)
	db.Create(&Topic{
		Code:    "test",
		Checked: false,
	})
	return &Golang{
		golangStg: golangStg,
		db:        db,
		ready:     true,
	}
}

func (n *Golang) GetTag() string {
	return n.golangStg.Tag
}

func (n *Golang) IsReady() bool {
	return n.ready
}

func (n *Golang) Query(query string) (string, error) {
	var topic Topic
	n.db.First(&topic)
	return topic.Code, nil
}
