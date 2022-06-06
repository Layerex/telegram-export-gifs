package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/3bl3gamer/tgclient/mtproto"
	"github.com/adrg/xdg"
)

const programName = "telegram-export-gifs"
const sessionFile = "telegram-export-gifs/tg.session"

func (t *Telegram) GetCurrentGifs() ([]mtproto.TL_document, error) {
	tl := t.Request(mtproto.TL_messages_getSavedGifs{})
	savedGifsRes, ok := tl.(mtproto.TL_messages_savedGifs)
	if !ok {
		return nil, errors.New("TL_messages_getSavedGifs failed")
	}

	savedGifs := make([]mtproto.TL_document, len(savedGifsRes.Gifs))
	for i, gifDocument := range savedGifsRes.Gifs {
		savedGifs[i] = gifDocument.(mtproto.TL_document)
	}

	return savedGifs, nil
}

func (t *Telegram) SaveGif(document mtproto.TL_inputDocument, unsave bool) {
	t.Request(mtproto.TL_messages_saveGif{ID: document, Unsave: EncodeBool(unsave)})
}

func (t *Telegram) GetAllGifs(savedGifsLimit int, clearGifs bool) []mtproto.TL_document {
	savedGifs := make([]mtproto.TL_document, 0)
	unsavedGifs := make([]mtproto.TL_inputDocument, 0)

	defer func() {
		errorString := recover()
		if errorString != nil {
			fmt.Printf("Error when getting gifs: %s\n", errorString)
			defer panic(errorString)
		}
		if !clearGifs {
			for i := range unsavedGifs {
				inputDocument := unsavedGifs[len(unsavedGifs)-i-1]
				fmt.Printf("(%d/%d) Resaving gifs\n", i+1, len(unsavedGifs))
				t.SaveGif(inputDocument, false) // Resave gif
			}
		}
	}()

	// Handle interrupts
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	fmt.Println("WARNING: this program uses a hack to get older gifs not visible from client: it temporarily removes newer gifs to get older ones. If something goes wrong (for example if program gets killed, power goes off or internet disconnects), removed gifs will be lost.")
	for {
		currentSavedGifs, err := t.GetCurrentGifs()
		if err != nil {
			panic(err)
		}
		for _, document := range currentSavedGifs {
			savedGifs = append(savedGifs, document)
		}
		if len(currentSavedGifs) != savedGifsLimit { // if there is not max amount of gifs, there are no more
			fmt.Printf("Unsaved %d gifs\n", len(savedGifs)-len(currentSavedGifs))
			fmt.Printf("Got %d gifs\n", len(savedGifs))
			return savedGifs
		}
		for i, document := range currentSavedGifs {
			select {
			case <-interruptChan:
				return make([]mtproto.TL_document, 0)
			default:
				inputDocument := mtproto.TL_inputDocument{ID: document.ID, AccessHash: document.AccessHash, FileReference: document.FileReference}
				// fmt.Println(inputDocument)
				fmt.Printf("(%d/%d+) Unsaving gifs\n", len(savedGifs)-len(currentSavedGifs)+i+1, len(savedGifs))
				unsavedGifs = append(unsavedGifs, inputDocument)
				// Gifs can be lost that way, but i am too lazy to implement anything to recover them
				t.SaveGif(inputDocument, true) // Temporarily unsave gif, so old gifs become visible
			}
		}
	}
}

func main() {
	args := ParseArgs()

	var sessionFilePath string
	if !args.DontSaveSession {
		var err error
		sessionFilePath, err = xdg.DataFile(sessionFile)
		if err != nil {
			panic(err)
		}
	}

	var t Telegram
	err := t.SignIn(args.AppID, args.AppHash, sessionFilePath)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(args.Directory, 0755)
	if err != nil {
		panic(err)
	}
	err = os.Chdir(args.Directory)
	if err != nil {
		panic(err)
	}

	tl := t.Request(mtproto.TL_help_getConfig{})
	serverConfig, ok := tl.(mtproto.TL_config)
	if !ok {
		panic("TL_help_getConfig failed")
	}

	savedGifs := t.GetAllGifs(int(serverConfig.SavedGifsLimit), false)

	for i, document := range savedGifs {
		filename := strconv.FormatInt(document.ID, 10) + ".mp4"

		fileInfo, err := os.Stat(filename)
		exists := !errors.Is(err, os.ErrNotExist)
		if exists && fileInfo.Size() == int64(document.Size) {
			fmt.Printf("(%d/%d) Gif %s already exported\n", i+1, len(savedGifs), filename)
		} else {
			fmt.Printf("(%d/%d) Exporting gif %s\n", i+1, len(savedGifs), filename)
			err := t.DownloadDocument(filename, document)
			if err != nil {
				fmt.Printf("Failed to export gif %s: %s\n", filename, err.Error())
			}
		}
	}
}
