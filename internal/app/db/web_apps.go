package db

import (
	"database/sql"
	"fmt"
	"time"

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

func (s *webAppsStore) CountByFilter(ftr types.Filter) (int, error) {
	var quantity int
	query := `SELECT COUNT(1) FROM WEB_APPS WHERE LOWER(USER_EMAIL) = LOWER($1)`

	if err := s.db.QueryRow(query, ftr.Params["userEmail"]).Scan(&quantity); err != nil {
		return -1, err
	}

	return quantity, nil
}

func (s *webAppsStore) DeleteByIdAndUserEmail(id, usr string) error {
	query := `DELETE FROM WEB_APPS WHERE ID = $1 AND LOWER(USER_EMAIL) = LOWER($2)`

	if _, err := s.db.Exec(query, id, usr); err != nil {
		return err
	}

	return nil
}

func (s *webAppsStore) GetByFilter(ftr types.Filter, pag types.Pagination) ([]*types.WebApp, error) {
	query := fmt.Sprintf(`SELECT ID, NAME, URL, STATUS, INTERVAL, USER_EMAIL USEREMAIL, HEALTHY, CREATED_AT CREATEDAT, MODIFIED_AT
		FROM WEB_APPS
		WHERE LOWER(USER_EMAIL) = LOWER($1)
		ORDER BY %s LIMIT $2 OFFSET $3`, ftr.Sort)
	rows, err := s.db.Query(query, ftr.Params["userEmail"], pag.Limit, pag.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var collection []*types.WebApp

	for rows.Next() {
		var result types.WebApp
		var createdAt int64
		var modifiedAt int64

		if err := rows.Scan(
			&result.Id,
			&result.Name,
			&result.Url,
			&result.Status,
			&result.Interval,
			&result.UserEmail,
			&result.Healthy,
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

func (s *webAppsStore) GetByIdAndUserEmail(id, usr string) (*types.WebApp, error) {
	var app types.WebApp
	var createdAt int64
	var modifiedAt int64
	query := `SELECT ID, NAME, URL, STATUS, INTERVAL, USER_EMAIL, HEALTHY, CREATED_AT, MODIFIED_AT FROM WEB_APPS
		WHERE ID = $1 AND LOWER(USER_EMAIL) = LOWER($2)`

	if err := s.db.QueryRow(query, id, usr).Scan(
		&app.Id,
		&app.Name,
		&app.Url,
		&app.Status,
		&app.Interval,
		&app.UserEmail,
		&app.Healthy,
		&createdAt,
		&modifiedAt,
	); err != nil {
		return nil, err
	}

	app.CreatedAt = time.Unix(createdAt, 0).UTC()
	app.ModifiedAt = time.Unix(modifiedAt, 0).UTC()

	return &app, nil
}

func (s *webAppsStore) GetByIntervalOrderByNameAsc(itv string) ([]*types.WebApp, error) {
	query := `SELECT ID, NAME, URL, STATUS, INTERVAL, USER_EMAIL, HEALTHY, CREATED_AT, MODIFIED_AT
		FROM WEB_APPS
		WHERE INTERVAL = $1
		ORDER BY NAME ASC`
	rows, err := s.db.Query(query, itv)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var collection []*types.WebApp

	for rows.Next() {
		var result types.WebApp
		var createdAt int64
		var modifiedAt int64

		if err := rows.Scan(
			&result.Id,
			&result.Name,
			&result.Url,
			&result.Status,
			&result.Interval,
			&result.UserEmail,
			&result.Healthy,
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

func (s *webAppsStore) Insert(wap *types.WebApp) (*types.WebApp, error) {
	var id string
	query := `INSERT INTO WEB_APPS (ID, NAME, URL, STATUS, INTERVAL, USER_EMAIL, HEALTHY, CREATED_AT, MODIFIED_AT)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING ID`

	if err := s.db.QueryRow(
		query,
		uuid.New().String(),
		wap.Name,
		wap.Url,
		wap.Status,
		wap.Interval,
		wap.UserEmail,
		wap.Healthy,
		wap.CreatedAt.Unix(),
		wap.ModifiedAt.Unix(),
	).Scan(&id); err != nil {
		return nil, err
	}

	wap.Id = id
	return wap, nil
}

func (s *webAppsStore) Update(wap *types.WebApp) (*types.WebApp, error) {
	query := `UPDATE WEB_APPS SET NAME = $1, URL = $2, STATUS = $3, INTERVAL = $4, USER_EMAIL = $5, HEALTHY = $6, MODIFIED_AT = $7
		WHERE ID = $8 AND USER_EMAIL = $5`

	if _, err := s.db.Exec(
		query,
		wap.Name,
		wap.Url,
		wap.Status,
		wap.Interval,
		wap.UserEmail,
		wap.Healthy,
		time.Now().UTC().Unix(),
		wap.Id,
	); err != nil {
		return nil, err
	}

	return wap, nil
}
