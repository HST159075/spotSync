package reservation

import (
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"
	"time"
)

type Model struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint       `gorm:"not null" json:"user_id"`
	ZoneID       uint       `gorm:"not null" json:"zone_id"`
	LicensePlate string     `gorm:"not null;size:15" json:"license_plate"`
	Status       string     `gorm:"not null;default:active" json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	User         user.Model `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone         zone.Model `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}

func (Model) TableName() string {
	return "reservations"
}
