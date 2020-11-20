package wiki

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

//https://ru.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exintro=&explaintext=&titles=java

type Wiki struct {
	url string
}

func NewWiki(url string) *Wiki {
	return &Wiki{url: url}
}

func (w *Wiki) Query(query string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", w.url, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	title := gjson.Get(string(body), "query.pages.*.extract").String()
	if title == "" {
		return "", fmt.Errorf("empty title")
	}

	return title, nil
}
