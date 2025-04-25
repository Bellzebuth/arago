.PHONY: proto build up down

protoc:
	$(MAKE) -C adserver protoc
	$(MAKE) -C tracker protoc

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down
