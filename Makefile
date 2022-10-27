image:
	docker pull postgres

postgres:
	docker run --name postgres1 -p 5432:5432 -e POSTGRES_USER=evo -e POSTGRES_PASSWORD=evo -d postgres:latest

server:
	go run .

.PHONY: image postgres server