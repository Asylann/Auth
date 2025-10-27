package service

import (
	"context"
	"github.com/Asylann/Auth/internal/model"
	"github.com/Asylann/Auth/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	Repo   repository.Repository
	Logger *logrus.Logger
}

func NewService(repo repository.Repository, logger *logrus.Logger) Service {
	return Service{Logger: logger, Repo: repo}
}

func (service *Service) RegisterUser(ctx context.Context, user model.User) (string, error) {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return "", err
	}

	user.Password = string(HashedPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	id, err := service.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (service *Service) GetUserById(ctx context.Context, id int) (model.User, error) {
	user, err := service.Repo.GetUserById(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (service *Service) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := service.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service *Service) GetListOfUsers(ctx context.Context) ([]model.User, error) {
	users, err := service.Repo.GetListOfUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, err
}
