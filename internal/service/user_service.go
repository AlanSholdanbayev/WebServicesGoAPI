package service

import (
	"context"
	"errors"
	"finalproject/internal/models"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	repo UserRepo
}

type UserRepo interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id int64) error
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{repo: r}
}

// Register создаёт нового пользователя с захешированным паролем
func (s *UserService) Register(ctx context.Context, u *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return s.repo.Create(ctx, u)
}

// Authenticate проверяет email и пароль и возвращает пользователя
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

// FindByID возвращает пользователя по ID
func (s *UserService) FindByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, u *models.User) error {
	// Если пароль передан, хэшируем его
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)
	}
	return s.repo.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
