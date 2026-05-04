-- +goose up
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);

CREATE INDEX idx_user_period ON subscriptions(user_id, start_date, end_date);
CREATE INDEX idx_service_name ON subscriptions(LOWER(service_name));