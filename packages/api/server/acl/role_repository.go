package acl

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server"
	"log"
)

type RoleRepository struct {
	server.Repository

	db *sql.DB
	tx *sql.Tx
}

func NewRoleRepository(db *sql.DB, tx *sql.Tx) *RoleRepository {
	repo := &RoleRepository{
		db: db,
		tx: tx,
	}
	return repo
}

func (repo *RoleRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewRoleRepository(repo.db, tx)
}

func (repo *RoleRepository) Find(name string) *Role {
	row := repo.QueryRow(
		"select id, name, created_at, updated_at, deleted_at from role where name = ?",
		name,
	)

	role := &Role{}
	err := row.Scan(role.Id, role.Name, role.CreatedAt, role.UpdatedAt, role.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	return role
}

func (repo *RoleRepository) Insert(role *Role) (*Role, error) {
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