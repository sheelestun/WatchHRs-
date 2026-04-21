-- +goose Up

CREATE TABLE IF NOT EXISTS managers (
    id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    manager_id UUID NOT NULL REFERENCES managers(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS photos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS screenshot_statistics (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    cnt_mouse_clicks INT NOT NULL CHECK (cnt_mouse_clicks >= 0),
    cnt_keyboard_clicks INT NOT NULL CHECK (cnt_keyboard_clicks >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_sessions (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    total_time INTERVAL
);

CREATE UNIQUE INDEX uniq_active_session
    ON work_sessions(employee_id)
    WHERE end_time IS NULL;

-- +goose Down

DROP TABLE IF EXISTS work_sessions;
DROP TABLE IF EXISTS screenshot_statistics;
DROP TABLE IF EXISTS photos;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS managers;