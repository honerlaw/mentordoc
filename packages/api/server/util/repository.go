package util

import "database/sql"


type Repository struct {
	Db *sql.DB
	Tx *sql.Tx
}

func (repo *Repository) Exec(query string, args ...interface{}) (sql.Result, error) {
	if repo.Tx != nil {
		return repo.Tx.Exec(query, args...)
	}
	return repo.Db.Exec(query, args...)
}

func (repo *Repository) QueryRow(query string, args ...interface{}) *sql.Row {
	if repo.Tx != nil {
		return repo.Tx.QueryRow(query, args...)
	}
	return repo.Db.QueryRow(query, args...)
}

func (repo *Repository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if repo.Tx != nil {
		return repo.Tx.Query(query, args...)
	}
	return repo.Db.Query(query, args...)
}
