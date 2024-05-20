build:
	@go build -o "halo-suster-be" cmd/main.go

migrate-up:
	migrate -database "postgres://postgres:password@localhost:5432/halo-suster?sslmode=disable" -path internal/db/migrations up

migrate-down:
	migrate -database "postgres://postgres:password@localhost:5432/halo-suster?sslmode=disable" -path internal/db/migrations down

run:
	./halo-suster-be

docker-build:
	docker build --no-cache -t halo-suster-be .

docker-run:
	docker compose -f "docker-compose.yaml" up -d --build