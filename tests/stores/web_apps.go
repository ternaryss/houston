package stores

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

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

func (s *inMemWebAppsStore) GetByIntervalOrderByNameAsc(itv string) ([]*types.WebApp, error) {
	var result []*types.WebApp

	for _, app := range s.data {
		if app.Interval == itv {
			result = append(result, app)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	return result, nil
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

func (s *inMemWebAppsStore) Update(wap *types.WebApp) (*types.WebApp, error) {
	app, exists := s.data[wap.Id]

	if !exists || !strings.EqualFold(app.UserEmail, wap.UserEmail) {
		return wap, nil
	}

	app.Name = wap.Name
	app.Url = wap.Url
	app.Status = wap.Status
	app.Interval = wap.Interval
	app.UserEmail = wap.UserEmail
	app.Healthy = wap.Healthy
	app.ModifiedAt = time.Now().UTC()

	return app, nil
}
