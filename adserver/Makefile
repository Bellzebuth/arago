.PHONY: build run proto up down logs

proto:
	protoc --go_out=proto/ad --go-grpc_out=proto/ad \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		proto/ad/ad.proto

build:
	docker build -t adserver .

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

run: proto up
