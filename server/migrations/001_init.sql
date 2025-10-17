CREATE TABLE universities (
    id UUID PRIMARY KEY,
    title VARCHAR(512) NOT NULL,
    login VARCHAR(256) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    inn VARCHAR(32) NOT NULL UNIQUE,
    confirmed BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS universities;
