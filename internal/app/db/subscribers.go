package db

import (
	"database/sql"
	"time"

	"github.com/ternaryss/houston/internal/app/types"
)

type subscribersStore struct {
	db *sql.DB
}

func NewSubscribersStore(db *sql.DB) *subscribersStore {
	return &subscribersStore{
		db: db,
	}
}

func (s *subscribersStore) Begin() (*types.DbCtx, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	return types.NewDbCtx(tx), nil
}

func (s *subscribersStore) Commit(ctx *types.DbCtx) error {
	if err := ctx.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *subscribersStore) Rollback(ctx *types.DbCtx) error {
	if err := ctx.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (s *subscribersStore) DeleteByWebAppId(wid string) error {
	query := `DELETE FROM SUBSCRIBERS WHERE WEB_APP_ID = $1`

	if _, err := s.db.Exec(query, wid); err != nil {
		return err
	}

	return nil
}

func (s *subscribersStore) GetByWebAppIdOrderByEmailAsc(wid string) ([]*types.Subscriber, error) {
	query := `SELECT ID, WEB_APP_ID, EMAIL, CREATED_AT, MODIFIED_AT
		FROM SUBSCRIBERS WHERE WEB_APP_ID = $1 ORDER BY EMAIL ASC`
	rows, err := s.db.Query(query, wid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var collection []*types.Subscriber

	for rows.Next() {
		var result types.Subscriber
		var createdAt int64
		var modifiedAt int64

		if err := rows.Scan(
			&result.Id,
			&result.WebAppId,
			&result.Email,
			&createdAt,
			&modifiedAt,
		); err != nil {
			return nil, err
		}

		result.CreatedAt = time.Unix(createdAt, 0).UTC()
		result.ModifiedAt = time.Unix(modifiedAt, 0).UTC()
		collection = append(collection, &result)
	}

	return collection, nil
}

func (s *subscribersStore) Insert(sub *types.Subscriber) (*types.Subscriber, error) {
	var id int64
	query := `INSERT INTO SUBSCRIBERS (WEB_APP_ID, EMAIL, CREATED_AT, MODIFIED_AT)
		VALUES ($1, $2, $3, $4) RETURNING ID`

	if err := s.db.QueryRow(
		query,
		sub.WebAppId,
		sub.Email,
		sub.CreatedAt.Unix(),
		sub.ModifiedAt.Unix(),
	).Scan(&id); err != nil {
		return nil, err
	}

	sub.Id = id
	return sub, nil
}
