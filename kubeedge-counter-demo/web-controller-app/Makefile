.PHONY: all docker clean

all:
	go build -o kubeedge-counter-controller main.go

docker:
	docker build . -t kubeedge/kubeedge-counter-app:v1.0.0

clean:
	rm -rf ./kubeedge-counter-controller
