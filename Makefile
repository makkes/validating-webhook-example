DOCKER_IMG ?= "makkes/webhook-example:dev"

docker-build:
	docker build -t $(DOCKER_IMG) .

docker-push: docker-build
	docker push $(DOCKER_IMG)

deploy-cert-manager:
	kubectl create namespace cert-manager
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.12.0-beta.1/cert-manager.yaml

deploy-ca:
	kubectl apply -f ca.yaml

deploy-webhook: deploy-ca
 	kubectl apply -f webhook.yaml

deploy: docker-push deploy-cert-manager deploy-ca deploy-webhook