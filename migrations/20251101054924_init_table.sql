-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

DROP  TABLE IF EXISTS drivers;
CREATE TABLE drivers (
    id VARCHAR(64) PRIMARY KEY, 
    name VARCHAR(64), 
    profile_picture VARCHAR(64), 
    car_plate VARCHAR(24), 
    package_slug VARCHAR(64),
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    user_id VARCHAR(64),
    location JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

DROP  TABLE IF EXISTS riders;
CREATE TABLE riders (
    id VARCHAR(64) PRIMARY KEY, 
    package_slug VARCHAR(64) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    user_id VARCHAR(64),
    location JSONB,
    destination JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

DROP  TABLE IF EXISTS fares;
CREATE TABLE fares (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    package_slug VARCHAR(50) NOT NULL,
    total_price_in_cents DOUBLE PRECISION NOT NULL DEFAULT 0,
    routes JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

DROP  TABLE IF EXISTS transactions;
CREATE TABLE transactions (
    id VARCHAR(64) PRIMARY KEY,
    rider_id VARCHAR(64),
    driver_id VARCHAR(64),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

DROP  TABLE IF EXISTS trips;
CREATE TABLE trips (
    id VARCHAR(64) PRIMARY KEY,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    rider_id VARCHAR(64),
    transaction_id VARCHAR(64),
    driver_id VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);


DROP  TABLE IF EXISTS users;
CREATE TABLE users (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    client_id VARCHAR(100)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE drivers;
DROP TABLE fares;
DROP TABLE transactions;
DROP TABLE trips;
DROP TABLE users;
-- +goose StatementEnd

