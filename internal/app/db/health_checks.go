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

func (s *healthChecksStore) CountByFilter(ftr types.Filter) (int, error) {
	var quantity int
	query := `SELECT COUNT(1) FROM HEALTH_CHECKS WHERE WEB_APP_ID = $1`

	if err := s.db.QueryRow(query, ftr.Params["webAppId"]).Scan(&quantity); err != nil {
		return -1, err
	}

	return quantity, nil
}

func (s *healthChecksStore) DeleteByCreatedAtLowerThan(cre time.Time) error {
	query := `DELETE FROM HEALTH_CHECKS WHERE CREATED_AT <= $1`

	if _, err := s.db.Exec(query, cre.Unix()); err != nil {
		return err
	}

	return nil
}

func (s *healthChecksStore) DeleteByWebAppId(wid string) error {
	query := `DELETE FROM HEALTH_CHECKS WHERE WEB_APP_ID = $1`

	if _, err := s.db.Exec(query, wid); err != nil {
		return err
	}

	return nil
}

func (s *healthChecksStore) GetByFilter(ftr types.Filter, pag types.Pagination) ([]*types.HealthCheck, error) {
	query := fmt.Sprintf(`SELECT ID, WEB_APP_ID, STATUS, CREATED_AT CREATEDAT, MODIFIED_AT
		FROM HEALTH_CHECKS
		WHERE WEB_APP_ID = $1
		ORDER BY %s LIMIT $2 OFFSET $3`, ftr.Sort)
	rows, err := s.db.Query(query, ftr.Params["webAppId"], pag.Limit, pag.Offset)

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

func (s *healthChecksStore) GetByWebAppId(wid string) ([]*types.HealthCheck, error) {
	query := `SELECT ID, WEB_APP_ID, STATUS, CREATED_AT, MODIFIED_AT
		FROM HEALTH_CHECKS WHERE WEB_APP_ID = $1`
	rows, err := s.db.Query(query, wid)

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

func (s *healthChecksStore) GetFirstByWebAppIdOrderByCreatedAtDesc(wid string) (*types.HealthCheck, error) {
	var health types.HealthCheck
	var createdAt int64
	var modifiedAt int64
	query := `SELECT ID, WEB_APP_ID, STATUS, CREATED_AT, MODIFIED_AT FROM HEALTH_CHECKS
		WHERE WEB_APP_ID = $1 ORDER BY CREATED_AT DESC LIMIT 1`

	if err := s.db.QueryRow(query, wid).Scan(
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

func (s *healthChecksStore) Insert(hck *types.HealthCheck) (*types.HealthCheck, error) {
	var id int64
	query := `INSERT INTO HEALTH_CHECKS (WEB_APP_ID, STATUS, CREATED_AT, MODIFIED_AT)
		VALUES ($1, $2, $3, $4) RETURNING ID`

	if err := s.db.QueryRow(
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
