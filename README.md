# Wether-Snapshot-Service

## Назначение

- Получение текущей погоды через внешний API (Open‑Meteo).
- Кэширование результатов на 5 минут с помощью in‑memory TTL cache.
- Сохранение каждого успешного ответа в PostgreSQL.
- Предоставление истории запросов по городу.
- Метрики Prometheus: количество запросов и количество cache hit/miss.
- Health эндпоинт.

## Требования

- Go 1.22+
- Docker/Docker Compose (для PostgreSQL)

## База данных

- Для миграции базы данных (для PostgreSQL) используется make-команда:
```make migrate-up```
- Для отката миграции (для PostgreSQL) используется make-команда:
```make migrate-down```

## Запуск сервиса

- Для запуска сервиса исользуется make-команда:
```make run-service```
- или вручную:
```go run cmd/weather-service/main.go```

## Переменные окружения

Создайте файл `.env` в корне проекта (или экспортируйте переменные вручную):
- Порт вашего сервера `PORT=`
- Ссылку на базу данных (PostgreSQL): `DATABASE_URL=`
- Ссылки на внешний API: `GEOCODING_ENDPOINT=` и `FORECAST_ENDPOINT=`
Продолжительно жизни кэша (в минутах) `TTL=`

## Пример запрсов
- Запрос:   
`curl http://localhost:3000/api/v1/weather/current?city=MURManSK`  
Ответ:  
`{
  "city": "murmansk",  
  "provider": "open-meteo",
  "temperature_celsius": 16.9,
  "wind_speed": 6.5,
  "observed_at": "2026-06-10T21:00:00Z",
  "cached": false
}`
- Запрос  
`curl http://localhost:3000/api/v1/weather/history?city=murmansk&limit=1`  
Ответ:  
`{
  "items": [
    {
      "city": "murmansk",
      "provider": "open-meteo",
      "temperature_celsius": 16.9,
      "wind_speed": 6.5,
      "observed_at": "2026-06-11T00:00:00+03:00",
      "cached": false
    }
  ]
}`

## Запуск тестов
- Запуск всех тестов:
`go test -v -race ./...`

## Возможные улучшения
- Graceful shutdown не реализован (сервер завершается мгновенно по Ctrl+C).
- Redis как опциональный кэш не реализован (только in‑memory).
- Интеграция миграций в код.