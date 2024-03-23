package database

import (
	"time"

	"gorm.io/gorm"
)

type MeterReadings struct {
	gorm.Model
	ID          string    `gorm:"primaryKey"`
	Nmi         string    `gorm:"size:12;uniqueIndex:idx_nmi_timestamp"`
	Timestamp   time.Time `gorm:"uniqueIndex:idx_nmi_timestamp"`
	Consumption float64
}
