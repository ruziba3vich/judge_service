-- db/queries/problems.sql

-- name: CreateProblem :exec
INSERT INTO problems (
    id, title, description, time_limit_ms, memory_limit_mb
) VALUES (
    $1, $2, $3, $4, $5
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
SELECT id, title, description, time_limit_ms, memory_limit_mb
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
    title = $2,
    description = $3,
    time_limit_ms = $4,
    memory_limit_mb = $5
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
