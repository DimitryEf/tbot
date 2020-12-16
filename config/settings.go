package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Settings struct {
	HelpText      string              `yaml:"help_text"`
	WikiStg       *WikiSettings       `yaml:"wiki"`
	NewtonStg     *NewtonSettings     `yaml:"newton"`
	PlaygroundStg *PlaygroundSettings `yaml:"playground"`
	MarkovStg     *MarkovSettings     `yaml:"markov"`
	GolangStg     *GolangSettings     `yaml:"golang"`

	Services map[string][]string `yaml:"services"`
}

type WikiSettings struct {
	Tag string `yaml:"tag"`
	Url string `yaml:"url"`
}

type NewtonSettings struct {
	Tag        string   `yaml:"tag"`
	Url        string   `yaml:"url"`
	Operations []string `yaml:"operations"`
}

type PlaygroundSettings struct {
	Tag string `yaml:"tag"`
	Url string `yaml:"url"`
}

type MarkovSettings struct {
	Tag  string `yaml:"tag"`
	File string `yaml:"file"`
}

type GolangSettings struct {
	Tag  string `yaml:"tag"`
	File string `yaml:"file"`
}

// LoadFromFile create configuration from file.
func NewSettings(fileName string) (*Settings, error) {
	stg := Settings{}

	configBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		return getDefaultSettings(), nil
	}

	err = yaml.Unmarshal(configBytes, &stg)
	if err != nil {
		return &stg, err
	}

	return &stg, nil
}

func getDefaultSettings() *Settings {
	return &Settings{
		HelpText: `
  1) wiki - find some title in wikipedia
  using: w [or W, в, В] <name>
  example: "w go"

  2) newton - powerful math calculator
  using: n [or N, н, Н] <operation> <expression>
  example: "n derive x^2+2x"
  list of available operations:
    simplify
    factor
    derive
    integrate
    zeroes
    tangent
    area
    cos
    sin
    tan
    arccos
    arcsin
    arctan
    abs
    log

  3) playground - write, compile and run go code
  using: p [or P, п, П] <code>
  example: p package main

             import (
               "fmt"
             )

             func main() {
               fmt.Println("Hello, World!")
             }

  4) markov - generates a sentence from the works of Plato using the Markov algorithm
  using: m [or M, м, М] <word> <count>
  example: m человек 10

  5) golang - tips and tricks
  using: go [or Go, го, Го] <tag1> <tag2> ...
  example: go extract prefix`,
		WikiStg: &WikiSettings{
			Tag: "wiki",
			Url: "https://ru.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exintro=&explaintext=&titles=",
		},
		NewtonStg: &NewtonSettings{
			Tag: "newton",
			Url: "https://newton.now.sh/api/v2/",
		},
		PlaygroundStg: &PlaygroundSettings{
			Tag: "playground",
			Url: "https://play.golang.org/compile",
		},
		MarkovStg: &MarkovSettings{
			Tag:  "markov",
			File: "services/markov/platon.txt",
		},
		GolangStg: &GolangSettings{
			Tag:  "golang",
			File: "databases/golang.db",
		},
		Services: map[string][]string{
			"wiki":       {"w ", "W ", "wiki ", "Wiki ", "в ", "В", "вики ", "Вики "},
			"newton":     {"n ", "N ", "newton ", "Newton ", "н ", "Н", "ньютон ", "Ньютон "},
			"playground": {"p ", "P ", "play ", "Play ", "п ", "П ", "плэй ", "Плэй "},
			"markov":     {"m ", "M ", "markov ", "Markov ", "м ", "М ", "марков ", "Марков "},
			"golang":     {"go ", "Go ", "golang ", "Golang ", "го ", "Го ", "голанг ", "Голанг "},
		},
	}
}
