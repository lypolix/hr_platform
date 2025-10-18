-- Universities

-- name: GetUniversityByID :one
SELECT
    id,
    login,
    password_hash,
    inn,
    title,
    confirmed,
    created_at,
    updated_at
FROM universities
WHERE id = @id;

-- name: GetUniversityByLogin :one
SELECT
    id,
    login,
    password_hash,
    inn,
    title,
    confirmed,
    created_at,
    updated_at
FROM universities
WHERE login = @login;

-- name: CreateUniversity :exec
INSERT INTO universities (id, title, login, password_hash, inn, confirmed, created_at, updated_at)
VALUES (@id, @title, @login, @password_hash, @inn, @confirmed, @created_at, @updated_at);

-- name: UpdateUniversity :exec
UPDATE universities
SET
    login = @login,
    password_hash = @password_hash,
    title = @title,
    inn = @inn,
    confirmed = @confirmed,
    created_at = @created_at,
    updated_at = @updated_at
WHERE id = @id;
