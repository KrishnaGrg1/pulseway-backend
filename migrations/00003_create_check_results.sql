-- +goose Up
CREATE TABLE check_results (
  id          BIGSERIAL PRIMARY KEY,
  monitor_id  BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  status      TEXT NOT NULL,
  latency_ms  INT NOT NULL,
  status_code INT,
  checked_at  TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_check_results_monitor_id ON check_results(monitor_id);
CREATE INDEX idx_check_results_checked_at ON check_results(checked_at);

-- +goose Down
DROP TABLE IF EXISTS check_results;