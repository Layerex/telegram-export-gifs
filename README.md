# telegram-export-gifs

Export gifs from telegram

**WARNING: this program uses a hack to get older gifs not visible from the client: it temporarily removes newer gifs to get older ones. If something goes wrong (for example if program gets killed, power goes off or internet disconnects), removed gifs will be lost**

## Installation

```sh
go build
sudo install telegram-export-gifs /usr/local/bin
```

## Usage

```text
usage: telegram-export-gifs [-h|--help] [-d|--directory "<value>"] [--app-id
                            <integer>] [--app-hash "<value>"]

                            Export gifs from telegram

Arguments:

  -h  --help       Print help information
  -d  --directory  Directory to export gifs to. Default: gifs
      --app-id     Test credentials are used by default. Default: 17349
      --app-hash   Test credentials are used by default. Default:
                   344583e45741c457fe1862106095a5eb
```
