package repository

import (
	"context"
	"finalproject/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id int64) error
}

type PostgresUserRepo struct {
	db *sqlx.DB
}

func NewPostgresUserRepo(db *sqlx.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, u *models.User) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO users (email, password_hash, name) VALUES (:email, :password_hash, :name)`, map[string]interface{}{
		"email":         u.Email,
		"password_hash": u.Password,
		"name":          u.Name,
	})
	return err
}

func (r *PostgresUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, "SELECT id, email, password_hash as password, name, created_at FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepo) FindByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, "SELECT id, email, password_hash as password, name, created_at FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepo) Update(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET name=$1 WHERE id=$2", u.Name, u.ID)
	return err
}

func (r *PostgresUserRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}
