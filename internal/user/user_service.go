package user

import (
	"context"
	"rostekus/golang-backend/internal/auth"
	"rostekus/golang-backend/pkg/util"
	"time"
)

type service struct {
	repository *repository
	timeout    time.Duration
}

func NewService(r *repository) *service {
	return &service{
		repository: r,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	err = util.CheckEmail(u.Email)
	if err != nil {
		return nil, err
	}
	r, err := s.repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserResponse{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	u, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserResponse{}, err
	}
	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserResponse{}, err
	}
	ss, err := auth.GenerateJWT(u.Email, u.Username, u.ID)
	if err != nil {
		return &LoginUserResponse{}, err
	}
	return &LoginUserResponse{AccessToken: ss, User: u.Username}, nil
}
