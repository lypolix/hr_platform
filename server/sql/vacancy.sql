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
