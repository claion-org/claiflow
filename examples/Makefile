EXAMPLE_DIR_NAME=${example}

go-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./${EXAMPLE_DIR_NAME}/${EXAMPLE_DIR_NAME} ./${EXAMPLE_DIR_NAME}

docker-build:
	docker build \
		--build-arg EXAMPLE_DIR_NAME=${EXAMPLE_DIR_NAME} \
		-t ${image}:${EXAMPLE_DIR_NAME} -f Dockerfile .

docker-push:
	docker push "${image}:${EXAMPLE_DIR_NAME}"

clean:
	rm ./${EXAMPLE_DIR_NAME}/${EXAMPLE_DIR_NAME}