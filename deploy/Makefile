CONFIG_APPLICATION_NAME=claiflow-config-application
CONFIG_ENIGMA_NAME=claiflow-config-enigma
CONFIG_DATABASE_NAME=claiflow-config-database

.PHONY: namespace
namespace:
	kubectl apply -f namespace.yaml

.PHONY: namespace-clean
namespace-clean:
	kubectl delete -f namespace.yaml

.PHONY: mariadb
mariadb: namespace
	helm install mariadb-claiflow \
		--namespace claiflow \
		--set auth.rootPassword=root \
		--set auth.database=claiflow \
		--set auth.username=claiflow \
		--set auth.password=claiflow \
		--version 11.4.5 \
		oci://registry-1.docker.io/bitnamicharts/mariadb

.PHONY: mariadb-clean
mariadb-clean:
	helm uninstall mariadb-claiflow \
		--namespace claiflow

.PHONY: config
config: namespace
	-kubectl -n claiflow create configmap ${CONFIG_APPLICATION_NAME} \
		--from-file config/application.yaml

	-kubectl -n claiflow create configmap ${CONFIG_ENIGMA_NAME} \
		--from-file config/enigma.yaml

	-kubectl -n claiflow create secret generic ${CONFIG_DATABASE_NAME} \
		--from-file config/database.yaml

.PHONY: config-clean
config-clean:
	-kubectl -n claiflow delete configmap ${CONFIG_APPLICATION_NAME}
	
	-kubectl -n claiflow delete configmap ${CONFIG_ENIGMA_NAME}

	-kubectl -n claiflow delete secret ${CONFIG_DATABASE_NAME}

.PHONY: deployment
deployment: namespace
	kubectl apply -f deployment.yaml

.PHONY: deployment-clean
deployment-clean:
	kubectl delete -f deployment.yaml

.PHONY: service
service: namespace deployment
	kubectl apply -f service.yaml

.PHONY: service-clean
service-clean:
	kubectl delete -f service.yaml

.PHONY: all
all: namespace mariadb config deployment service

.PHONY: clean
clean: service-clean deployment-clean config-clean mariadb-clean namespace-clean
