package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user_service/internal/config"
	"user_service/internal/domain"
	"user_service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        *repository.UserRepository
	cfg         config.Config
	minioClient *minio.Client
}

func NewUserService(repo *repository.UserRepository, cfg config.Config, mc *minio.Client) *UserService {
	return &UserService{repo: repo, cfg: cfg, minioClient: mc}
}

func (s *UserService) Register(username, password, fullName string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Username: username,
		Password: string(hashed),
		FullName: fullName,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("Username exists")
	}
	return user, nil
}

func (s *UserService) Login(username, password string) (string, *domain.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("Wrong login information")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("Wrong login information")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(s.cfg.JWTSecret))
	return tokenStr, user, err
}

func (s *UserService) GetProfile(userID uint) (*domain.User, error) {
	return s.repo.FindByID(userID)
}

func (s *UserService) GenerateUploadURL(userID uint, fileName string) (string, error) {
	objectName := fmt.Sprintf("%d_%s", userID, fileName)
	url, err := s.minioClient.PresignedPutObject(
		context.Background(),
		"avatars",
		objectName,
		5*time.Minute,
	)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
