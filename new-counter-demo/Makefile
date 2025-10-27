.PHONY: clean build-webUI build-mapper deploy-webUI deploy-mapper undeploy-webUI undeploy-mapper
clean: undeploy-webUI undeploy-mapper
build-webUI:
	docker build -f ./webUI/Dockerfile -t counter-webui:v1.0 ./webUI
#	kind load docker-image counter-webui:v1.0 --nodes=kind-control-plane
build-mapper:
	docker build -f ./virtualprotocol/Dockerfile_nostream -t counter-mapper:v1.0 ./virtualprotocol
#	kind load docker-image counter-mapper:v1.0 --nodes=kind-worker

deploy-webUI:
	kubectl apply -f ./webUI/resource/deploy.yaml

deploy-mapper:
	kubectl apply -f ./virtualprotocol/resource/configmap.yaml
	kubectl apply -f ./virtualprotocol/resource/deployment.yaml
	kubectl wait --for=condition=available --timeout=120s deployment/counter-mapper
	kubectl apply -f ./crds/counter-model.yaml
	kubectl apply -f ./crds/counter-instance.yaml

undeploy-webUI:
	kubectl delete -f ./webUI/resource
undeploy-mapper:
	kubectl delete -f ./virtualprotocol/resource/configmap.yaml
	kubectl delete -f ./virtualprotocol/resource/deployment.yaml
	kubectl delete -f ./crds/counter-model.yaml
	kubectl delete -f ./crds/counter-instance.yaml