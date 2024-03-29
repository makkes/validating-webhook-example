DOCKER_IMG ?= "makkes/minimal-webhook:dev"

.DEFAULT_GOAL := docker-build

deploy-cert-manager:
	kubectl create namespace cert-manager
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.10.1/cert-manager.yaml

docker-build:
	docker build -t $(DOCKER_IMG) .

docker-push: docker-build
	docker push $(DOCKER_IMG)

deploy-webhook:
	cat webhook.yaml | sed s~##DOCKER_IMG##~$(DOCKER_IMG)~ | kubectl apply -f -

deploy: docker-push deploy-webhook

undeploy:
	kubectl delete -f webhook.yaml