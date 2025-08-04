package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ruziba3vich/judge_service/db/sqlc"
	"github.com/ruziba3vich/judge_service/genprotos/judge_service"
	"github.com/ruziba3vich/judge_service/internal/pkg/helper"
)

// PostgresRepository implements JudgeRepository for PostgreSQL.
type PostgresRepository struct {
	queries *sqlc.Queries // sqlc generated queries object
	db      *sql.DB       // Raw DB connection for transaction management
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

// CreateProblem inserts a new problem and its associated data into the database within a transaction.
func (r *PostgresRepository) CreateProblem(ctx context.Context, p *judge_service.CreateProblemRequest) (*judge_service.Problem, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %s", err.Error())
	}
	qtx := r.queries.WithTx(tx)
	var finalErr error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback
		}
		if finalErr != nil { // If an error occurred in the main logic
			tx.Rollback() // Rollback the transaction
		} else {
			finalErr = tx.Commit() // Otherwise, commit and capture any commit error
		}
	}()

	ID := helper.GenerateTimeUUID()

	// Insert into problems table
	finalErr = qtx.CreateProblem(ctx, sqlc.CreateProblemParams{
		ID:            ID,
		Title:         p.Title,
		Description:   p.Description,
		TimeLimitMs:   p.Limits.TimeLimitMs,
		MemoryLimitMb: p.Limits.MemoryLimitMb,
	})
	if finalErr != nil {
		return nil, fmt.Errorf("failed to create problem: %s", finalErr.Error())
	}

	// Insert language harnesses
	for lang, code := range p.LanguageHarnesses {
		finalErr = qtx.CreateProblemLanguageHarness(ctx, sqlc.CreateProblemLanguageHarnessParams{
			ProblemID:   ID,
			Language:    lang,
			HarnessCode: code,
		})
		if finalErr != nil {
			return nil, fmt.Errorf("failed to create language harness for problem %s: %s", ID, finalErr.Error())
		}
	}

	// Insert test cases
	for _, tc := range p.TestCases {
		finalErr = qtx.CreateTestCase(ctx, sqlc.CreateTestCaseParams{
			ProblemID:      ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
		})
		if finalErr != nil {
			return nil, fmt.Errorf("failed to create test case for problem %s: %s", ID, finalErr.Error())
		}
	}

	return &judge_service.Problem{
		Id:                ID,
		Title:             p.Title,
		Description:       p.Description,
		LanguageHarnesses: p.LanguageHarnesses,
		TestCases:         p.TestCases,
		Limits:            p.Limits,
	}, finalErr
}

// GetProblem retrieves a problem and its associated data from the database.
func (r *PostgresRepository) GetProblem(ctx context.Context, problemID string) (*judge_service.Problem, error) {
	// Get main problem data
	sqlcProblem, err := r.queries.GetProblem(ctx, problemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("problem not found: %s", err.Error())
		}
		return nil, fmt.Errorf("failed to get problem: %s", err.Error())
	}

	repoProblem := judge_service.Problem{
		Id:          problemID,
		Title:       sqlcProblem.Title,
		Description: sqlcProblem.Description,
		Limits: &judge_service.ResourceLimits{
			TimeLimitMs:   sqlcProblem.TimeLimitMs,
			MemoryLimitMb: sqlcProblem.MemoryLimitMb,
		},
		LanguageHarnesses: make(map[string]string),
	}

	// Get language harnesses
	harnesses, err := r.queries.GetProblemLanguageHarnesses(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get language harnesses for problem %s: %w", problemID, err)
	}
	for _, h := range harnesses {
		repoProblem.LanguageHarnesses[h.Language] = h.HarnessCode
	}

	// Get test cases
	testCases, err := r.queries.GetTestCases(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases for problem %s: %w", problemID, err)
	}
	for _, tc := range testCases {
		repoProblem.TestCases = append(repoProblem.TestCases, &judge_service.TestCase{
			Id:             int32(tc.ID),
			ProblemId:      problemID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
		})
	}

	return &repoProblem, nil
}

// UpdateProblem (apply similar defer logic)
func (r *PostgresRepository) UpdateProblem(ctx context.Context, p *judge_service.UpdateProblemRequest) (*judge_service.Problem, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %s", err)
	}
	qtx := r.queries.WithTx(tx)

	var finalErr error // Declare named return error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		if finalErr != nil {
			tx.Rollback()
		} else {
			finalErr = tx.Commit()
		}
	}()

	// Update main problem data
	finalErr = qtx.UpdateProblem(ctx, sqlc.UpdateProblemParams{
		ID:            p.Problem.Id,
		Title:         p.Problem.Title,
		Description:   p.Problem.Description,
		TimeLimitMs:   p.Problem.Limits.TimeLimitMs,
		MemoryLimitMb: p.Problem.Limits.MemoryLimitMb,
	})
	if finalErr != nil {
		return nil, fmt.Errorf("failed to update problem %s: %w", p.Problem.Id, finalErr)
	}

	// Delete existing language harnesses and insert new ones
	finalErr = qtx.DeleteProblemLanguageHarnessesByProblemID(ctx, p.Problem.Id)
	if finalErr != nil {
		return nil, fmt.Errorf("failed to delete existing language harnesses for problem %s: %w", p.Problem.Id, finalErr)
	}
	for lang, code := range p.Problem.LanguageHarnesses {
		finalErr = qtx.CreateProblemLanguageHarness(ctx, sqlc.CreateProblemLanguageHarnessParams{
			ProblemID:   p.Problem.Id,
			Language:    lang,
			HarnessCode: code,
		})
		if finalErr != nil {
			return nil, fmt.Errorf("failed to create new language harness for problem %s: %w", p.Problem.Id, finalErr)
		}
	}

	// Delete existing test cases and insert new ones
	finalErr = qtx.DeleteTestCasesByProblemID(ctx, p.Problem.Id)
	if finalErr != nil {
		return nil, fmt.Errorf("failed to delete existing test cases for problem %s: %w", p.Problem.Id, finalErr)
	}
	for _, tc := range p.Problem.TestCases {
		finalErr = qtx.CreateTestCase(ctx, sqlc.CreateTestCaseParams{
			ProblemID:      p.Problem.Id,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
		})
		if finalErr != nil {
			return nil, fmt.Errorf("failed to create new test case for problem %s: %w", p.Problem.Id, finalErr)
		}
	}

	return p.Problem, finalErr
}

// DeleteProblem deletes a problem and all its associated data.
// Due to ON DELETE CASCADE, related records in `problem_language_harnesses` and `test_cases`
// tables will be automatically deleted by the database.
func (r *PostgresRepository) DeleteProblem(ctx context.Context, problemID string) error {
	err := r.queries.DeleteProblem(ctx, problemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("problem not found for deletion: %w", err)
		}
		return fmt.Errorf("failed to delete problem %s: %w", problemID, err)
	}
	return nil
}
