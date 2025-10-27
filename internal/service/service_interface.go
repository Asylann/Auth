package service

import (
	"context"
	"github.com/Asylann/Auth/internal/model"
)

type ServiceIT interface {
	GetListOfUsers(ctx context.Context) ([]model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserById(ctx context.Context, id int) (model.User, error)
	RegisterUser(ctx context.Context, user model.User) (string, error)
}
