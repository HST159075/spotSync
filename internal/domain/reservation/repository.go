package reservation

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("zone is full")

type Repository interface {
	CreateWithLock(userID uint, zoneID uint, licensePlate string) (*Model, error)
	FindByUserID(userID uint) ([]Model, error)
	FindByID(id uint) (*Model, error)
	Cancel(id uint) error
	FindAll() ([]Model, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateWithLock(userID uint, zoneID uint, licensePlate string) (*Model, error) {
	var res Model
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone struct {
			TotalCapacity int
		}
		if err := tx.Table("parking_zones").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("total_capacity").
			Where("id = ?", zoneID).
			First(&zone).Error; err != nil {
			return errors.New("zone not found")
		}

		var count int64
		if err := tx.Model(&Model{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&count).Error; err != nil {
			return err
		}

		if int(count) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		res = Model{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}
		return tx.Create(&res).Error
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) FindByUserID(userID uint) ([]Model, error) {
	var reservations []Model
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Find(&reservations).Error
	return reservations, err
}

func (r *repository) FindByID(id uint) (*Model, error) {
	var res Model
	err := r.db.First(&res, id).Error
	return &res, err
}

func (r *repository) Cancel(id uint) error {
	return r.db.Model(&Model{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (r *repository) FindAll() ([]Model, error) {
	var reservations []Model
	err := r.db.Preload("Zone").Preload("User").
		Find(&reservations).Error
	return reservations, err
}
