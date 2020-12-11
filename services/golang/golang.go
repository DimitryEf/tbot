package golang

import (
	"fmt"
	"gorm.io/gorm"
	"sort"
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

	//initialize(db)

	return &Golang{
		golangStg: golangStg,
		db:        db,
		ready:     true,
	}
}

func initialize(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&Topic{}, &Tag{})
	errors.PanicIfErr(err)

	create(db, "Get executable dir\n(tags: get executable dir)\n---\n\nex, err := os.Executable()\ndir := filepath.Dir(ex)\nfmt.Println(\"dir:\", dir)\n")
	create(db, "Extract beginning of string (prefix)\n(tags: extract beginning string prefix)\n---\n\nt := string([]rune(s)[:5])")
	create(db, "Extract string suffix\n(tags: extract string suffix)\n---\n\nt := string([]rune(s)[len([]rune(s))-5:])")
	create(db, "Exec other program\n(tags: exec program)\n---\n\nerr := exec.Command(\"program\", \"arg1\", \"arg2\").Run()")
	create(db, "Telegram message markdown\n(tags: telegram message markdown)\n---\n\n*полужирный*\n_курсив_\n[ссылка](http://www.example.com/)\n`строчный моноширинный`\n```text\nблочный моноширинный (можно писать код)\n```\n\n[https://github.com/go-telegram-bot-api/telegram-bot-api](import \"github.com/go-telegram-bot-api/telegram-bot-api\")\n\nmsg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)\nmsg.ParseMode = \"markdown\" //msg.ParseMode = tgbotapi.ModeMarkdown")
	create(db, "Telegram message html\n(tags: telegram message html)\n---\n\n<b>полужирный</b>, <strong>полужирный</strong>\n<i>курсив</i>\n<a href=\"http://www.example.com/\">ссылка</a>\n<code>строчный моноширинный</code>\n<pre>блочный моноширинный (можно писать код)</pre>\n\n[https://github.com/go-telegram-bot-api/telegram-bot-api](import \"github.com/go-telegram-bot-api/telegram-bot-api\")\n\nmsg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)\nmsg.ParseMode = \"HTML\" //msg.ParseMode = tgbotapi.ModeHTML")
	create(db, "Iterate over map entries ordered by keys\n(tags: iterate map order key)\n---\n\nkeys := make([]string, 0, len(mymap))\nfor k := range mymap {\n        keys = append(keys, k)\n}\nsort.Strings(keys)\nfor _, k := range keys {\n        x := mymap[k]\n        fmt.Println(\"Key =\", k, \", Value =\", x)\n}\n")
	create(db, "Iterate over map entries ordered by values\n(tags: iterate map order value)\n---\n\ntype entry struct {\n        key   string\n        value int\n}\nentries := make([]entry, 0, len(mymap))\nfor k, x := range mymap {\n        entries = append(entries, entry{key: k, value: x})\n}\nsort.Slice(entries, func(i, j int) bool {\n        return entries[i].value < entries[j].value\n})\nfor _, e := range entries {\n        fmt.Println(\"Key =\", e.key, \", Value =\", e.value)\n}")
	create(db, "Slice to set\n(tags: slice set)\n---\n\ny := make(map[T]struct{}, len(x))\nfor _, v := range x {\n        y[v] = struct{}{}\n}")
	create(db, "Deduplicate slice\n(tags: deduplicate slice remove duplicate)\n---\n\nseen := make(map[T]bool)\nj := 0\nfor _, v := range x {\n        if !seen[v] {\n                x[j] = v\n                j++\n                seen[v] = true\n        }\n}\nfor i := j; i < len(x); i++ {\n        x[i] = nil\n}\nx = x[:j]")
	create(db, "Shuffle a slice\n(tags: slice shuffle)\n---\n\ny := make(map[T]struct{}, len(x))\nfor _, v := range x {\n        y[v] = struct{}{}\n}")
	create(db, "Sort slice asc\n(tags: sort slice asc)\n---\n\nsort.Slice(items, func(i, j int) bool {\n        return items[i].p < items[j].p\n})")
	create(db, "Sort slice desc\n(tags: sort slice desc)\n---\n\nsort.Slice(items, func(i, j int) bool {\n        return items[i].p > items[j].p\n})")
	create(db, "Remove item from slice by index\n(tags: remove item slice index)\n---\n\nitems = append(items[:i], items[i+1:]...)")
	create(db, "Graph with adjacency lists\n(tags: graph struct)\n---\n\ntype Vertex struct{\n        Id int\n        Label string\n        Neighbours map[*Vertex]bool\n}\ntype Graph []*Vertex")
	create(db, "Reverse a string\n(tags: string reverse)\n---\n\nrunes := []rune(s)\nfor i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {\n   runes[i], runes[j] = runes[j], runes[i]\n}\nt := string(runes)")
	create(db, "Insert item in slice\n(tags: insert item slice)\n---\n\ns = append(s, 0)\ncopy(s[i+1:], s[i:])\ns[i] = x")
	create(db, "Filter slice\n(tags: filter slice)\n---\n\ny := make([]T, 0, len(x))\nfor _, v := range x{\n        if p(v){\n                y = append(y, v)\n        }\n}")
	create(db, "File content to string\n(tags: file content string)\n---\n\nb, err := ioutil.ReadFile(f)\nlines := string(b)")

}

func create(db *gorm.DB, query string) Topic {
	topic := ConvertQueryToTopic(query)
	tags := topic.Tags
	for i, tag := range tags {
		db.Debug().Where("name = ?", tag.Name).Find(&tag)
		if tag.Id == 0 {
			db.Debug().Create(&tag)
		}
		tags[i] = tag
	}
	topic.Tags = tags
	db.Debug().Create(topic)
	return *topic
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
		newTopic := create(n.db, query)
		return fmt.Sprintf(formatStr, newTopic.Title, newTopic.GetTagsString(), newTopic.Code), nil
	}

	if strings.HasPrefix(query, "*") {
		var topics []Topic
		n.db.Find(&topics)
		//var topic Topic
		n.db.Model(&topics).Association("Tags").Find(&topics)
		res := ""
		for _, topic := range topics {
			n.db.Model(&topic).Association("Tags").Find(&topic.Tags)
			res += "\n\n===\n" + fmt.Sprintf(formatStr, topic.Title, topic.GetTagsString(), topic.Code)
		}
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
	for _, val := range matches {
		res += "\n\n===\n" + fmt.Sprintf(formatStr, val.topic.Title, val.topic.GetTagsString(), val.topic.Code)
	}

	return res, nil
}

type matchTopic struct {
	match int
	topic Topic
}
