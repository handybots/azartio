include .env
export

run:
	go run ./cmd/bot 
	
compose:
	docker-compose up -d

migrate:
	docker-compose up --force-recreate migrate
