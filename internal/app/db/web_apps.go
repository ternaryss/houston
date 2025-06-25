package db

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/ternaryss/houston/internal/app/types"
)

type webAppsStore struct {
	db *sql.DB
}

func NewWebAppsStore(db *sql.DB) *webAppsStore {
	return &webAppsStore{
		db: db,
	}
}

func (s *webAppsStore) Begin() (*types.DbCtx, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	return types.NewDbCtx(tx), nil
}

func (s *webAppsStore) Commit(ctx *types.DbCtx) error {
	if err := ctx.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *webAppsStore) Rollback(ctx *types.DbCtx) error {
	if err := ctx.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (s *webAppsStore) Insert(wap *types.WebApp) (*types.WebApp, error) {
	var id string
	query := `INSERT INTO WEB_APPS (ID, NAME, URL, USER_EMAIL, CREATED_AT, MODIFIED_AT) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID`

	if err := s.db.QueryRow(
		query,
		uuid.New().String(),
		wap.Name,
		wap.Url,
		wap.UserEmail,
		wap.CreatedAt.Unix(),
		wap.ModifiedAt.Unix(),
	).Scan(&id); err != nil {
		return nil, err
	}

	wap.Id = id
	return wap, nil
}
