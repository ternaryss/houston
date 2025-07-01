package types

import "database/sql"

type DbCtx struct {
	Tx *sql.Tx
}

func NewDbCtx(tx *sql.Tx) *DbCtx {
	return &DbCtx{Tx: tx}
}

type dbStore interface {
	Begin() (*DbCtx, error)
	Commit(ctx *DbCtx) error
	Rollback(ctx *DbCtx) error
}

type UsersStore interface {
	dbStore
	GetByEmail(eml string) (*User, error)
	Insert(usr *User) (*User, error)
}

type WebAppsStore interface {
	dbStore
	CountByFilter(ftr Filter) (int, error)
	DeleteByIdAndUserEmail(id, usr string) error
	GetByFilter(ftr Filter, pag Pagination) ([]*WebApp, error)
	GetByIdAndUserEmail(id, usr string) (*WebApp, error)
	Insert(wap *WebApp) (*WebApp, error)
}
