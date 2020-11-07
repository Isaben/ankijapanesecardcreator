package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

// Kanji struct, containing all data fetched from kanjiapi
type Kanji struct {
	Symbol       string   `json:"kanji"`
	Meanings     []string `json:"meanings"`
	KunReadings  []string `json:"kun_readings"`
	OnReadings   []string `json:"on_readings"`
	NameReadings []string `json:"name_readings"`
}

const kanjiAPIPath = "https://kanjiapi.dev/v1"

// GetKanjisInfo fetch kanji readings and meanings from kanjiapi.dev
func GetKanjisInfo(term string) ([]Kanji, error) {
	kanjisInfo := make([]Kanji, 0, 255)
	if len(term) > 255 {
		return kanjisInfo, errors.New("Stop")
	}

	removeAllButKanji := regexp.MustCompile(`[^\p{Han}]+`)
	term = removeAllButKanji.ReplaceAllString(term, "")

	for _, char := range term {
		stringified := string(char)
		getResponse, err := http.Get(kanjiAPIPath + "/kanji/" + url.QueryEscape(stringified))

		if err != nil {
			return kanjisInfo, err
		}

		defer getResponse.Body.Close()
		body, err := ioutil.ReadAll(getResponse.Body)

		if err != nil {
			return kanjisInfo, err
		}

		var thisKanjiInfo Kanji
		err = json.Unmarshal(body, &thisKanjiInfo)

		if err != nil {
			return kanjisInfo, err
		}

		kanjisInfo = append(kanjisInfo, thisKanjiInfo)
	}

	return kanjisInfo, nil
}
