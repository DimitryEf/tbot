package newton

import (
	"fmt"
	"github.com/tidwall/gjson"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tbot/config"
)

type Newton struct {
	newtonStg *config.NewtonSettings
	ready     bool
}

func NewNewton(newtonStg *config.NewtonSettings) *Newton {
	return &Newton{
		newtonStg: newtonStg,
		ready:     true,
	}
}

func (n *Newton) GetTag() string {
	return n.newtonStg.Tag
}

func (n *Newton) IsReady() bool {
	return n.ready
}

func (n *Newton) Query(query string) (string, error) {
	client := &http.Client{}
	query = strings.ToLower(query)

	operation := ""
	for _, op := range n.newtonStg.Operations {
		if strings.HasPrefix(query, op) {
			operation = op
			break
		}
	}
	if operation == "" {
		return "", fmt.Errorf("message has not some operation")
	}

	expression := query[len(operation):]
	expression = template.URLQueryEscaper(expression)
	if expression == "" {
		return "", fmt.Errorf("message has not some expression")
	}

	url := fmt.Sprintf("%s/%s/%s", n.newtonStg.Url, operation, expression)
	log.Printf("url:%s\n", url)
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

	title := gjson.Get(string(body), "result").String()
	if title == "" {
		return "", fmt.Errorf("empty result")
	}

	return title, nil
}
