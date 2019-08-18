package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
)

type UserRoleRepository struct {
	server.Repository

	db *sql.DB
	tx *sql.Tx
}

func NewUserRoleRepository(db *sql.DB, tx *sql.Tx) *UserRoleRepository {
	repo := &UserRoleRepository{
		db: db,
		tx: tx,
	}
	return repo
}

func (repo *UserRoleRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserRoleRepository(repo.db, tx)
}

func (repo *UserRoleRepository) Link(user model.User, role *model.Role) error {

	return nil
}

func (repo *UserRoleRepository) Unlink(user model.User, role *model.Role) error {

	return nil
}