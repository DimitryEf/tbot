package golang

import (
	"fmt"
	"gorm.io/gorm"
	"sort"
	"strings"
	"tbot/config"
	"tbot/tools"
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

func create(db *gorm.DB, query string) Topic {
	topic, err := ConvertQueryToTopic(query)
	if err != nil {
		return Topic{}
	}
	tags := topic.Tags
	for i, tag := range tags {
		db.Where("name = ?", tag.Name).Find(&tag)
		if tag.Id == 0 {
			db.Create(&tag)
		}
		tags[i] = tag
	}
	topic.Tags = tags
	db.Create(topic)
	return *topic
}

func (n *Golang) GetTag() string {
	return n.golangStg.Tag
}

func (n *Golang) IsReady() bool {
	return n.ready
}

func (n *Golang) Query(query string) (string, error) {
	formatStr := "*%s*\n_\\(tags:%v\\)_\n\\-\\-\\-\n`%s`"

	// create new
	if strings.HasPrefix(query, "+") {
		query = query[1:]
		newTopic := create(n.db, query)
		return fmt.Sprintf(formatStr, tools.EscapeMarkdownV2(newTopic.Title), newTopic.GetTagsString(), newTopic.Code), nil
	}

	// get all
	if strings.HasPrefix(query, "*") {
		var tags []Tag
		n.db.Find(&tags)
		tagsStr := make([]string, 0, len(tags))
		for _, tag := range tags {
			tagsStr = append(tagsStr, tag.Name)
		}
		sort.Strings(tagsStr)
		res := "all available tags:\n" + strings.Join(tagsStr, "\n")
		return res, nil
	}

	queryTags := strings.Split(strings.ToLower(query), " ")

	// get matched tags
	var tags []Tag
	n.db.Where("name IN ?", queryTags).Find(&tags)

	// get associated topics by tags
	var topics []Topic
	n.db.Model(&tags).Association("Topics").Find(&topics)

	// make set deduplicate topics
	set := make(map[string]Topic)
	for _, topic := range topics {
		set[topic.Title] = topic
	}

	// make slice for counting matches
	matches := make([]matchTopic, 0, len(set))
	for _, topic := range set {
		// add tags to topic struct
		n.db.Model(&topic).Association("Tags").Find(&topic.Tags)
		match := 0
		for _, tag := range topic.Tags {
			for _, queryTag := range queryTags {
				if tag.Name == queryTag {
					match++
				}
			}
		}
		matches = append(matches, matchTopic{match: match, topic: topic})
	}

	// sort slice by matches desc
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].match > matches[j].match
	})

	res := ""
	for i, val := range matches {
		res += "\n\n\\=\\=\\=\n" + fmt.Sprintf(formatStr, tools.EscapeMarkdownV2(val.topic.Title), val.topic.GetTagsString(), val.topic.Code)
		if i > 3 {
			break
		}
		if len(res) > 4096 {
			break
		}
	}

	return res, nil
}

type matchTopic struct {
	match int
	topic Topic
}
