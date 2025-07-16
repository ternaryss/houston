package stores

import (
	"database/sql"
	"fmt"

	"github.com/ternaryss/houston/internal/app/types"
)

type inMemUsersStore struct {
	data map[string]*types.User
}

func NewInMemUsersStore() *inMemUsersStore {
	return &inMemUsersStore{
		data: make(map[string]*types.User),
	}
}

func (s *inMemUsersStore) Begin() (*types.DbCtx, error) {
	return types.NewDbCtx(nil), nil
}

func (s *inMemUsersStore) Commit(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemUsersStore) Rollback(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemUsersStore) GetByEmail(eml string, ctx *types.DbCtx) (*types.User, error) {
	user, exists := s.data[eml]

	if !exists {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

func (s *inMemUsersStore) Insert(usr *types.User, ctx *types.DbCtx) (*types.User, error) {
	if _, exists := s.data[usr.Email]; exists {
		return nil, fmt.Errorf("unique constraint violated: %s", usr.Email)
	}

	s.data[usr.Email] = usr

	return usr, nil
}
