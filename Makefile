.PHONY: migrate-up migrate-down

migrate-up:
	docker exec -i weather-postgres psql -U weather -d weather -f - < migrations/001_create_weather_snapshot_table.up.sql

migrate-down:
	docker exec -i weather-postgres psql -U weather -d weather -f - < migrations/001_create_weather_snapshot_table.down.sql