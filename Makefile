.PHONY: build run

build:
	go build .

run:
	go build . && ./ebiten-demo

clean:
	rm main
