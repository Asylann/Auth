package repository

import (
	"context"
	"github.com/Asylann/Auth/internal/model"
)

type RepositoryIT interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GetUserById(ctx context.Context, id int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetListOfUsers(ctx context.Context) ([]model.User, error)
}
