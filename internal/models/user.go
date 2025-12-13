package models

import "time"

type User struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email" validate:"required,email"`
	Password  string    `db:"password,omitempty" json:"password,omitempty" validate:"required,min=6"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
