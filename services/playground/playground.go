package playground

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"tbot/config"
)

type Playground struct {
	playgroundStg *config.PlaygroundSettings
}

func NewPlayground(playgroundStg *config.PlaygroundSettings) *Playground {
	return &Playground{playgroundStg: playgroundStg}
}

func (p *Playground) GetTag() string {
	return p.playgroundStg.Tag
}

func (p *Playground) Query(query string) (string, error) {
	client := &http.Client{}

	payload := "version=2&body=" + url.QueryEscape(query) + "&withVet=true"
	data := []byte(payload)

	log.Printf("url:%s\n", p.playgroundStg.Url)
	log.Printf("data:%s\n", string(data))
	req, err := http.NewRequest("POST", p.playgroundStg.Url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Add("Host", "play.golang.org")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Header.Add("Origin", "https://play.golang.org")
	req.Header.Add("Referer", "https://play.golang.org/")

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

	pgr := PlaygroundResponse{}
	err = json.Unmarshal(body, &pgr)
	if err != nil {
		return "", err
	}

	log.Printf("body:%s\n", string(body))
	//title := gjson.Get(string(body), "Events.0.Object.Message").String()
	//if title == "" {
	//	return "", fmt.Errorf("empty title")
	//}

	output := ""
	for _, out := range pgr.Events {
		output += out.Message
	}

	return output, nil
}

type PlaygroundResponse struct {
	Errors string `json:"Errors"`
	Events []struct {
		Message string `json:"Message"`
		Kind    string `json:"Kind"`
		Delay   int    `json:"Delay"`
	} `json:"Events"`
	Status      int  `json:"Status"`
	IsTest      bool `json:"IsTest"`
	TestsFailed int  `json:"TestsFailed"`
	VetOK       bool `json:"VetOK"`
}
