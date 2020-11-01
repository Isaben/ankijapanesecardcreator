**Disclaimer: Learning project, use with caution**

# Japanese flashcard creator for Anki

This is a really simple (and probably broken, depending on **when** you're reading this) program made in Go using the Fyne package for the GUI. It's ugly and barely functional (well, so far it is working), but it does what I was trying to do, so that's a win I guess.

## Running

You'll need:

- Internet connection, as it uses [Jisho](https://jisho.org/) and [KanjiAPI](https://kanjiapi.dev/) to get the card contents.
- A font file capable of displaying japanese characters. Unfortunately Fyne defaults to a font not capable of doing so for some unknown reason. There is a free to use available in this repository.
- [AnkiConnect](https://foosoft.net/projects/anki-connect/) installed on your Anki.
- Patience I guess.

After downloading the repo, just run a `go build` and it'll generate a runnable binary for you. Fyne is supposedly compatible with Windows, Mac OS and Linux, so it should work just as expected. I haven't tested it with anything but Windows though.

## Motivation

I'm learning japanese for a while now, and at first I was using the amazing [Yomichan](https://foosoft.net/projects/yomichan/) extension with the AnkiConnect integration to populate my mining Anki decks. But unfortunately it doesn't provide the kanji readings on the back of the cards it creates, which prompted me to start creating my own cards manually, using Jisho. Of course, that turned into a chore faster than expected, and now here I am, trying to make something to automate the whole process.

Also, as a note: I'm fully aware that Go isn't the right tool to make a GUI application. This is a learning project, I went with Go because I wanted to use the language for something, and this was the only thing on my mind at the time.