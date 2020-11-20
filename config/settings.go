package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Settings struct {
	HelpText  string          `yaml:"help_text"`
	WikiStg   *WikiSettings   `yaml:"wiki"`
	NewtonStg *NewtonSettings `yaml:"newton"`

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
  using: w [or W, в, В] <some_name>
  example: "w go"

  2) newton - powerful math calculator
  using: n [or N, н, Н] <operation> <expression>
  example: "n derive x^2+2x"`,
		WikiStg: &WikiSettings{
			Tag: "wiki",
			Url: "https://ru.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exintro=&explaintext=&titles=",
		},
		NewtonStg: &NewtonSettings{
			Tag: "newton",
			Url: "https://newton.now.sh/api/v2/",
		},
		Services: map[string][]string{
			"wiki":   {"w ", "W ", "wiki ", "Wiki ", "в ", "В", "вики ", "Вики "},
			"newton": {"n ", "N ", "newton ", "Newton ", "н ", "Н", "ньютон ", "Ньютон "},
		},
	}
}
