build:
	go build -v ./cmd/s3/

debug:
	make build
	./s3.exe -debug

protoc:
	protoc ./proto/s3.proto --go_out=. --go-grpc_out=.

deploy:
	docker-compose down
	docker-compose build
	docker-compose up -d

.DEFAULT_GOAL := debug