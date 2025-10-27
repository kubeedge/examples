.PHONY: test build-webUI build-modbus build-hardware

test: deploy-hardware deploy-modbus deploy-webUI

clean:
	kubectl delete -f ./webUI/resource
	kubectl delete -f ./hardware/resource
	make -C ./modbus clean

deploy-crds:
	kubectl apply -f ./crds/temperature-model.yaml
	kubectl apply -f ./crds/temperature-instance.yaml
deploy-webUI:
	kubectl apply -f ./webUI/resource

deploy-hardware:
	kubectl apply -f ./hardware/resource

deploy-modbus: 
	make -C ./modbus deploy

build-webUI:
	docker build -t temperature-webui:v1.0 ./webUI
#	kind load docker-image temperature-webui:v1.0 --nodes=kind-control-plane
build-modbus:
	make -C ./modbus docker-build

build-hardware:
	docker build -t temperature-sensor:v1.0 ./hardware
#	kind load docker-image temperature-sensor:v1.0 --nodes=kind-worker

test-webUI:
	go build -o main ./main.go
	docker cp ./main kind-control-plane:./srv/temperature/main
	docker cp ./static kind-control-plane:./srv/temperature/static