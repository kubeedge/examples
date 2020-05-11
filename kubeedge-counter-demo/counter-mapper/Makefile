.PHONY: all pi-execute-app docker clean

all: pi-execute-app

pi-execute-app:
	GOARCH=arm64 go build -o pi-counter-app main.go

docker:
	docker build . -t kubeedge/kubeedge-pi-counter:v1.0.0

clean:
	rm -f pi-counter-app
