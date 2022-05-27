package main

import (
	"fmt"
	"os"
	"strconv"
)

const helpMessage =
`usage: %s [-h] [-d DIRECTORY] [--app-id APP_ID] [--app-hash APP_HASH]

Export saved gifs from telegram.

options:
  -h, --help            Show this help message and exit
  -d DIRECTORY, --directory DIRECTORY
                        Directory to export gifs to
  --app-id APP_ID       Test credentials are used by default
  --app-hash APP_HASH   Test credentials are used by default

WARNING: this program uses a hack to get older gifs not visible from client: it temporarily removes newer gifs to get older ones. If something goes wrong (for example if program gets killed, power goes off or internet disconnects), removed gifs will be lost.
`

type Args struct {
	Directory       string
	AppID           int32
	AppHash         string
}

func ParseArgs() Args {

	var args Args
	end := len(os.Args) - 1
	for i := 1; i < len(os.Args); i++ {
		inc := func() {
			if i == end {
				panic(fmt.Sprintf("Option %s requires a value", os.Args[i]))
			}
			i++
		}
		switch os.Args[i] {
		case "-d", "--directory":
			inc()
			args.Directory = os.Args[i]
		case "--app-id":
			inc()
			argument, err := strconv.Atoi(os.Args[i])
			if err != nil {
				panic("--app-id value has to be a 32-bit integer")
			}
			args.AppID = int32(argument)
		case "--app-hash":
			inc()
			if len(os.Args[i]) != 32 || !IsHex(os.Args[i]) {
				panic("--app-hash value has to be a hex string of 32 characters")
			}
			args.AppHash = os.Args[i]
		case "-h", "--help":
			fmt.Printf(helpMessage, os.Args[0])
			os.Exit(0)
		default:
			panic(fmt.Sprintf("Unexpected argument: %s", os.Args[i]))
		}
	}

	if args.Directory == "" {
		args.Directory = "gifs/"
	}
	if args.AppID == 0 {
		if args.AppHash != "" {
			panic("--app-hash is provided, but --app-id isn't")
		}
		args.AppID = 17349;
	}
	if args.AppHash == "" {
		if args.AppID == 0 {
			panic("--app-id is provided, but --app-hash isn't")
		}
		args.AppHash = "344583e45741c457fe1862106095a5eb";
	}
	return args
}
