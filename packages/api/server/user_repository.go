package server

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/model"
	"log"
	"strings"
)

type UserRepository struct {
	Repository
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

func (repo *UserRepository) Insert(user *model.User) (*model.User, error) {
	user.CreatedAt = NowUnix()
	user.UpdatedAt = NowUnix()

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

func (repo *UserRepository) Update(user *model.User) (*model.User, error) {
	user.UpdatedAt = NowUnix()

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

func (repo *UserRepository) FindByEmail(email string) *model.User {
	row := repo.QueryRow(
		"select id, email, password, created_at, updated_at, deleted_at from user where email = ?",
		strings.TrimSpace(strings.ToLower(email)),
	)
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}
