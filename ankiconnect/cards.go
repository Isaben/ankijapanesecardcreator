package ankiconnect

import (
	"ankijapanesecardcreator/data"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// Options of a card
type Options struct {
	AllowDuplicate bool
	DuplicateScope string
}

// Fields of a card
type Fields struct {
	Front string
	Back  string
}

// Card structure from AnkiConnect
type Card struct {
	DeckName  string
	ModelName string
	Fields    Fields
	Options   Options
	Tags      []string
}

// UserInput what we want from the User
type UserInput struct {
	Term     string
	Sentence string
	DeckName string
}

type requestResult struct {
	Result int
	Error  interface{}
}

const ankiConnectAddress = "http://localhost:8765"

// CreateCard using AnkiConnect
func CreateCard(infos UserInput) (Card, error) {
	var returnedCard Card

	removeAllButKanjiAndHiragana := regexp.MustCompile(`[^\p{Han}\p{Hiragana}\p{Katakana}]+`)
	infos.Term = removeAllButKanjiAndHiragana.ReplaceAllString(infos.Term, "")

	jishoData, err := data.GetTermData(infos.Term)

	if err != nil {
		return returnedCard, err
	}

	if len(jishoData.Data) == 0 {
		return returnedCard, errors.New("No data found for the requested input. Try again with something else")
	}

	kanjiData, err := data.GetKanjisInfo(infos.Term)

	if err != nil {
		return returnedCard, err
	}

	valueThatMatters := jishoData.Data[0]
	fields := Fields{
		Front: infos.Term + "<br><br>" + infos.Sentence,
		Back:  "",
	}

	// add readings
	for _, value := range valueThatMatters.Japanese {
		fields.Back += value.Word + "【" + value.Reading + "】 "
	}
	fields.Back += "<br><br>"

	// add meanings
	for _, value := range valueThatMatters.Senses {
		speech := strings.Join(value.PartsOfSpeech, ", ")
		meaning := strings.Join(value.EnglishDefinitions, ", ")
		info := ""

		if len(value.Info) > 0 {
			info = strings.Join(value.Info, ", ")
		}

		if !strings.Contains(fields.Back, speech) {
			fields.Back += speech + "<br>"
		}

		fields.Back += "• " + meaning + " " + info + "<br>"
	}
	fields.Back += "<br><br>"

	// add kanji infos
	for _, value := range kanjiData {
		fields.Back += value.Symbol + "<br>"
		fields.Back += strings.Join(value.Meanings, ", ") + "<br>"
		fields.Back += "Kun: " + strings.Join(value.KunReadings, ", ") + "<br>"
		fields.Back += "On: " + strings.Join(value.OnReadings, ", ") + "<br>"
		fields.Back += "Name: " + strings.Join(value.NameReadings, ", ") + "<br><br>"
	}

	returnedCard = Card{
		DeckName:  infos.DeckName,
		ModelName: "Basic",
		Fields:    fields,
		Options: Options{
			AllowDuplicate: false,
			DuplicateScope: "deck",
		},
		Tags: []string{"AnkiCreatorSoftware"},
	}

	return returnedCard, nil
}

// AddCardToDeck using AnkiConnect
func AddCardToDeck(card Card) error {
	postBody, _ := json.Marshal(map[string]interface{}{
		"action":  "addNote",
		"version": 6,
		"params": map[string]interface{}{
			"note": map[string]interface{}{
				"deckName":  card.DeckName,
				"modelName": card.ModelName,
				"fields": map[string]interface{}{
					"Front": card.Fields.Front,
					"Back":  card.Fields.Back,
				},
				"options": map[string]interface{}{
					"allowDuplicate": card.Options.AllowDuplicate,
					"duplicateScope": card.Options.DuplicateScope,
				},
				"tags": card.Tags,
			},
		},
	})
	responseBody := bytes.NewBuffer(postBody)

	res, err := http.Post(ankiConnectAddress, "application/json", responseBody)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var resBody requestResult
	json.Unmarshal(body, &resBody)

	if resBody.Error != nil {
		if value, isString := resBody.Error.(string); isString {
			return errors.New(value)
		}
		return errors.New(string(body))
	}

	return nil
}
