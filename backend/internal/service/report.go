package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ReportService interface {
	GenerateFieldReport(fieldID int, from, to time.Time) (models.FieldReport, error)
}

type reportService struct {
	fieldService repository.FieldRepository
	metricRepo   repository.MetricRepository
	forecastRepo repository.ForecastRepository
	log          *zap.Logger
}

func NewReportService(
	fieldRepo repository.FieldRepository,
	metricRepo repository.MetricRepository,
	forecastRepo repository.ForecastRepository,
	log *zap.Logger,
) ReportService {
	return &reportService{
		fieldService: fieldRepo,
		metricRepo:   metricRepo,
		forecastRepo: forecastRepo,
		log:          log,
	}
}

func (s *reportService) GenerateFieldReport(fieldID int, from, to time.Time) (models.FieldReport, error) {
	if from.After(to) {
		return models.FieldReport{}, fmt.Errorf("'from' date cannot be after 'to' date")
	}

	field, err := s.fieldService.GetByIDWithCrop(fieldID)
	if err != nil {
		return models.FieldReport{}, fmt.Errorf("failed to get field info: %w", err)
	}

	aggregates, err := s.metricRepo.GetAggregatedMetricsByField(fieldID, from, to)
	if err != nil {
		return models.FieldReport{}, fmt.Errorf("failed to get metric aggregates: %w", err)
	}

	var temperatureSummary models.MetricSummary
	var humiditySummary models.MetricSummary
	for _, row := range aggregates {
		switch row.SensorType {
		case "temperature":
			temperatureSummary = models.MetricSummary{Avg: row.Avg, Min: row.Min, Max: row.Max}
		case "air_humidity":
			humiditySummary = models.MetricSummary{Avg: row.Avg, Min: row.Min, Max: row.Max}
		}
	}

	forecastStats, err := s.forecastRepo.GetForecastStatisticsByField(fieldID, from, to)
	if err != nil {
		return models.FieldReport{}, fmt.Errorf("failed to get forecast statistics: %w", err)
	}

	return models.FieldReport{
		FieldName:                  field.Name,
		Temperature:                temperatureSummary,
		AirHumidity:                humiditySummary,
		ForecastAverageProbability: forecastStats.AverageProbability,
	}, nil
}
