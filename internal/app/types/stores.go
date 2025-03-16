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
