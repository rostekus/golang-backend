package user

import (
	"context"
	"log"
	"rostekus/golang-backend/db"

	"github.com/google/uuid"
)

type repository struct {
	db db.SQLDatabase
}

func NewRepository(db db.SQLDatabase) *repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	userId := uuid.New().String()
	query := "INSERT INTO users(id, username, password, email) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, userId, user.Username, user.Password, user.Email)
	if err != nil {
		log.Fatal(err)
	}
	user.ID = userId
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &User{}, err
	}
	return &u, err
}
