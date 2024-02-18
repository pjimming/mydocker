DST=./target/bin
GO=$(shell which go)

build:
	@mkdir -p ${DST}
	@GOOS=linux GOARCH=amd64 $(GO) build -o ${DST} .


upload:
	@make build
	@docker cp target/bin/mydocker mydocker:/home