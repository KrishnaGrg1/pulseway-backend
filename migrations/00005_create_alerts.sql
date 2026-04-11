-- +goose Up
CREATE TABLE alerts (
  id          BIGSERIAL PRIMARY KEY,
  monitor_id  BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  type        TEXT NOT NULL,
  destination TEXT NOT NULL,
  created_at  TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS alerts;