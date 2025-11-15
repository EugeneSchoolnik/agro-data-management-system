package models

import "time"

type Crop struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Variety     string `db:"variety" json:"variety"`
	Description string `db:"description" json:"description"`
}

type Field struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Area      float64   `db:"area" json:"area"`
	Location  string    `db:"location" json:"location"`
	CropID    int       `db:"crop_id" json:"crop_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}