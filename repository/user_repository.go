package repository

import (
	"github.com/azharf99/portofolio-api/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepository) GetByID(id uint) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepository) Update(id uint, user *domain.User) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
