package user

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
	"strings"
)

type UserRepository struct {
	util.Repository
}

func NewUserRepository(db *sql.DB, tx *sql.Tx) *UserRepository {
	repo := &UserRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo;
}

func (repo *UserRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserRepository(repo.Db, tx)
}

func (repo *UserRepository) Insert(user *shared.User) (*shared.User, error) {
	user.CreatedAt = util.NowUnix()
	user.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into user (id, email, password, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?)",
		user.Id,
		strings.TrimSpace(strings.ToLower(user.Email)),
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert user")
	}

	return user, nil;
}

func (repo *UserRepository) Update(user *shared.User) (*shared.User, error) {
	user.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"update user set email = ?, password = ?, updated_at = ?, deleted_at = ? where id = ?",
		strings.TrimSpace(strings.ToLower(user.Email)),
		user.Password,
		user.UpdatedAt,
		user.DeletedAt,
		user.Id,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to update user")
	}

	return user, nil;
}

func (repo *UserRepository) FindByEmail(email string) *shared.User {
	row := repo.QueryRow(
		"select id, email, password, created_at, updated_at, deleted_at from user where email = ?",
		strings.TrimSpace(strings.ToLower(email)),
	)
	user := &shared.User{}
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}

func (repo *UserRepository) FindById(id string) *shared.User {
	row := repo.QueryRow(
		"select id, email, password, created_at, updated_at, deleted_at from user where id = ?",
		id,
	)
	user := &shared.User{}
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}
