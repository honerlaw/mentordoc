package service

import (
	"database/sql"
	"errors"
	"log"
	"server"
	"server/model"
	"strings"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(user *model.User) (*model.User, error) {
	user.CreatedAt = server.NowUnix()
	user.UpdatedAt = server.NowUnix()

	_, err := dao.db.Exec(
		"insert into user (id, email, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?)",
		user.Id,
		strings.TrimSpace(strings.ToLower(user.Email)),
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

func (dao *UserDao) Update(user *model.User) (*model.User, error) {
	user.UpdatedAt = server.NowUnix()

	_, err := dao.db.Exec(
		"update user set email = ?, updated_at = ?, deleted_at = ? where id = ?",
		strings.TrimSpace(strings.ToLower(user.Email)),
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

func (dao *UserDao) FindByEmail(email string) *model.User {
	row := dao.db.QueryRow(
		"select id, email, created_at, updated_at, deleted_at from user where email = ?",
		strings.TrimSpace(strings.ToLower(email)),
	)
	user := &model.User{}
	err := row.Scan(user.Id, user.Email, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
	if err != nil {
		return nil
	}
	return user
}
