CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    login TEXT UNIQUE,
    password TEXT,
    salt TEXT,
    role TEXT
);