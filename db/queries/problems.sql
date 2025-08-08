-- db/queries/problems.sql

-- name: CreateProblem :exec
INSERT INTO problems (
    id, title, description, time_limit_ms, memory_limit_mb, difficulty
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CreateProblemLanguageHarness :exec
INSERT INTO problem_language_harnesses (
    problem_id, language, harness_code
) VALUES (
    $1, $2, $3
);

-- name: CreateTestCase :exec
INSERT INTO test_cases (
    problem_id, input, expected_output, is_hidden
) VALUES (
    $1, $2, $3, $4
);

-- name: GetProblem :one
SELECT id, title, description, time_limit_ms, memory_limit_mb, difficulty
FROM problems
WHERE id = $1;

-- name: GetProblemLanguageHarnesses :many
SELECT language, harness_code
FROM problem_language_harnesses
WHERE problem_id = $1;

-- name: GetTestCases :many
SELECT id, input, expected_output, is_hidden
FROM test_cases
WHERE problem_id = $1;

-- name: UpdateProblem :exec
UPDATE problems
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    time_limit_ms = COALESCE($4, time_limit_ms),
    memory_limit_mb = COALESCE($5, memory_limit_mb),
    difficulty = COALESCE($6, difficulty)
WHERE id = $1;

-- name: DeleteProblemLanguageHarnessesByProblemID :exec
DELETE FROM problem_language_harnesses
WHERE problem_id = $1;

-- name: DeleteTestCasesByProblemID :exec
DELETE FROM test_cases
WHERE problem_id = $1;

-- name: DeleteProblem :exec
DELETE FROM problems
WHERE id = $1;

-- name: GetProblemsWithFilter :many
SELECT 
    id, 
    title, 
    description, 
    difficulty
FROM problems
WHERE 
    ($1::text IS NULL OR difficulty = $1)
ORDER BY id
LIMIT $2 OFFSET $3;
