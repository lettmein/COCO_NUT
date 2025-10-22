.PHONY: start stop test clean

start:
	@echo "Stopping and cleaning all containers, volumes and images..."
	docker compose down -v 2>/dev/null || true
	docker system prune -af --volumes 2>/dev/null || true
	@echo "Starting services..."
	docker compose up --build -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services are running!"
	@echo "Request Service: http://localhost:8081/health"
	@echo "Office Service: http://localhost:8082/health"
	@docker compose ps

stop:
	@echo "Stopping all services..."
	docker compose down
	@echo "Services stopped!"

test:
	@echo "Running tests..."
	@bash ./test.sh

clean:
	@echo "Cleaning everything..."
	docker compose down -v
	docker system prune -af --volumes
	@echo "Cleanup complete!"

logs:
	docker compose logs -f

status:
	@docker compose ps
	@echo ""
	@echo "Health checks:"
	@curl -s http://localhost:8081/health | jq . 2>/dev/null || echo "Request Service not available"
	@curl -s http://localhost:8082/health | jq . 2>/dev/null || echo "Office Service not available"
