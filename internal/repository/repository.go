package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruziba3vich/judge_service/genprotos/judge_service"
	"github.com/ruziba3vich/judge_service/internal/pkg/helper"
)

// Problem represents the internal data model for a problem in the repository.
// This closely mirrors the database schema.
type Problem struct {
	ID                string
	Title             string
	Description       string
	TimeLimitMs       int32
	MemoryLimitMb     int32
	LanguageHarnesses map[string]string
	TestCases         []*TestCase
}

// TestCase represents a single test case, also mirroring the database schema.
type TestCase struct {
	ID             int32 // This is the database's auto-incrementing ID, not exposed in proto
	ProblemID      string
	Input          string
	ExpectedOutput string
	IsHidden       bool
}

// JudgeRepository defines the interface for data access operations related to problems.
type JudgeRepository interface {
	CreateProblem(ctx context.Context, p Problem) (Problem, error)
	GetProblem(ctx context.Context, problemID uuid.UUID) (Problem, error)
	UpdateProblem(ctx context.Context, p Problem) (Problem, error)
	DeleteProblem(ctx context.Context, problemID uuid.UUID) error
}

func FromProtoCreateProblemRequest(req *judge_service.CreateProblemRequest) Problem {
	repoProblem := Problem{
		ID:                helper.GenerateTimeUUID(),
		Title:             req.GetTitle(),
		Description:       req.GetDescription(),
		LanguageHarnesses: req.GetLanguageHarnesses(),
	}
	if req.GetLimits() != nil {
		repoProblem.TimeLimitMs = req.GetLimits().GetTimeLimitMs()
		repoProblem.MemoryLimitMb = req.GetLimits().GetMemoryLimitMb()
	}

	for _, tc := range req.GetTestCases() {
		repoProblem.TestCases = append(repoProblem.TestCases, &TestCase{
			Input:          tc.GetInput(),
			ExpectedOutput: tc.GetExpectedOutput(),
			IsHidden:       tc.GetIsHidden(),
		})
	}
	return repoProblem
}

// FromProtoUpdateProblemRequest converts a pb.UpdateProblemRequest into a repository.Problem.
func FromProtoUpdateProblemRequest(req *judge_service.UpdateProblemRequest) (Problem, error) {
	p := req.GetProblem()
	repoProblem := Problem{
		ID:                p.Id,
		Title:             p.GetTitle(),
		Description:       p.GetDescription(),
		LanguageHarnesses: p.GetLanguageHarnesses(),
	}
	if p.GetLimits() != nil {
		repoProblem.TimeLimitMs = p.GetLimits().GetTimeLimitMs()
		repoProblem.MemoryLimitMb = p.GetLimits().GetMemoryLimitMb()
	}

	for _, tc := range p.GetTestCases() {
		repoProblem.TestCases = append(repoProblem.TestCases, &TestCase{
			Input:          tc.GetInput(),
			ExpectedOutput: tc.GetExpectedOutput(),
			IsHidden:       tc.GetIsHidden(),
		})
	}
	return repoProblem, nil
}

// ToProtoProblem converts a repository.Problem back into a pb.Problem.
func ToProtoProblem(p Problem) *judge_service.Problem {
	protoProblem := &judge_service.Problem{
		Id:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Limits: &judge_service.ResourceLimits{
			TimeLimitMs:   p.TimeLimitMs,
			MemoryLimitMb: p.MemoryLimitMb,
		},
		LanguageHarnesses: p.LanguageHarnesses,
	}

	for _, tc := range p.TestCases {
		// Note: The internal TestCase.ID (from DB) is not exposed in the proto message.
		protoProblem.TestCases = append(protoProblem.TestCases, &judge_service.TestCase{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
		})
	}
	return protoProblem
}
