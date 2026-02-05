package service

import (
	"context"
	"time"

	"github.com/jackyansen22/crud-category/internal/model"
	"github.com/jackyansen22/crud-category/internal/repository"
)

type ReportService interface {
	GetToday(ctx context.Context) (*model.ReportResponse, error)
	GetByRange(ctx context.Context, start, end time.Time) (*model.ReportResponse, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GetToday(ctx context.Context) (*model.ReportResponse, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	return s.repo.GetReport(ctx, start, end)
}

func (s *reportService) GetByRange(
	ctx context.Context,
	start, end time.Time,
) (*model.ReportResponse, error) {
	return s.repo.GetReport(ctx, start, end)
}
