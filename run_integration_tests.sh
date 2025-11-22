docker compose -f docker-compose.tests.yaml down -v
docker compose -f docker-compose.tests.yaml up --build -d --wait
go test ./tests/integration
docker compose -f docker-compose.tests.yaml down -v
