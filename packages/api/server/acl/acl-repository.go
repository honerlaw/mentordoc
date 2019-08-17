package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/transaction"
)

type AclRepository struct {
	transaction.Transactionable

	db *sql.DB
	tx *sql.Tx
}

func NewAclRepository(db *sql.DB, tx *sql.Tx) *AclRepository {
	repo := &AclRepository{
		db: db,
		tx: tx,
	}
	repo.CloneWithTransaction = func(tx *sql.Tx) interface{} {
		return NewAclRepository(repo.db, tx)
	}
	return repo
}

func (*AclRepository) Testing() {

}
