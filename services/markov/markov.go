package markov

import (
	"fmt"
	"strconv"
	"strings"
	"tbot/config"
)

type Markov struct {
	markovStg *config.MarkovSettings
	states    []State
	ready     bool
}

func NewMarkov(markovStg *config.MarkovSettings) *Markov {
	m := &Markov{
		markovStg: markovStg,
		ready:     false,
	}

	go m.load()

	//states := initialize(markovStg.File)
	return m
}

func (m *Markov) GetTag() string {
	return m.markovStg.Tag
}

func (m *Markov) IsReady() bool {
	return m.ready
}

func (m *Markov) Query(query string) (string, error) {
	wordAndCount := strings.Split(query, " ")
	word := wordAndCount[0]
	if len(wordAndCount) > 2 {
		return "", fmt.Errorf("there are needs one word or one word with count. Example: человек 10")
	}
	count := 5
	if len(wordAndCount) == 2 {
		var err error
		count, err = strconv.Atoi(wordAndCount[1])
		if err != nil {
			return "", fmt.Errorf("second parameter is not a digit. Example: человек 10")
		}

	}
	result := generateText(m.states, word, count)
	return result, nil

}
