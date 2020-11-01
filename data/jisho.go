package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const jishoAPIPath = "https://jisho.org/api/v1"

// Japanese infos from Jisho
type Japanese struct {
	Word    string
	Reading string
}

// Senses english definitions and grammar infos from Jisho
type Senses struct {
	EnglishDefinitions []string `json:"english_definitions"`
	PartsOfSpeech      []string `json:"parts_of_speech"`
	Info               []string `json:"info"`
}

// JishoData containing everything fetched from Jisho API
type JishoData struct {
	Senses   []Senses   `json:"senses"`
	Japanese []Japanese `json:"japanese"`
}

// APIResponse literally what the API returns
type APIResponse struct {
	Data []JishoData `json:"data"`
}

// GetTermData fetches everything about a term from Jisho API
func GetTermData(term string) (APIResponse, error) {
	var response APIResponse

	if len(term) == 0 {
		return response, errors.New("No input found")
	}

	getResponse, err := http.Get(jishoAPIPath + "/search/words?keyword=" + url.QueryEscape(term))

	if err != nil {
		return response, err
	}

	defer getResponse.Body.Close()
	body, err := ioutil.ReadAll(getResponse.Body)

	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	return response, err
}
