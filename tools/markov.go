package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"tbot/internal/errors"
)

type State struct {
	Id         int
	Word       string
	Count      int
	Prob       float64
	NextStates []State
}

func main() {
	states := initialize("platon.txt")

	//os.Create("platon.db")

	db, err := sql.Open("sqlite3", "platon.db")
	errors.PanicIfErr(err)
	defer db.Close()

	psStates, err := db.Prepare("insert into states (id, word, counts, prob) values ($1, $2, $3, $4)")
	errors.PanicIfErr(err)
	psNextStates, err := db.Prepare("insert into next_states (id, next_id) values ($1, $2)")
	errors.PanicIfErr(err)
	for _, state := range states {
		_, err = psStates.Exec(state.Id, state.Word, state.Count, state.Prob)
		errors.PanicIfErr(err)
		for _, nextState := range state.NextStates {
			_, err = psStates.Exec(nextState.Id, nextState.Word, nextState.Count, nextState.Prob)
			_, err = psNextStates.Exec(state.Id, nextState.Id)
		}
	}

	//fmt.Println(states)

}

func initialize(file string) []State {
	text := readTxt(file)
	return train(text)
}

func readTxt(path string) string {
	data, err := ioutil.ReadFile(path)
	errors.PanicIfErr(err)
	dataClean := strings.Replace(string(data), "\n", " ", -1)
	dataClean = strings.Replace(dataClean, ".", "", -1)
	dataClean = strings.Replace(dataClean, ",", "", -1)
	dataClean = strings.Replace(dataClean, "-", "", -1)
	dataClean = strings.Replace(dataClean, "!", "", -1)
	dataClean = strings.Replace(dataClean, "?", "", -1)
	dataClean = strings.Replace(dataClean, "[", "", -1)
	dataClean = strings.Replace(dataClean, "]", "", -1)
	dataClean = strings.ToLower(dataClean)
	dataClean = strings.Replace(dataClean, "бог", "Бог", -1)
	return dataClean
}

func printLoading(n int, total int) {
	var bar []string
	tantPerFourty := int((float64(n) / float64(total)) * 40)
	tantPerCent := int((float64(n) / float64(total)) * 100)
	for i := 0; i < tantPerFourty; i++ {
		bar = append(bar, "█")
	}
	progressBar := strings.Join(bar, "")
	fmt.Printf("\r " + progressBar + " - " + strconv.Itoa(tantPerCent) + "")
}

var id int

func addWordToStates(states []State, word string) ([]State, int) {
	iState := -1
	for i := 0; i < len(states); i++ {
		if states[i].Word == word {
			iState = i
		}
	}
	if iState >= 0 {
		states[iState].Count++
	} else {
		var tempState State
		tempState.Word = word
		tempState.Count = 1

		id++
		tempState.Id = id

		states = append(states, tempState)
		iState = len(states) - 1

	}
	return states, iState
}

func calcMarkovStates(words []string) []State {
	var states []State
	//count words
	for i := 0; i < len(words)-1; i++ {
		var iState int
		states, iState = addWordToStates(states, words[i])
		if iState < len(words) {
			states[iState].NextStates, _ = addWordToStates(states[iState].NextStates, words[i+1])
		}

		printLoading(i, len(words))
	}

	//count prob
	for i := 0; i < len(states); i++ {
		states[i].Prob = (float64(states[i].Count) / float64(len(words)) * 100)
		for j := 0; j < len(states[i].NextStates); j++ {
			states[i].NextStates[j].Prob = (float64(states[i].NextStates[j].Count) / float64(len(words)) * 100)
		}
	}
	fmt.Println("\ntotal words computed: " + strconv.Itoa(len(words)))
	//fmt.Println(states)
	return states
}

func textToWords(text string) []string {
	s := strings.Split(text, " ")
	return s
	words := make([]string, 0, len(s))
	patternDigits := *regexp.MustCompile(`[\d]`)
	patternEngWord := *regexp.MustCompile(`[a-zA-Z]`)
	for _, word := range s {
		if word == "" || patternDigits.Match([]byte(word)) || patternEngWord.Match([]byte(word)) {
			continue
		}
		words = append(words, word)
	}
	return words
}

func train(text string) []State {

	words := textToWords(text)
	states := calcMarkovStates(words)
	//fmt.Println(states)

	return states
}

//-----------------

func getNextMarkovState(states []State, word string) string {
	iState := -1
	for i := 0; i < len(states); i++ {
		if states[i].Word == word {
			iState = i
		}
	}
	if iState < 0 {
		return "word no exist on the memory"
	}
	var next State
	next = states[iState].NextStates[0]
	next.Prob = rand.Float64() * states[iState].Prob
	for i := 0; i < len(states[iState].NextStates); i++ {
		if (rand.Float64()*states[iState].NextStates[i].Prob) > next.Prob && states[iState-1].Word != states[iState].NextStates[i].Word {
			next = states[iState].NextStates[i]
		}
	}
	return next.Word
}

func generateText(states []State, initWord string, count int) string {
	var generatedText []string
	word := initWord
	generatedText = append(generatedText, word)
	for i := 0; i < count; i++ {
		word = getNextMarkovState(states, word)
		if word == "word no exist on the memory" {
			return "word no exist on the memory"
		}
		generatedText = append(generatedText, word)
	}
	//generatedText = append(generatedText, ".")
	text := strings.Join(generatedText, " ")
	return text
}
