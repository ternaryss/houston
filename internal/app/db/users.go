package db

import (
	"database/sql"
	"time"

	"github.com/ternaryss/houston/internal/app/types"
)

type usersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *usersStore {
	return &usersStore{
		db: db,
	}
}

func (s *usersStore) Begin() (*types.DbCtx, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	return types.NewDbCtx(tx), nil
}

func (s *usersStore) Commit(ctx *types.DbCtx) error {
	if err := ctx.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *usersStore) Rollback(ctx *types.DbCtx) error {
	if err := ctx.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (s *usersStore) GetByEmail(eml string, ctx *types.DbCtx) (*types.User, error) {
	var user types.User
	var createdAt int64
	var modifiedAt int64
	exec := s.db.QueryRow
	query := `SELECT EMAIL, "PASSWORD", CREATED_AT, MODIFIED_AT FROM USERS
				WHERE LOWER(EMAIL) = LOWER($1)`

	if ctx != nil {
		exec = ctx.Tx.QueryRow
	}

	if err := exec(query, eml).Scan(
		&user.Email,
		&user.Password,
		&createdAt,
		&modifiedAt,
	); err != nil {
		return nil, err
	}

	user.CreatedAt = time.Unix(createdAt, 0).UTC()
	user.ModifiedAt = time.Unix(modifiedAt, 0).UTC()

	return &user, nil
}

func (s *usersStore) Insert(usr *types.User, ctx *types.DbCtx) (*types.User, error) {
	exec := s.db.Exec
	query := `INSERT INTO USERS (EMAIL, "PASSWORD", CREATED_AT, MODIFIED_AT) VALUES ($1, $2, $3, $4)`

	if ctx != nil {
		exec = ctx.Tx.Exec
	}

	if _, err := exec(query, usr.Email, usr.Password, usr.CreatedAt.Unix(), usr.ModifiedAt.Unix()); err != nil {
		return nil, err
	}

	return usr, nil
}
