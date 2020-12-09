package golang

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"tbot/config"
	"tbot/internal/errors"
)

type Golang struct {
	golangStg *config.GolangSettings
	db        *gorm.DB
	ready     bool
}

func NewGolang(golangStg *config.GolangSettings, db *gorm.DB) *Golang {
	err := db.AutoMigrate(&Topic{}, &Tag{})
	errors.PanicIfErr(err)
	var topics []Topic
	topics = append(topics, ConvertQueryToTopic("Get executable dir\n(tags: get executable dir)\n---\n\nex, err := os.Executable()\ndir := filepath.Dir(ex)\nfmt.Println(\"dir:\", dir)\n"))
	topics = append(topics, ConvertQueryToTopic("Extract beginning of string (prefix)\n(tags: extract beginning string prefix)\n---\n\nt := string([]rune(s)[:5])"))
	topics = append(topics, ConvertQueryToTopic("Extract string suffix\n(tags: extract string suffix)\n---\n\nt := string([]rune(s)[len([]rune(s))-5:])"))
	topics = append(topics, ConvertQueryToTopic("Exec other program\n(tags: exec program)\n---\n\nerr := exec.Command(\"program\", \"arg1\", \"arg2\").Run()"))
	db.Create(&topics)

	//	topic := Topic{
	//		Title: "Get executable dir path",
	//		Code: `ex, err := os.Executable()
	//dir := filepath.Dir(ex)
	//fmt.Println("dir:", dir)`,
	//		Checked: false,
	//		Tags: []Tag{
	//			{Name: "get"},
	//			{Name: "executable"},
	//			{Name: "dir"},
	//			{Name: "path"},
	//		},
	//	}
	//	db.Create(&topic)

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
	if strings.HasPrefix(query, "+") {
		query = query[1:]
		newTopic := ConvertQueryToTopic(query)
		n.db.Create(&newTopic)
		tagsStr := ""
		for _, tag := range newTopic.Tags {
			tagsStr += " " + tag.Name
		}
		return fmt.Sprintf("%s\n(tags:%v)\n---\n%s", newTopic.Title, tagsStr, newTopic.Code), nil
	}

	queryTags := strings.Split(query, " ")
	var tags []Tag
	for _, tag := range queryTags {
		tags = append(tags, Tag{Name: tag})
	}
	n.db.Find(&tags)

	var topics []Topic
	n.db.Model(&tags).Association("Topics").Find(&topics)
	res := ""
	set := make(map[string]Topic)
	for _, topic := range topics {
		set[topic.Title] = topic
	}
	for _, topic := range set {
		res += "\n\n***\n" + fmt.Sprintf("%s\n(tags:%v)\n---\n%s", topic.Title, "tagsStr", topic.Code)
	}
	return res, nil
}
