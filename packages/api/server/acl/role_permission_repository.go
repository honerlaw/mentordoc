package acl

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	"log"
)

type RolePermissionRepository struct {
	server.Repository

	db *sql.DB
	tx *sql.Tx
}

func NewRolePermissionRepository(db *sql.DB, tx *sql.Tx) *RolePermissionRepository {
	repo := &RolePermissionRepository{
		db: db,
		tx: tx,
	}
	return repo
}

func (repo *RolePermissionRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewRolePermissionRepository(repo.db, tx)
}

func (repo *RolePermissionRepository) Link(role *model.Role, permission *model.Permission) error {

	// check if its already been linked
	rows, err := repo.Query(
		"select role_id, permission_id from role_permission where role_id = ? and permission_id = ?",
		role.Id,
		permission.Id,
	)
	if err != nil {
		log.Print(err)
		return errors.New("failed to link permission to role")
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count += 1
	}

	// permission is already linked, so do nothing
	if count > 0 {
		return nil
	}

	// otherwise attempt to link the permission to the role
	_, err = repo.Exec(
		"insert into role_permission (role_id, permission_id) values (?, ?)",
		role.Id,
		permission.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to link permission to role")
	}

	return nil
}