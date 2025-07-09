package stores

import (
	"fmt"
	"sort"

	"github.com/ternaryss/houston/internal/app/types"
)

type inMemSubscribersStore struct {
	sequence int64
	data     map[int64]*types.Subscriber
}

func NewInMemSubscribersStore() *inMemSubscribersStore {
	return &inMemSubscribersStore{
		sequence: 0,
		data:     make(map[int64]*types.Subscriber),
	}
}

func (s *inMemSubscribersStore) Begin() (*types.DbCtx, error) {
	return types.NewDbCtx(nil), nil
}

func (s *inMemSubscribersStore) Commit(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemSubscribersStore) Rollback(ctx *types.DbCtx) error {
	return nil
}

func (s *inMemSubscribersStore) DeleteByWebAppId(wid string) error {
	for id, subscriber := range s.data {
		if subscriber.WebAppId == wid {
			delete(s.data, id)
		}
	}

	return nil
}

func (s *inMemSubscribersStore) GetByWebAppIdOrderByEmailAsc(wid string) ([]*types.Subscriber, error) {
	var collection []*types.Subscriber

	for _, subscriber := range s.data {
		if subscriber.WebAppId == wid {
			collection = append(collection, subscriber)
		}
	}

	sort.Slice(collection, func(i, j int) bool {
		return collection[i].Email < collection[j].Email
	})

	return collection, nil
}

func (s *inMemSubscribersStore) Insert(sub *types.Subscriber) (*types.Subscriber, error) {
	s.sequence++

	if _, exists := s.data[s.sequence]; exists {
		return nil, fmt.Errorf("unique constraint violated: %d", s.sequence)
	}

	sub.Id = s.sequence
	s.data[sub.Id] = sub

	return sub, nil
}
