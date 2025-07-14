package types

import (
	"database/sql"
	"time"
)

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
	GetByIntervalOrderByNameAsc(itv string) ([]*WebApp, error)
	Insert(wap *WebApp) (*WebApp, error)
	Update(wap *WebApp) (*WebApp, error)
}

type SubscribersStore interface {
	dbStore
	DeleteByWebAppId(wid string) error
	GetByWebAppIdOrderByEmailAsc(wid string) ([]*Subscriber, error)
	Insert(sub *Subscriber) (*Subscriber, error)
}

type HealthChecksStore interface {
	dbStore
	CountByFilter(ftr Filter) (int, error)
	DeleteByCreatedAtLowerThan(cre time.Time) error
	DeleteByWebAppId(wid string) error
	GetByFilter(ftr Filter, pag Pagination) ([]*HealthCheck, error)
	GetByWebAppId(wid string) ([]*HealthCheck, error)
	GetFirstByWebAppIdOrderByCreatedAtDesc(wid string) (*HealthCheck, error)
	Insert(hck *HealthCheck) (*HealthCheck, error)
}
