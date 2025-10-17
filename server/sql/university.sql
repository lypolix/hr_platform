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

-- name: GetUniversityByINN :one
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
WHERE inn = @inn;

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



-- Companies

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



-- Vacancies

-- name: GetVacancyByID :one
SELECT
    id,
    company_id,
    title,
    description,
    contacts,
    requirements,
    responsibilities,
    conditions,
    employment,
    schedule,
    experience,
    education,
    location,
    is_active,
    created_at,
    updated_at
FROM vacancies
WHERE id = @id;

-- name: ListVacanciesByCompany :many
SELECT
    id,
    company_id,
    title,
    description,
    contacts,
    requirements,
    responsibilities,
    conditions,
    employment,
    schedule,
    experience,
    education,
    location,
    is_active,
    created_at,
    updated_at
FROM vacancies
WHERE company_id = @company_id
ORDER BY created_at DESC;

-- name: ListActiveVacancies :many
SELECT
    id,
    company_id,
    title,
    description,
    contacts,
    requirements,
    responsibilities,
    conditions,
    employment,
    schedule,
    experience,
    education,
    location,
    is_active,
    created_at,
    updated_at
FROM vacancies
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: CreateVacancy :exec
INSERT INTO vacancies (
    id,
    company_id,
    title,
    description,
    contacts,
    requirements,
    responsibilities,
    conditions,
    employment,
    schedule,
    experience,
    education,
    location,
    is_active,
    created_at,
    updated_at
) VALUES (
    @id,
    @company_id,
    @title,
    @description,
    @contacts,
    @requirements,
    @responsibilities,
    @conditions,
    @employment,
    @schedule,
    @experience,
    @education,
    @location,
    @is_active,
    @created_at,
    @updated_at
);

-- name: UpdateVacancy :exec
UPDATE vacancies
SET
    company_id = @company_id,
    title = @title,
    description = @description,
    contacts = @contacts,
    requirements = @requirements,
    responsibilities = @responsibilities,
    conditions = @conditions,
    employment = @employment,
    schedule = @schedule,
    experience = @experience,
    education = @education,
    location = @location,
    is_active = @is_active,
    created_at = @created_at,
    updated_at = @updated_at
WHERE id = @id;



-- Responses

-- name: GetResponseByID :one
SELECT
    id,
    vacancy_id,
    full_name,
    email,
    phone,
    cover_letter,
    resume_url,
    status,
    created_at,
    updated_at
FROM responses
WHERE id = @id;

-- name: ListResponsesByVacancy :many
SELECT
    id,
    vacancy_id,
    full_name,
    email,
    phone,
    cover_letter,
    resume_url,
    status,
    created_at,
    updated_at
FROM responses
WHERE vacancy_id = @vacancy_id
ORDER BY created_at DESC;

-- name: CreateResponse :exec
INSERT INTO responses (
    id,
    vacancy_id,
    full_name,
    email,
    phone,
    cover_letter,
    resume_url,
    status,
    created_at,
    updated_at
) VALUES (
    @id,
    @vacancy_id,
    @full_name,
    @email,
    @phone,
    @cover_letter,
    @resume_url,
    @status,
    @created_at,
    @updated_at
);

-- name: UpdateResponseStatus :exec
UPDATE responses
SET
    status = @status,
    updated_at = @updated_at
WHERE id = @id;
