package golang

import (
	"gorm.io/gorm"
	"strings"
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

func ConvertQueryToTopic(query string) *Topic {
	topic := Topic{}
	title := query[:strings.Index(query, "\n")]
	topic.Title = title
	query = query[len(title):]
	tagsStr := query[:strings.Index(query, "---")]
	tagsStr2 := tagsStr[len("(tags:  "):strings.Index(tagsStr, ")")]
	tags := strings.Split(strings.ToLower(tagsStr2), " ")
	for _, tag := range tags {
		topic.Tags = append(topic.Tags, Tag{
			Name: tag,
		})
	}
	query = query[len(tagsStr):]
	code := query[5:]
	topic.Code = code
	return &topic
}

func (t Topic) GetTagsString() string {
	res := ""
	for _, tag := range t.Tags {
		res += " " + tag.Name
	}
	return res
}
