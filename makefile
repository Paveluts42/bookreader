run-app:
	cd backend && DOCKER_BUILDKIT=0 docker-compose up --build -d
	cd frontend && npm install && npm run dev
gen-proto:
	cd backend && buf generate
	cd frontend  && npx buf generate ../backend/api
logs:
	cd backend && docker-compose logs -f