package golang

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"tbot/config"
)

type Golang struct {
	golangStg *config.GolangSettings
	db        *gorm.DB
	ready     bool
}

func NewGolang(golangStg *config.GolangSettings, db *gorm.DB) *Golang {

	initialize(db)

	return &Golang{
		golangStg: golangStg,
		db:        db,
		ready:     true,
	}
}

func initialize(db *gorm.DB) {
	//err := db.AutoMigrate(&Topic{}, &Tag{})
	//errors.PanicIfErr(err)
	//db.Create(ConvertQueryToTopic("Get executable dir\n(tags: get executable dir)\n---\n\nex, err := os.Executable()\ndir := filepath.Dir(ex)\nfmt.Println(\"dir:\", dir)\n"))
	//db.Create(ConvertQueryToTopic("Extract beginning of string (prefix)\n(tags: extract beginning string prefix)\n---\n\nt := string([]rune(s)[:5])"))
	//db.Create(ConvertQueryToTopic("Extract string suffix\n(tags: extract string suffix)\n---\n\nt := string([]rune(s)[len([]rune(s))-5:])"))
	//db.Create(ConvertQueryToTopic("Exec other program\n(tags: exec program)\n---\n\nerr := exec.Command(\"program\", \"arg1\", \"arg2\").Run()"))

	//db.Create(ConvertQueryToTopic("Exec other program\n(tags: exec program)\n---\n\nerr := exec.Command(\"program\", \"arg1\", \"arg2\").Run()"))
}

func (n *Golang) GetTag() string {
	return n.golangStg.Tag
}

func (n *Golang) IsReady() bool {
	return n.ready
}

func (n *Golang) Query(query string) (string, error) {
	formatStr := "*%s*\n_(tags:%v)_\n---\n`%s`"
	if strings.HasPrefix(query, "+") {
		query = query[1:]
		newTopic := ConvertQueryToTopic(query)
		n.db.Create(newTopic)
		tagsStr := ""
		for _, tag := range newTopic.Tags {
			tagsStr += " " + tag.Name
		}
		return fmt.Sprintf(formatStr, newTopic.Title, tagsStr, newTopic.Code), nil
	}

	if strings.HasPrefix(query, "*") {
		var topics []Topic
		n.db.Find(&topics)
		//var topic Topic
		n.db.Model(&topics).Association("Tags").Find(&topics)
		res := ""
		for _, topic := range topics {
			n.db.Model(&topic).Association("Tags").Find(&topic.Tags)
			topicTags := ""
			for _, tag := range topic.Tags {
				topicTags += " " + tag.Name
			}
			res += "\n\n===\n" + fmt.Sprintf(formatStr, topic.Title, topicTags, topic.Code)
		}
		return res, nil
	}

	queryTags := strings.Split(query, " ")
	var tags []Tag
	var tagsWhere []string
	for _, tag := range queryTags {
		//tags = append(tags, Tag{Name: tag})
		tagsWhere = append(tagsWhere, tag)
	}
	n.db.Where("name IN ?", tagsWhere).Find(&tags)

	var topics []Topic
	n.db.Model(&tags).Association("Topics").Find(&topics)
	res := ""
	set := make(map[string]Topic)
	for _, topic := range topics {
		set[topic.Title] = topic
	}
	for _, topic := range set {
		n.db.Model(&topic).Association("Tags").Find(&topic.Tags)
		topicTags := ""
		for _, tag := range topic.Tags {
			topicTags += " " + tag.Name
		}
		res += "\n\n===\n" + fmt.Sprintf(formatStr, topic.Title, topicTags, topic.Code)
	}
	return res, nil
}
