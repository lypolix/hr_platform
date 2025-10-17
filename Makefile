up-dev:
	docker compose -f docker-compose.dev.yml up -d

down-dev:
	docker compose -f docker-compose.dev.yml down

include ./server/.env

migrate:
	tern migrate \
	    --host=0.0.0.0 \
		--user=${POSTGRES_USER} \
		--password=${POSTGRES_PASSWORD} \
		--database=${POSTGRES_DB} \
		-m ./server/migrations
