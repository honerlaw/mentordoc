package acl

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
)

type PermissionRepository struct {
	util.Repository
}

func NewPermissionRepository(db *sql.DB, tx *sql.Tx) *PermissionRepository {
	repo := &PermissionRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *PermissionRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewPermissionRepository(repo.Db, tx)
}

func (repo *PermissionRepository) Find(resourcePath string, action string) *Permission {
	row := repo.QueryRow(
		"select id, resource_path, action, created_at, updated_at, deleted_at from permission where resource_path = ? and action = ?",
		resourcePath,
		action,
	)

	var permission Permission
	err := row.Scan(&permission.Id, &permission.ResourcePath, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt, &permission.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &permission
}

func (repo *PermissionRepository) Insert(permission *Permission) (*Permission, error) {
	existing := repo.Find(permission.ResourcePath, permission.Action)
	if existing != nil {
		return existing, nil;
	}

	permission.CreatedAt = util.NowUnix()
	permission.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into permission (id, resource_path, action, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?)",
		permission.Id,
		permission.ResourcePath,
		permission.Action,
		permission.CreatedAt,
		permission.UpdatedAt,
		permission.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert permission")
	}

	return permission, nil;
}