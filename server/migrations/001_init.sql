-- Up

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

CREATE TABLE companies (
    id UUID PRIMARY KEY,
    title VARCHAR(512) NOT NULL,
    description TEXT NOT NULL,
    contacts TEXT NOT NULL,
    inn VARCHAR(32) NOT NULL UNIQUE,
    address TEXT NOT NULL,
    approved BOOLEAN NOT NULL,
    representative_id UUID NOT NULL,
    login VARCHAR(256) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX companies_inn_idx ON companies(inn);
CREATE INDEX companies_login_idx ON companies(login);
CREATE INDEX companies_representative_idx ON companies(representative_id);
CREATE INDEX companies_approved_idx ON companies(approved);

CREATE TABLE vacancies (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title VARCHAR(512) NOT NULL,
    description TEXT NOT NULL,
    contacts TEXT NOT NULL,
    requirements TEXT NOT NULL,
    responsibilities TEXT NOT NULL,
    conditions TEXT NOT NULL,
    employment TEXT NOT NULL,
    schedule TEXT NOT NULL,
    experience TEXT NOT NULL,
    education TEXT NOT NULL,
    location TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX vacancies_company_idx ON vacancies(company_id);
CREATE INDEX vacancies_active_idx ON vacancies(is_active);
CREATE INDEX vacancies_location_idx ON vacancies(location);

CREATE TABLE responses (
    id UUID PRIMARY KEY,
    vacancy_id UUID NOT NULL REFERENCES vacancies(id) ON DELETE CASCADE,
    full_name VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL,
    phone VARCHAR(64) NOT NULL,
    cover_letter TEXT NOT NULL,
    resume_url TEXT NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX responses_vacancy_idx ON responses(vacancy_id);
CREATE INDEX responses_status_idx ON responses(status);

---- create above / drop below ----

-- Down

DROP TABLE IF EXISTS responses;
DROP TABLE IF EXISTS vacancies;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS universities;
