package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Settings struct {
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
