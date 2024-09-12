CREATE TABLE IF  NOT EXISTS users(
    id   VARCHAR(255) PRIMARY KEY,
    fullname VARCHAR(100)  NOT NULL,
    email VARCHAR(65) UNIQUE,
    hash_password VARCHAR(255) NOT NULL,
    rate   FLOAT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS ranks (
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    commentor_id  VARCHAR(255) UNIQUE,
    commentor_name VARCHAR(100) NOT NULL,
    rate FLOAT NOT NULL
);



