package service

import (
	"context"
	"github.com/spinkymax/image-loader/internal/model"
)

type repository interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, modelUser model.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type Controller struct {
	repo repository
}

func NewController(repo repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) AddUser(ctx context.Context, user model.User) error {
	return c.repo.AddUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, id int64) (model.User, error) {
	return c.repo.GetUser(ctx, id)
}

func (c *Controller) UpdateUser(ctx context.Context, user model.User) error {
	return c.repo.UpdateUser(ctx, user)
}

func (c *Controller) DeleteUser(ctx context.Context, id int64) error {
	return c.repo.DeleteUser(ctx, id)
}

func (c *Controller) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return c.repo.GetAllUsers(ctx)
}
