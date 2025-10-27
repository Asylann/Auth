package repository

import (
	"github.com/Asylann/Auth/internal/model"
	"golang.org/x/net/context"
)

func (repo *Repository) CreateUser(ctx context.Context, user model.User) (string, error) {
	if user.Role == "" {
		user.Role = "user"
	}
	var id string
	err := repo.Pool.QueryRow(ctx, "INSERT INTO users(email, password,role) VALUES($1,$2,$3) RETURNING id", user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, err
}

func (repo *Repository) GetUserById(ctx context.Context, id int) (model.User, error) {
	var user model.User
	err := repo.Pool.QueryRow(ctx, "SELECT * FROM users WHERE id=$1", id).Scan(&user.ID, &user.Email, &user.Password, &user.CreateAt, &user.Role)

	if err != nil {
		return model.User{}, err
	}

	return user, err
}

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := repo.Pool.QueryRow(ctx, "SELECT * FROM users WHERE email=$1", email).Scan(&user.ID, &user.Email, &user.Password, &user.CreateAt, &user.Role)
	if err != nil {
		return user, err
	}

	return user, err
}

func (repo *Repository) GetListOfUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	rows, err := repo.Pool.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreateAt, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
