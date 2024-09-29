.PHONEY: all, pack, clear, format, run
PATH = ./cmd/server/server.go
NAME = ./song_server

all: clear format	pack run

format:
	go fmt ./...

pack:
	go build -o ${NAME} ${PATH}

run:
	${NAME}

clear:
	rm -f ${NAME}