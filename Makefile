.PHONY: all up migrate generate run proto

all: up migrate

run:
	docker-compose up database rabbitmq redis api

up:
	docker-compose up -d

migrate:
	docker-compose exec database sh -c 'psql -U casino < /db/migrations/00001.create_base.sql'

generator:
	docker-compose run --rm generator

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/event.proto
