package model

import "github.com/lib/pq"

type Project struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"not null"`
	Images      pq.StringArray `gorm:"type:text[]"`
	Description string         `gorm:"type:text"`
	Tools       pq.StringArray `gorm:"type:text[]"`
}
