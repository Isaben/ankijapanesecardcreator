package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"ankijapanesecardcreator/ankiconnect"
	"net/http"
)

var isAnkiConnected bool

func isDeckNameValid(deckName string, decks []string) bool {
	found := false

	for _, element := range decks {
		if deckName == element {
			found = true
			break
		}
	}

	return found
}

func updateDeckList(combo *widget.Select, label *widget.Label, client *http.Client) {
	for true {
		time.Sleep(3 * time.Second)
		list, err := ankiconnect.GetDecks(client)

		if err != nil {
			isAnkiConnected = false
			label.Show()
			continue
		}
		isAnkiConnected = true
		label.Hide()
		combo.Options = list.Result
		selectedValue := combo.Selected

		if !isDeckNameValid(selectedValue, list.Result) {
			combo.Selected = combo.PlaceHolder
		}
	}
}

func main() {
	directory, _ := os.Open("./")
	files, _ := directory.Readdirnames(0)
	fontInDir := false

	for _, file := range files {
		if file == "font.ttf" {
			fontInDir = true
			os.Setenv("FYNE_FONT", "./font.ttf")
			break
		}
	}

	a := app.New()
	w := a.NewWindow("Anki Japanese Card Creator")
	w.SetMaster()

	transport := &http.Transport{
		MaxIdleConns:    0,
		IdleConnTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: transport}

	word := widget.NewLabel("Word: ")
	sentence := widget.NewLabel("Phrase: ")
	deck := widget.NewLabel("Deck: ")
	ankiDisconnectedLabel := widget.NewLabel("Anki is not connected!")
	ankiDisconnectedLabel.TextStyle.Bold = true
	ankiDisconnectedLabel.Hide()
	wordField := widget.NewEntry()
	wordFieldContainer := widget.NewHScrollContainer(wordField)
	sentenceField := widget.NewEntry()
	sentenceFieldContainer := widget.NewHScrollContainer(sentenceField)

	logs := widget.NewLabel("")
	logsLabel := widget.NewLabel("Logs: ")
	logsContainer := widget.NewVScrollContainer(logs)
	logsContainer.SetMinSize(fyne.NewSize(800, 200))

	selectedDeck := ""

	combo := widget.NewSelect([]string{}, func(value string) {
		selectedDeck = value
	})
	addCardButton := widget.NewButton("Add", func() {
		if !isAnkiConnected {
			msg := "Anki is not connected! Verify your AnkiConnect settings, and a try again"
			logs.Text += msg + "\n"
			logs.Refresh()
			dialog.ShowError(errors.New(msg), w)
			return
		}

		if !isDeckNameValid(selectedDeck, combo.Options) {
			msg := "No deck selected! Please choose one and try again"
			logs.Text += msg + "\n"
			logs.Refresh()
			dialog.ShowError(errors.New(msg), w)
			return
		}

		card, err := ankiconnect.CreateCard(ankiconnect.UserInput{
			Term:     wordField.Text,
			Sentence: sentenceField.Text,
			DeckName: selectedDeck,
		}, client)

		if err != nil {
			logs.Text += err.Error() + "\n"
			logs.Refresh()
			logsContainer.ScrollToBottom()
			dialog.ShowError(err, w)
			return
		}

		errAdd := ankiconnect.AddCardToDeck(card, client)

		if errAdd != nil {
			logs.Text += errAdd.Error() + "\n"
			logs.Refresh()
			logsContainer.ScrollToBottom()
			dialog.ShowError(errAdd, w)
			return
		}

		logs.Text += "Term " + wordField.Text + " added succesfully to the " + selectedDeck + "\n"
		logs.Refresh()
		logsContainer.ScrollToBottom()
	})

	grid := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		word,
		wordFieldContainer,
		sentence,
		sentenceFieldContainer,
		addCardButton,
		deck,
		combo,
		ankiDisconnectedLabel,
		logsLabel,
		logsContainer,
	)

	w.SetContent(grid)

	w.Resize(fyne.NewSize(800, 600))

	decks, err := ankiconnect.GetDecks(client)

	if err != nil {
		fmt.Println(err)
		dialog.ShowError(err, w)
		isAnkiConnected = false
	} else {
		combo.Options = decks.Result
		isAnkiConnected = true
	}

	//go updateDeckList(combo, ankiDisconnectedLabel, client)

	if !fontInDir {
		dialog.ShowError(errors.New("font.ttf file not found! You won't be seeing japanese characters without one"), w)
	}
	w.ShowAndRun()
}
