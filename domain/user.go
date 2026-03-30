package domain

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUsecase interface {
	Login(username, password string) (string, error)
	Update(id uint, user *User) error // BARU
	Delete(id uint) error             // BARU
}

type UserRepository interface {
	GetByUsername(username string) (User, error)
	GetByID(id uint) (User, error)    // BARU: Untuk mengecek user sebelum diupdate
	Update(id uint, user *User) error // BARU
	Delete(id uint) error             // BARU
}
