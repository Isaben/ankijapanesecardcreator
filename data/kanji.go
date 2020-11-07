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

type channelResult struct {
	Kanji Kanji
	Error error
}

const kanjiAPIPath = "https://kanjiapi.dev/v1"

func fetchFromAPI(kanji rune, channel chan channelResult) {
	var resultKanji Kanji
	stringified := string(kanji)
	getResponse, err := http.Get(kanjiAPIPath + "/kanji/" + url.QueryEscape(stringified))

	if err != nil {
		channel <- channelResult{
			resultKanji,
			err,
		}
		return
	}

	defer getResponse.Body.Close()
	body, err := ioutil.ReadAll(getResponse.Body)

	if err != nil {
		channel <- channelResult{
			resultKanji,
			err,
		}
		return
	}

	err = json.Unmarshal(body, &resultKanji)

	if err != nil {
		channel <- channelResult{
			resultKanji,
			err,
		}
		return
	}

	channel <- channelResult{
		resultKanji,
		nil,
	}
}

// GetKanjisInfo fetch kanji readings and meanings from kanjiapi.dev
func GetKanjisInfo(term string) ([]Kanji, error) {
	kanjisInfo := make([]Kanji, 0, 255)
	if len(term) > 255 {
		return kanjisInfo, errors.New("Stop")
	}

	removeAllButKanji := regexp.MustCompile(`[^\p{Han}]+`)
	term = removeAllButKanji.ReplaceAllString(term, "")
	channelBuffer := len(term) / 3
	channel := make(chan channelResult, channelBuffer)

	for _, char := range term {
		go fetchFromAPI(char, channel)
	}

	for finalizedRequests := 0; finalizedRequests < channelBuffer; {
		select {
		case result := <-channel:
			if result.Error != nil {
				return kanjisInfo, result.Error
			}

			kanjisInfo = append(kanjisInfo, result.Kanji)
			finalizedRequests++
		}
	}

	return kanjisInfo, nil
}
