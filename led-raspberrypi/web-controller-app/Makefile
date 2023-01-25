.PHONY: all docker clean

all:
	go build -o kubeedge-light-controller main.go

docker:
	docker build . -t kubeedge/kubeedge-light-app:v1.0.0

clean:
	rm -rf ./kubeedge-light-controller
