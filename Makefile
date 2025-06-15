up:
	docker-compose up
build:
	docker-compose up --build
migration:
	migrate create -ext sql -dir db/migrations $(name)

