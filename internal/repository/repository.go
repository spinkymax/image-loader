package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/spinkymax/image-loader/internal/config"
	"github.com/spinkymax/image-loader/internal/model"
)

type user struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Login       string         `db:"login"`
	Password    string         `db:"password"`
	Description sql.NullString `db:"description"`
}

type UserRepo struct {
	db  *sqlx.DB
	cfg *config.DB
}

func NewUserRepo(db *sqlx.DB, cfg *config.DB) *UserRepo {
	return &UserRepo{
		db:  db,
		cfg: cfg,
	}
}

func (u *UserRepo) RunMigrations() error {
	driver, err := postgres.WithInstance(u.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get migration tool driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		u.cfg.Driver, driver)
	if err != nil {
		return fmt.Errorf("failed to connect migration tool: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (u *UserRepo) AddUser(ctx context.Context, modelUser model.User) error {
	query := `INSERT INTO users(name, description, login, password) VALUES (:name, :description, :login, :password)`

	user := convertUser(modelUser)

	_, err := u.db.NamedExecContext(ctx, query, &user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (u *UserRepo) GetUser(ctx context.Context, id int64) (model.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var us user

	row := u.db.QueryRowxContext(ctx, query, id)

	err := row.StructScan(&us)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to scan user: %w", err)
	}

	return us.toModel(), nil

}

func (u *UserRepo) UpdateUser(ctx context.Context, modelUser model.User) error {
	query := `UPDATE users set (name, description, login, password) = (:name, :description, :login, :password) 
             WHERE id = :id`

	_, err := u.db.NamedExecContext(ctx, query, convertUser(modelUser))
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users where id = $1`

	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context) ([]model.User, error) {

	users := make([]model.User, 0)

	rows, err := u.db.Queryx("SELECT * FROM users")

	if err != nil {
		return []model.User{}, fmt.Errorf("failed to scan users: %w", err)
	}

	for rows.Next() {
		var userEntity user

		err = rows.StructScan(&userEntity)
		users = append(users, userEntity.toModel())
	}

	err = rows.Err()
	if err != nil {
		return []model.User{}, fmt.Errorf("failed to mapping users: %w", err)
	}

	return users, nil
}

func convertUser(modelUser model.User) user {
	return user{
		ID:       modelUser.ID,
		Name:     modelUser.Name,
		Login:    modelUser.Login,
		Password: modelUser.Password,
		Description: sql.NullString{
			String: modelUser.Description,
			Valid:  true,
		},
	}
}

func (u user) toModel() model.User {
	return model.User{
		ID:          u.ID,
		Name:        u.Name,
		Login:       u.Login,
		Password:    u.Password,
		Description: u.Description.String,
	}
}
