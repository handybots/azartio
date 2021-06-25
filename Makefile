include .env
export

run:
	go run ./cmd/bot

compose:
	docker-compose up -d

migrate:
	docker-compose up --force-recreate migrate

deploy-arm:
	docker buildx build --platform linux/amd64 -t docker.pkg.github.com/handybots/azartio/bot:$$(git rev-parse --short HEAD) .
	docker push docker.pkg.github.com/handybots/azartio/bot:$$(git rev-parse --short HEAD);

