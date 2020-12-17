package golang

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"tbot/tools"
)

type Topic struct {
	gorm.Model
	Id      int    `gorm:"primaryKey"`
	Title   string `gorm:"unique"`
	Code    string
	Checked bool
	Tags    []Tag `gorm:"many2many:topic_tags;"`
}

//type TopicTags struct {
//	gorm.Model
//	Id int `gorm:"primaryKey"`
//	Topics []Topic `gorm:"many2many:topic_tags;"`
//	Tags    []Tag `gorm:"many2many:topic_tags;"`
//}

type Tag struct {
	gorm.Model
	Id     int     `gorm:"primaryKey"`
	Name   string  `gorm:"unique"`
	Topics []Topic `gorm:"many2many:topic_tags;"`
}

func ConvertQueryToTopic(query string) (*Topic, error) {
	topic := Topic{}
	divider := "---"
	if strings.Index(query, "\n") == -1 {
		return nil, fmt.Errorf("query '%s' has no symbol '\\n'", query)
	}
	title := query[:strings.Index(query, "\n")]
	topic.Title = title
	query = query[len(title):]
	if strings.Index(query, divider) == -1 {
		return nil, fmt.Errorf("query '%s' has no symbol '%s'", query, divider)
	}
	tagsStr := query[:strings.Index(query, divider)]
	if len("(tags:  ") > len(tagsStr) {
		return nil, fmt.Errorf("len(\"(tags:  \") > len(tagsStr)")
	}
	if strings.Index(tagsStr, ")") == -1 {
		return nil, fmt.Errorf("query '%s' has no symbol ')'", query)
	}
	tagsStr2 := tagsStr[len("(tags:  "):strings.Index(tagsStr, ")")]
	tags := strings.Split(strings.ToLower(tagsStr2), " ")
	for _, tag := range tags {
		topic.Tags = append(topic.Tags, Tag{
			Name: tag,
		})
	}
	if len(tagsStr) > len(query) {
		return nil, fmt.Errorf("len(tagsStr) > len(query)")
	}
	query = query[len(tagsStr):]
	if len(divider) > len(query) {
		return nil, fmt.Errorf("len(divider) > len(query)")
	}
	code := query[len(divider):]
	code = strings.TrimPrefix(code, "\n")

	code = tools.EscapeMarkdownV2(code)

	topic.Code = code
	return &topic, nil
}

func (t Topic) GetTagsString() string {
	res := ""
	for _, tag := range t.Tags {
		res += " " + tag.Name
	}
	return res
}
