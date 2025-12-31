package models

import "time"

type Pest struct {
	ID             int    `db:"id" json:"id"`
	Name           string `db:"name" json:"name" validate:"required,min=3"`
	ScientificName string `db:"scientific_name" json:"scientific_name" validate:"required,min=5"`
}

type Forecast struct {
	ID             int       `db:"id" json:"id"`
	FieldID        int       `db:"field_id" json:"field_id"`
	PestID         int       `db:"pest_id" json:"pest_id"`
	Probability    float64   `db:"probability" json:"probability"`
	Recommendation string    `db:"recommendation" json:"recommendation"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type ForecastRequest struct {
	CropName string    `json:"crop_name"`
	Variety  string    `json:"variety"`
	Metrics  []float64 `json:"metrics"`
	PestName string    `json:"pest_name"`
}

type ForecastResponse struct {
	Probability    float64 `json:"probability"`
	Recommendation string  `json:"recommendation"`
}
