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
const sessionFile = programName + "/tg.session"

func (t *Telegram) GetCurrentGIFs() ([]mtproto.TL_document, error) {
	tl := t.Request(mtproto.TL_messages_getSavedGIFs{})
	savedGIFsRes, ok := tl.(mtproto.TL_messages_savedGIFs)
	if !ok {
		return nil, errors.New("TL_messages_getSavedGIFs failed")
	}

	savedGIFs := make([]mtproto.TL_document, len(savedGIFsRes.GIFs))
	for i, GIFDocument := range savedGIFsRes.GIFs {
		savedGIFs[i] = GIFDocument.(mtproto.TL_document)
	}

	return savedGIFs, nil
}

func (t *Telegram) SaveGIF(document mtproto.TL_inputDocument, unsave bool) {
	t.Request(mtproto.TL_messages_saveGIF{ID: document, Unsave: unsave})
}

func (t *Telegram) GetAllGIFs(savedGIFsLimit int, clearGIFs bool) []mtproto.TL_document {
	savedGIFs := make([]mtproto.TL_document, 0)
	unsavedGIFs := make([]mtproto.TL_inputDocument, 0)

	defer func() {
		errorString := recover()
		if errorString != nil {
			fmt.Printf("Error when getting GIFs: %s\n", errorString)
			defer panic(errorString)
		}
		if !clearGIFs {
			for i := range unsavedGIFs {
				inputDocument := unsavedGIFs[len(unsavedGIFs)-i-1]
				fmt.Printf("(%d/%d) Resaving GIFs\n", i+1, len(unsavedGIFs))
				t.SaveGIF(inputDocument, false) // Resave GIF
			}
		}
	}()

	// Handle interrupts
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	fmt.Println("WARNING: this program uses a hack to get older GIFs not visible from client: it temporarily removes newer GIFs to get older ones. If something goes wrong (for example if program gets killed, power goes off or internet disconnects), removed GIFs will be lost.")
	for {
		currentSavedGIFs, err := t.GetCurrentGIFs()
		if err != nil {
			panic(err)
		}
		savedGIFs = append(savedGIFs, currentSavedGIFs...)
		if len(currentSavedGIFs) != savedGIFsLimit { // if there is not max amount of GIFs, there are no more
			fmt.Printf("Unsaved %d GIFs\n", len(savedGIFs)-len(currentSavedGIFs))
			fmt.Printf("Got %d GIFs\n", len(savedGIFs))
			return savedGIFs
		}
		for i, document := range currentSavedGIFs {
			select {
			case <-interruptChan:
				return make([]mtproto.TL_document, 0)
			default:
				inputDocument := mtproto.TL_inputDocument{ID: document.ID, AccessHash: document.AccessHash, FileReference: document.FileReference}
				// fmt.Println(inputDocument)
				fmt.Printf("(%d/%d+) Unsaving GIFs\n", len(savedGIFs)-len(currentSavedGIFs)+i+1, len(savedGIFs))
				unsavedGIFs = append(unsavedGIFs, inputDocument)
				// GIFs can be lost that way, but i am too lazy to implement anything to recover them
				t.SaveGIF(inputDocument, true) // Temporarily unsave GIF, so old GIFs become visible
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

	tl := t.Request(mtproto.TL_users_getUsers{ID: []mtproto.TL{mtproto.TL_inputUserSelf{}}})
	user, ok := tl.(mtproto.VectorObject)[0].(mtproto.TL_user)
	if !ok {
		panic("TL_users_getUsers failed")
	}
	savedGIFsLimit := 200
	if user.Premium {
		savedGIFsLimit = 400
	}

	savedGIFs := t.GetAllGIFs(savedGIFsLimit, false)

	for i, document := range savedGIFs {
		filename := strconv.FormatInt(document.ID, 10) + ".mp4"

		fileInfo, err := os.Stat(filename)
		exists := !errors.Is(err, os.ErrNotExist)
		if exists && fileInfo.Size() == int64(document.Size) {
			fmt.Printf("(%d/%d) GIF %s already exported\n", i+1, len(savedGIFs), filename)
		} else {
			fmt.Printf("(%d/%d) Exporting GIF %s\n", i+1, len(savedGIFs), filename)
			err := t.DownloadDocument(filename, document)
			if err != nil {
				fmt.Printf("Failed to export GIF %s: %s\n", filename, err.Error())
			}
		}
	}
}
