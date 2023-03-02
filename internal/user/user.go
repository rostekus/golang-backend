package user

import "context"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Email    string `json:"email"  binding:"required"`
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserRequest struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
type LoginUserResponse struct {
	AccessToken string `json:"token"`
	User        string `json:"user"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	Login(c context.Context, req *LoginUserRequest) (*LoginUserResponse, error)
}
