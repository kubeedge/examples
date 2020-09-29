# make traffic
.PHONY: default image

default: image

image: bin
	docker build -t kubeedge/traffic-light:arm32-v1 .

.PHONY: bin
bin:
	go build -o light

.PHONY: run
run: bin
	./light
