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
	GetByEmail(eml string, ctx *DbCtx) (*User, error)
	Insert(usr *User, ctx *DbCtx) (*User, error)
}

type WebAppsStore interface {
	dbStore
	CountByFilter(ftr Filter, ctx *DbCtx) (int, error)
	DeleteByIdAndUserEmail(id, usr string, ctx *DbCtx) error
	GetByFilter(ftr Filter, pag Pagination, ctx *DbCtx) ([]*WebApp, error)
	GetByIdAndUserEmail(id, usr string, ctx *DbCtx) (*WebApp, error)
	GetByIntervalOrderByNameAsc(itv string, ctx *DbCtx) ([]*WebApp, error)
	Insert(wap *WebApp, ctx *DbCtx) (*WebApp, error)
	Update(wap *WebApp, ctx *DbCtx) (*WebApp, error)
}

type SubscribersStore interface {
	dbStore
	DeleteByWebAppId(wid string, ctx *DbCtx) error
	DeleteByWebAppIdAndEmail(wid, eml string, ctx *DbCtx) error
	GetByWebAppIdOrderByEmailAsc(wid string, ctx *DbCtx) ([]*Subscriber, error)
	Insert(sub *Subscriber, ctx *DbCtx) (*Subscriber, error)
}

type HealthChecksStore interface {
	dbStore
	CountByFilter(ftr Filter, ctx *DbCtx) (int, error)
	CountByWebAppIdAndNotStatus(wid string, sts int, ctx *DbCtx) (int, error)
	CountByWebAppIdAndStatus(wid string, sts int, ctx *DbCtx) (int, error)
	DeleteByCreatedAtLowerThan(cre time.Time, ctx *DbCtx) error
	DeleteByWebAppId(wid string, ctx *DbCtx) error
	GetByFilter(ftr Filter, pag Pagination, ctx *DbCtx) ([]*HealthCheck, error)
	GetByWebAppId(wid string, ctx *DbCtx) ([]*HealthCheck, error)
	GetFirstByWebAppIdOrderByCreatedAtDesc(wid string, ctx *DbCtx) (*HealthCheck, error)
	Insert(hck *HealthCheck, ctx *DbCtx) (*HealthCheck, error)
}
