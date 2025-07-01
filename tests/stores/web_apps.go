package stores

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ternaryss/houston/internal/app/types"
)

type inMemWebAppsStore struct {
	data map[string]*types.WebApp
}

func NewInMemWebAppsStore() *inMemWebAppsStore {
	return &inMemWebAppsStore{
		data: make(map[string]*types.WebApp),
	}
}

func (s *inMemWebAppsStore) Begin() (*types.DbCtx, error) {
	return types.NewDbCtx(nil), nil
}

func (s *inMemWebAppsStore) Commit(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemWebAppsStore) Rollback(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemWebAppsStore) CountByFilter(ftr types.Filter) (int, error) {
	return -1, nil
}

func (s *inMemWebAppsStore) DeleteByIdAndUserEmail(id, usr string) error {
	app, exists := s.data[id]

	if !exists || !strings.EqualFold(app.UserEmail, usr) {
		return nil
	}

	delete(s.data, app.Id)

	return nil
}

func (s *inMemWebAppsStore) GetByFilter(ftr types.Filter, pag types.Pagination) ([]*types.WebApp, error) {
	return []*types.WebApp{}, nil
}

func (s *inMemWebAppsStore) GetByIdAndUserEmail(id, usr string) (*types.WebApp, error) {
	app, exists := s.data[id]

	if !exists || !strings.EqualFold(app.UserEmail, usr) {
		return nil, sql.ErrNoRows
	}

	return app, nil
}

func (s *inMemWebAppsStore) Insert(wap *types.WebApp) (*types.WebApp, error) {
	id := uuid.New().String()

	if _, exists := s.data[id]; exists {
		return nil, fmt.Errorf("unique constraint violated: %s", id)
	}

	wap.Id = id
	s.data[id] = wap

	return wap, nil
}
