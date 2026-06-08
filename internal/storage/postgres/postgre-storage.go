package postgres

import (
	"context"
	"time"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *storage {
	return &storage{db: db}
}

func (s *storage) Save(ctx context.Context, snapshot weather.WeatherSnapshot) error {
	sqlE := ` 
	INSERT INTO weather (id, city, provider, temperature_celsius, wind_speed, observed_at, raw_payload, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)		 
	`
	id := uuid.New()

	_, err := s.db.Exec(ctx,
		sqlE,
		id,
		weather.CityNormalize(snapshot.City),
		snapshot.Provider,
		snapshot.TemperatureCelsius,
		snapshot.WindSpeed,
		snapshot.ObservedAt,
		snapshot.RawPayload,
		time.Now(),
	)
	return err
}

func (s *storage) History(ctx context.Context, city string, limit int) ([]weather.WeatherSnapshot, error) {
	city = weather.CityNormalize(city)
	if limit > 50 {
		limit = 50
	}
	if limit <= 0 {
		limit = 10
	}

	sqlQ := `
	SELECT city, provider, temperature_celsius, wind_speed, observed_at, raw_payload FROM weather
	WHERE city = $1
	ORDER BY created_at DESC
	LIMIT $2
	`

	rows, err := s.db.Query(ctx, sqlQ, city, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snapshots []weather.WeatherSnapshot

	for rows.Next() {
		var snap weather.WeatherSnapshot
		err := rows.Scan(
			&snap.City,
			&snap.Provider,
			&snap.TemperatureCelsius,
			&snap.WindSpeed,
			&snap.ObservedAt,
			&snap.RawPayload,
		)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snap)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}
