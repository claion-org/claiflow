PACKAGE=github.com/claion-org/claiflow/pkg
VERSION=$(shell sed -n 's/VERSION=//p' properties)
COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS=-X $(PACKAGE)/version.Version=$(VERSION) -X $(PACKAGE)/version.Commit=$(COMMIT) -X $(PACKAGE)/version.BuildDate=$(BUILD_DATE)

grpc-build:
	go generate

swag-prep:
	go install github.com/swaggo/swag/cmd/swag@v1.8.7

swag-build: swag-prep
	cd pkg/server/route;go generate

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/claiflow ./cmd/claiflow

docker-login:
	docker login ${register} -u ${user}

docker-build:
	docker build -t ${image}:$(VERSION) -f Dockerfile .

docker-push:
	docker push ${image}:$(VERSION)

docker-buildx-and-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t ${image}:${VERSION} -f Dockerfile --push .

clean:
	rm ./bin/claiflow
