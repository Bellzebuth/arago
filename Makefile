.PHONY: proto build up down

proto:
	$(MAKE) -C adserver proto
	$(MAKE) -C tracker proto

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down
