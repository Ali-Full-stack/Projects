CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hash_password VARCHAR(255) NOT NULL,
    created_at VARCHAR(255) DEFAULT NOW()
);