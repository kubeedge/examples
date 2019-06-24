.PHONY: all build kubeedge-wechat-app docker clean

all:build

build:kubeedge-wechat-app

kubeedge-wechat-app:
	go build -o kubeedge-wechat-app main.go

docker:kubeedge-wechat-app
	docker build . -t xiangxinyong/kubeedge-wechat-app:v2.0

clean:
	rm -rf ./kubeedge-wechat-app
