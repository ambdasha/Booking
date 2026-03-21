COMPOSE = docker compose -f deployments/docker/docker-compose.yml
DB_URL  = postgres://postgres:postgres@localhost:5432/booking?sslmode=disable

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f api

migrate-up:
	docker run --rm -v "$(PWD)/migrations:/migrations" migrate/migrate \
		-path=/migrations -database "$(DB_URL)" up

migrate-down:
	docker run --rm -v "$(PWD)/migrations:/migrations" migrate/migrate \
		-path=/migrations -database "$(DB_URL)" down 1