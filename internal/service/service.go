package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7"
	"github.com/spinkymax/image-loader/internal/config"
	"github.com/spinkymax/image-loader/internal/constants"
	"github.com/spinkymax/image-loader/internal/model"
	"io"
	"strconv"
	"time"
)

type repository interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, modelUser model.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	CheckAuth(ctx context.Context, login, password string) (model.User, error)
}

type Controller struct {
	repo  repository
	cfg   *config.Config
	minio *minio.Client
}

func NewController(repo repository, cfg *config.Config, m *minio.Client) *Controller {
	return &Controller{
		repo:  repo,
		cfg:   cfg,
		minio: m,
	}
}

func (c *Controller) AddUser(ctx context.Context, user model.User) error {
	return c.repo.AddUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, id int64) (model.User, error) {
	return c.repo.GetUser(ctx, id)
}

func (c *Controller) UpdateUser(ctx context.Context, user model.User) error {
	id := ctx.Value(constants.IdCtxKey)

	if id != user.ID {
		return fmt.Errorf("user do not match")
	}
	return c.repo.UpdateUser(ctx, user)
}

func (c *Controller) DeleteUser(ctx context.Context, id int64) error {
	return c.repo.DeleteUser(ctx, id)
}

func (c *Controller) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return c.repo.GetAllUsers(ctx)
}

func (c *Controller) Authorize(ctx context.Context, login, password string) (string, error) {
	user, err := c.repo.CheckAuth(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authorize user: %w", err)
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		Subject:   "authorized",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        strconv.Itoa(int(user.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.cfg.JWTKeyword))
	if err != nil {
		return "", fmt.Errorf("failed to sigh token: %w", err)
	}
	return tokenString, nil
}

func (c *Controller) AddFile(ctx context.Context, filename string, file io.Reader) error {
	_, err := c.minio.PutObject(ctx, c.cfg.Minio.Bucket, filename, file, -1, minio.PutObjectOptions{})

	return err
}
