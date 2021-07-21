include .env
export

run:
	go run ./cmd/bot

compose:
	docker-compose up -d

migrate:
	docker-compose up --force-recreate migrate

deploy-from-arm:
	docker buildx build \
		--platform linux/amd64 \
		-t docker.pkg.github.com/handybots/azartio/bot:$$(git rev-parse --short HEAD) \
		-t docker.pkg.github.com/handybots/azartio/bot:latest .
	docker push docker.pkg.github.com/handybots/azartio/bot:$$(git rev-parse --short HEAD)
	docker push docker.pkg.github.com/handybots/azartio/bot:latest

ideas-squick:
	export SQUICK_DRIVER=postgres SQUICK_URL="host=localhost port=5432 user=postgres password=postgres dbname=azartio_ideas sslmode=disable";\
	cd ./cmd/ideas ;\
	squick make -table ideas insert get:id select:used,deleted set:used,deleted ;\
	squick make -table votes insert get:done set:updated_at,message_id,days_left,done ;\
	squick make -table voters insert delete select:vote_id