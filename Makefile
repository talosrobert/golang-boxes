main_package_path=./cmd/web
binary_name=boxes

## tidy: tidy modfiles and format .go files
PHONY: tidy
tidy:
	go mod tidy -v
	go fmt './...'

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${binary_name} ${main_package_path}

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${binary_name}
