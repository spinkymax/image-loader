package repository

import (
	"context"
	"fmt"
)

type tgAuth struct {
	ID         int   `db:"id"`
	UserID     int   `db:"user_id"`
	TelegramID int64 `db:"telegram_id"`
}

func (u *UserRepo) CheckTgAuth(ctx context.Context, tgID int64) (int, error) {
	query := `SELECT * 
              FROM tg_auth tga 
              WHERE tga.telegram_id = $1`

	var tgAuth tgAuth

	row := u.db.QueryRowxContext(ctx, query, tgID)

	err := row.StructScan(&tgAuth)
	if err != nil {
		return 0, fmt.Errorf("failed to tg auth: %w", err)
	}

	return tgAuth.UserID, nil
}

func (u *UserRepo) AuthorizeTG(ctx context.Context, userID int, telegramID int64) error {
	query := `INSERT INTO tg_auth(user_id, telegram_id) VALUES ($1, $2)`

	_, err := u.db.ExecContext(ctx, query, userID, telegramID)
	if err != nil {
		return fmt.Errorf("failed to tg auth: %w", err)
	}

	return nil
}
