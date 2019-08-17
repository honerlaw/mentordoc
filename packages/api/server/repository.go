package server

import "database/sql"

type Repository struct {
	db *sql.DB
	tx *sql.Tx
}

func (repo *Repository) Exec(query string, args ...interface{}) (sql.Result, error) {
	if repo.tx != nil {
		return repo.tx.Exec(query, args)
	}
	return repo.db.Exec(query, args)
}

func (repo *Repository) QueryRow(query string, args ...interface{}) *sql.Row {
	if repo.tx != nil {
		return repo.tx.QueryRow(query, args)
	}
	return repo.db.QueryRow(query, args)
}
