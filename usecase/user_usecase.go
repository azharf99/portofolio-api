package usecase

import (
	"errors"
	"time"

	"github.com/azharf99/portofolio-api/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewUserUsecase(repo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	return &userUsecase{repo, jwtSecret}
}

func (u *userUsecase) Login(username, password string) (string, error) {
	user, err := u.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("username atau password salah")
	}

	// Buat JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(u.jwtSecret))
}

func (u *userUsecase) Update(id uint, user *domain.User) error {
	// Jika password ikut diupdate, kita WAJIB melakukan hashing ulang demi keamanan!
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("gagal mengenkripsi password baru")
		}
		user.Password = string(hashedPassword)
	}
	return u.repo.Update(id, user)
}

func (u *userUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
