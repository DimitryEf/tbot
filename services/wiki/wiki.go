package wiki

import (
	"compress/gzip"
	"fmt"
	"github.com/tidwall/gjson"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"tbot/config"
)

//https://ru.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exintro=&explaintext=&titles=java

type Wiki struct {
	wikiStg *config.WikiSettings
}

func NewWiki(wikiStg *config.WikiSettings) *Wiki {
	return &Wiki{wikiStg: wikiStg}
}

func (w *Wiki) GetTag() string {
	return w.wikiStg.Tag
}

func (w *Wiki) Query(query string) (string, error) {
	client := &http.Client{}

	//url := strings.ToLower(query)
	url := template.URLQueryEscaper(query)
	url = fmt.Sprintf("%s%s", w.wikiStg.Url, url)
	log.Printf("url:%s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	title := gjson.Get(string(body), "query.pages.*.extract").String()
	if title == "" {
		return "", fmt.Errorf("empty title")
	}

	return title, nil
}
