-- db/migrations/001_initial_schema.up.sql

CREATE TABLE IF NOT EXISTS problems (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    time_limit_ms INTEGER NOT NULL,
    memory_limit_mb INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS problem_language_harnesses (
    problem_id TEXT NOT NULL,
    language TEXT NOT NULL,
    harness_code TEXT NOT NULL,
    PRIMARY KEY (problem_id, language),
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS test_cases (
    id SERIAL PRIMARY KEY,
    problem_id TEXT NOT NULL,
    input TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_hidden BOOLEAN NOT NULL,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_test_cases_problem_id ON test_cases(problem_id);

CREATE INDEX IF NOT EXISTS idx_harnesses_problem_id ON problem_language_harnesses(problem_id);

CREATE INDEX IF NOT EXISTS idx_harnesses_language ON problem_language_harnesses(language);
