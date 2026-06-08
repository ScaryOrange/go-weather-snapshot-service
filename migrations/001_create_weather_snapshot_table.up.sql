create table weather_snapshot (
  id uuid primary key,
  city text not null,
  provider text not null,
  temperature_celsius numeric(5,2),
  wind_speed numeric(6,2),
  raw_payload jsonb,
  observed_at timestamptz not null,
  created_at timestamptz not null
);

create index idx_weather_snapshot_city_created_at
  on weather_snapshot (city, created_at desc);