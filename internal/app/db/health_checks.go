package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ternaryss/houston/internal/app/types"
)

type healthChecksStore struct {
	db *sql.DB
}

func NewHealthChecksStore(db *sql.DB) *healthChecksStore {
	return &healthChecksStore{
		db: db,
	}
}

func (s *healthChecksStore) Begin() (*types.DbCtx, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	return types.NewDbCtx(tx), nil
}

func (s *healthChecksStore) Commit(ctx *types.DbCtx) error {
	if err := ctx.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *healthChecksStore) Rollback(ctx *types.DbCtx) error {
	if err := ctx.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (s *healthChecksStore) CountByFilter(ftr types.Filter, ctx *types.DbCtx) (int, error) {
	var quantity int
	exec := s.db.QueryRow
	query := `SELECT COUNT(1) FROM HEALTH_CHECKS WHERE WEB_APP_ID = $1`

	if ctx != nil {
		exec = ctx.Tx.QueryRow
	}

	if err := exec(query, ftr.Params["webAppId"]).Scan(&quantity); err != nil {
		return -1, err
	}

	return quantity, nil
}

func (s *healthChecksStore) DeleteByCreatedAtLowerThan(cre time.Time, ctx *types.DbCtx) error {
	exec := s.db.Exec
	query := `DELETE FROM HEALTH_CHECKS WHERE CREATED_AT <= $1`

	if ctx != nil {
		exec = ctx.Tx.Exec
	}

	if _, err := exec(query, cre.Unix()); err != nil {
		return err
	}

	return nil
}

func (s *healthChecksStore) DeleteByWebAppId(wid string, ctx *types.DbCtx) error {
	exec := s.db.Exec
	query := `DELETE FROM HEALTH_CHECKS WHERE WEB_APP_ID = $1`

	if ctx != nil {
		exec = ctx.Tx.Exec
	}

	if _, err := exec(query, wid); err != nil {
		return err
	}

	return nil
}

func (s *healthChecksStore) GetByFilter(ftr types.Filter, pag types.Pagination, ctx *types.DbCtx) ([]*types.HealthCheck, error) {
	exec := s.db.Query
	query := fmt.Sprintf(`SELECT ID, WEB_APP_ID, STATUS, CREATED_AT CREATEDAT, MODIFIED_AT
		FROM HEALTH_CHECKS
		WHERE WEB_APP_ID = $1
		ORDER BY %s LIMIT $2 OFFSET $3`, ftr.Sort)

	if ctx != nil {
		exec = ctx.Tx.Query
	}

	rows, err := exec(query, ftr.Params["webAppId"], pag.Limit, pag.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var collection []*types.HealthCheck

	for rows.Next() {
		var result types.HealthCheck
		var createdAt int64
		var modifiedAt int64

		if err := rows.Scan(
			&result.Id,
			&result.WebAppId,
			&result.Status,
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

func (s *healthChecksStore) GetByWebAppId(wid string, ctx *types.DbCtx) ([]*types.HealthCheck, error) {
	exec := s.db.Query
	query := `SELECT ID, WEB_APP_ID, STATUS, CREATED_AT, MODIFIED_AT
		FROM HEALTH_CHECKS WHERE WEB_APP_ID = $1`

	if ctx != nil {
		exec = ctx.Tx.Query
	}

	rows, err := exec(query, wid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var collection []*types.HealthCheck

	for rows.Next() {
		var result types.HealthCheck
		var createdAt int64
		var modifiedAt int64

		if err := rows.Scan(
			&result.Id,
			&result.WebAppId,
			&result.Status,
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

func (s *healthChecksStore) GetFirstByWebAppIdOrderByCreatedAtDesc(wid string, ctx *types.DbCtx) (*types.HealthCheck, error) {
	var health types.HealthCheck
	var createdAt int64
	var modifiedAt int64
	exec := s.db.QueryRow
	query := `SELECT ID, WEB_APP_ID, STATUS, CREATED_AT, MODIFIED_AT FROM HEALTH_CHECKS
		WHERE WEB_APP_ID = $1 ORDER BY CREATED_AT DESC LIMIT 1`

	if ctx != nil {
		exec = ctx.Tx.QueryRow
	}

	if err := exec(query, wid).Scan(
		&health.Id,
		&health.WebAppId,
		&health.Status,
		&createdAt,
		&modifiedAt,
	); err != nil {
		return nil, err
	}

	health.CreatedAt = time.Unix(createdAt, 0).UTC()
	health.ModifiedAt = time.Unix(modifiedAt, 0).UTC()

	return &health, nil
}

func (s *healthChecksStore) Insert(hck *types.HealthCheck, ctx *types.DbCtx) (*types.HealthCheck, error) {
	var id int64
	exec := s.db.QueryRow
	query := `INSERT INTO HEALTH_CHECKS (WEB_APP_ID, STATUS, CREATED_AT, MODIFIED_AT)
		VALUES ($1, $2, $3, $4) RETURNING ID`

	if ctx != nil {
		exec = ctx.Tx.QueryRow
	}

	if err := exec(
		query,
		hck.WebAppId,
		hck.Status,
		hck.CreatedAt.Unix(),
		hck.ModifiedAt.Unix(),
	).Scan(&id); err != nil {
		return nil, err
	}

	hck.Id = id
	return hck, nil
}
