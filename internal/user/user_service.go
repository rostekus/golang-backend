package user

import (
	"context"
	"os"
	"rostekus/golang-backend/util"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = os.Getenv("SECRET_JWT")

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

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserResponse{}, err
	}

	return &LoginUserResponse{AccessToken: ss, User: u.Username}, nil
}
