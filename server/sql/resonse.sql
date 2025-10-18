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
