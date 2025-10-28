
# make led_light_mapper
.PHONY: default led_light_mapper
led_light_mapper:
	go build light_mapper.go
	docker build -t led-light-mapper:v1.1 .
