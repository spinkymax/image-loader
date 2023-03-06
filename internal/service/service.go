package service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
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
	CheckTgAuth(ctx context.Context, tgID int64) (int, error)
	AuthorizeTG(ctx context.Context, userID int, telegramID int64) error
}

type imageRepository interface {
	AddImage(ctx context.Context, modelImage model.Image) error
	GetImages(ctx context.Context, userID int) ([]model.Image, error)
}

type fileStorage interface {
	PutObject(ctx context.Context, image model.Image) error
	GetUrls(ctx context.Context, images []model.Image) ([]string, error)
	GetObjects(ctx context.Context, images []model.Image) ([]io.Reader, error)
}

type Controller struct {
	repo      repository
	imageRepo imageRepository
	cfg       *config.Config
	minio     fileStorage
}

func NewController(repo repository, imageRepo imageRepository, cfg *config.Config, m fileStorage) *Controller {
	return &Controller{
		repo:      repo,
		imageRepo: imageRepo,
		cfg:       cfg,
		minio:     m,
	}
}

func (c *Controller) AddUser(ctx context.Context, user model.User) error {
	return c.repo.AddUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := c.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	images, err := c.imageRepo.GetImages(ctx, user.ID)
	if err != nil {
		return model.User{}, err
	}

	urls, err := c.minio.GetUrls(ctx, images)
	if err != nil {
		return model.User{}, err
	}

	user.ImageUrls = urls

	return user, nil
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

	users, err := c.repo.GetAllUsers(ctx)
	if err != nil {
		return []model.User{}, err
	}

	for i, user := range users {
		images, err := c.imageRepo.GetImages(ctx, user.ID)

		if err != nil {
			return []model.User{}, err
		}
		urls, err := c.minio.GetUrls(ctx, images)
		if err != nil {
			return []model.User{}, err
		}
		users[i].ImageUrls = urls

	}
	return users, err
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

func (c *Controller) AddFile(ctx context.Context, image model.Image) error {
	imageName, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to generate image name: %w", err)
	}

	image.Name = imageName.String() + image.Extension

	err = c.imageRepo.AddImage(ctx, image)
	if err != nil {
		return fmt.Errorf("failed to save image data to db: %w", err)
	}

	err = c.minio.PutObject(ctx, image)
	if err != nil {
		return fmt.Errorf("failed to put image to fileStore: %w", err)
	}

	return err
}
func (c *Controller) AuthorizeTG(ctx context.Context, tgID int64, login, password string) error {
	user, err := c.repo.CheckAuth(ctx, login, password)
	if err != nil {
		return err
	}

	_, err = c.repo.CheckTgAuth(ctx, tgID)
	if err == nil {
		return err
	}

	err = c.repo.AuthorizeTG(ctx, user.ID, tgID)
	if err != nil {
		return err
	}

	return nil
}
func (c *Controller) GetImageObjects(ctx context.Context, tgID int64) ([]io.Reader, error) {
	userID, err := c.repo.CheckTgAuth(ctx, tgID)
	if err != nil {
		return nil, err
	}

	images, err := c.imageRepo.GetImages(ctx, userID)
	if err != nil {
		return nil, err
	}

	objects, err := c.minio.GetObjects(ctx, images)
	if err != nil {
		return nil, err
	}

	return objects, nil
}
