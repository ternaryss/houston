package stores

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ternaryss/houston/internal/app/types"
)

type inMemHealthChecksStore struct {
	sequence int64
	data     map[int64]*types.HealthCheck
}

func NewInMemHealthChecksStore() *inMemHealthChecksStore {
	return &inMemHealthChecksStore{
		sequence: 0,
		data:     make(map[int64]*types.HealthCheck),
	}
}

func (s *inMemHealthChecksStore) Begin() (*types.DbCtx, error) {
	return types.NewDbCtx(nil), nil
}

func (s *inMemHealthChecksStore) Commit(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemHealthChecksStore) Rollback(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemHealthChecksStore) DeleteByCreatedAtLowerThan(cre time.Time) error {
	for id, health := range s.data {
		if health.CreatedAt.Before(cre) || health.CreatedAt.Equal(cre) {
			delete(s.data, id)
		}
	}

	return nil
}

func (s *inMemHealthChecksStore) GetFirstByWebAppIdOrderByCreatedAtDesc(wid string) (*types.HealthCheck, error) {
	var latest *types.HealthCheck

	for _, health := range s.data {
		if health.WebAppId != wid {
			continue
		}

		if latest == nil || health.CreatedAt.After(latest.CreatedAt) {
			latest = health
		}
	}

	if latest == nil {
		return nil, sql.ErrNoRows
	}

	return latest, nil
}

func (s *inMemHealthChecksStore) Insert(hck *types.HealthCheck) (*types.HealthCheck, error) {
	s.sequence++

	if _, exists := s.data[s.sequence]; exists {
		return nil, fmt.Errorf("unique constraint violated: %d", s.sequence)
	}

	hck.Id = s.sequence
	s.data[hck.Id] = hck

	return hck, nil
}
