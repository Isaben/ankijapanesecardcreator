package ankiconnect

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Deck representando o retorno da requisição para decks
type Deck struct {
	Result     []string
	ErrorValue interface{} `json:"error"`
}

// GetDecks requisita a lista de decks via anki connect
func GetDecks(client *http.Client) (Deck, error) {
	var deckResult Deck
	postBody, _ := json.Marshal(map[string]interface{}{
		"action":  "deckNames",
		"version": 6,
	})
	responseBody := bytes.NewBuffer(postBody)

	client.CloseIdleConnections()
	res, err := client.Post(ankiConnectAddress, "application/json", responseBody)

	if err != nil {
		return deckResult, errors.New("Couldn't connect to Anki. Verify your AnkiConnect installation and try again")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return deckResult, errors.New("Something went wrong while reading the response body of AnkiConnect")
	}

	unmarshalErr := json.Unmarshal(body, &deckResult)

	if unmarshalErr != nil {
		return deckResult, errors.New("Something went wrong while unmarshalling the JSON")
	}

	return deckResult, nil
}
