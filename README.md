# telegram-export-gifs

Export gifs from telegram

**WARNING: this program uses a hack to get older gifs not visible from client: it temporarily removes newer gifs to get older ones. If something goes wrong (for example if program gets killed, power goes off or internet disconnects), removed gifs will be lost.**

## Running

```sh
go build
./telegram-export-gifs
```

## Usage

```text
usage: ./telegram-export-gifs [-h] [-d DIRECTORY] [--dont-save-session] [--app-id APP_ID] [--app-hash APP_HASH]

Export saved gifs from telegram.

options:
  -h, --help            Show this help message and exit
  -d DIRECTORY, --directory DIRECTORY
                        Directory to export gifs to
  --dont-save-session   Don't save session file (and don't use already saved one)
  --app-id APP_ID       Test credentials are used by default
  --app-hash APP_HASH   Test credentials are used by default

Session file is saved to /home/layerex/.local/share/telegram-export-gifs/tg.session

WARNING: this program uses a hack to get older gifs not visible from client: it temporarily removes newer gifs to get older ones. If something goes wrong (for example if program gets killed, power goes off or internet disconnects), removed gifs will be lost.
```
