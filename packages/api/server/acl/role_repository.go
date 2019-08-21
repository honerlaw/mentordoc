package acl

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	"log"
)

type RoleRepository struct {
	server.Repository
}

func NewRoleRepository(db *sql.DB, tx *sql.Tx) *RoleRepository {
	repo := &RoleRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *RoleRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewRoleRepository(repo.Db, tx)
}

func (repo *RoleRepository) Find(name string) *model.Role {
	row := repo.QueryRow(
		"select id, name, created_at, updated_at, deleted_at from role where name = ?",
		name,
	)

	role := &model.Role{}
	err := row.Scan(role.Id, role.Name, role.CreatedAt, role.UpdatedAt, role.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	return role
}

func (repo *RoleRepository) Insert(role *model.Role) (*model.Role, error) {
	existing := repo.Find(role.Name);
	if existing != nil {
		return existing, nil
	}
	role.CreatedAt = server.NowUnix()
	role.UpdatedAt = server.NowUnix()

	_, err := repo.Exec(
		"insert into role (id, name, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?)",
		role.Id,
		role.Name,
		role.CreatedAt,
		role.UpdatedAt,
		role.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert role")
	}

	return role, nil;
}