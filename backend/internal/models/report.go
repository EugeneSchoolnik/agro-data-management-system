package models

type MetricSummary struct {
	Avg float64 `db:"avg_value" json:"avg"`
	Min float64 `db:"min_value" json:"min"`
	Max float64 `db:"max_value" json:"max"`
}

type MetricAggregate struct {
	SensorType string  `db:"sensor_type"`
	Avg        float64 `db:"avg_value"`
	Min        float64 `db:"min_value"`
	Max        float64 `db:"max_value"`
}

type ForecastStats struct {
	AverageProbability float64 `db:"avg_probability"`
}

type FieldReport struct {
	FieldName                  string        `json:"field_name"`
	Temperature                MetricSummary `json:"temperature"`
	AirHumidity                MetricSummary `json:"air_humidity"`
	ForecastAverageProbability float64       `json:"forecast_average_probability"`
}
