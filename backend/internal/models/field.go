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

type FieldWithCrop struct {
	Field
	CropName    string `db:"crop_name" json:"crop_name"`
	CropVariety string `db:"crop_variety" json:"crop_variety"`
}
