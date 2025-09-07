up:
    docker compose up -d --build

down:
    docker compose down

logs:
    docker compose logs -f

gen-proto:
    cd backend && buf generate
    cd frontend && npx buf generate ../backend/api

db-clean:
    docker compose down -v