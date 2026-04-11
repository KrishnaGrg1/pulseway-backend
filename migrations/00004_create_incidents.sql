-- +goose Up
CREATE TABLE incidents (
  id          BIGSERIAL PRIMARY KEY,
  monitor_id  BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  started_at  TIMESTAMPTZ DEFAULT now(),
  resolved_at TIMESTAMPTZ,
  notified    BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE IF EXISTS incidents;