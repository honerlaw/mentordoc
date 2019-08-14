package server

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Insert(user *User) (*User, error) {
	user.CreatedAt = NowUnix()
	user.UpdatedAt = NowUnix()

	_, err := repo.db.Exec(
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

func (repo *UserRepository) Update(user *User) (*User, error) {
	user.UpdatedAt = NowUnix()

	_, err := repo.db.Exec(
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

func (repo *UserRepository) FindByEmail(email string) *User {
	row := repo.db.QueryRow(
		"select id, email, password, created_at, updated_at, deleted_at from user where email = ?",
		strings.TrimSpace(strings.ToLower(email)),
	)
	user := &User{}
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}
