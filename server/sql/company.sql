-- name: GetCompanyByID :one
SELECT
    id,
    title,
    description,
    contacts,
    inn,
    address,
    approved,
    representative_id,
    login,
    password_hash,
    created_at,
    updated_at
FROM companies
WHERE id = @id;

-- name: GetCompanyByLogin :one
SELECT
    id,
    title,
    description,
    contacts,
    inn,
    address,
    approved,
    representative_id,
    login,
    password_hash,
    created_at,
    updated_at
FROM companies
WHERE login = @login;

-- name: GetCompanyByINN :one
SELECT
    id,
    title,
    description,
    contacts,
    inn,
    address,
    approved,
    representative_id,
    login,
    password_hash,
    created_at,
    updated_at
FROM companies
WHERE inn = @inn;

-- name: CreateCompany :exec
INSERT INTO companies (
    id, title, description, contacts, inn, address, approved, representative_id, login, password_hash, created_at, updated_at
) VALUES (
    @id, @title, @description, @contacts, @inn, @address, @approved, @representative_id, @login, @password_hash, @created_at, @updated_at
);

-- name: UpdateCompany :exec
UPDATE companies
SET
    title = @title,
    description = @description,
    contacts = @contacts,
    inn = @inn,
    address = @address,
    approved = @approved,
    representative_id = @representative_id,
    login = @login,
    password_hash = @password_hash,
    created_at = @created_at,
    updated_at = @updated_at
WHERE id = @id;
