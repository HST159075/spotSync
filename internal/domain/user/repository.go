package user

import "gorm.io/gorm"

type Repository interface {
	Create(u *Model) error
	FindByEmail(email string) (*Model, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(u *Model) error {
	return r.db.Create(u).Error
}

func (r *repository) FindByEmail(email string) (*Model, error) {
	var u Model
	err := r.db.Where("email = ?", email).First(&u).Error
	return &u, err
}