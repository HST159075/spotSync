package zone

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(z *Model) error
	FindAll() ([]Model, error)
	FindByID(id uint) (*Model, error)
	Update(z *Model) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(z *Model) error {
	return r.db.Create(z).Error
}

func (r *repository) FindAll() ([]Model, error) {
	var zones []Model
	err := r.db.Find(&zones).Error
	return zones, err
}

func (r *repository) FindByID(id uint) (*Model, error) {
	var z Model
	err := r.db.First(&z, id).Error
	return &z, err
}

func (r *repository) Update(z *Model) error {
	return r.db.Save(z).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Model{}, id).Error
}

func (r *repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM reservations 
		WHERE zone_id = ? AND status = 'active'
	`, zoneID).Scan(&count).Error
	return count, err
}
