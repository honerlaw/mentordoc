package server

import (
	"database/sql"
	"errors"
	"log"
)

type OrganizationRepository struct {
	Transactionable
	db *sql.DB
	tx *sql.Tx
}

func NewOrganizationRepository(db *sql.DB, tx *sql.Tx) *OrganizationRepository {
	repo := &OrganizationRepository{
		db: db,
		tx: tx,
	}
	repo.cloneWithTransaction = func(tx *sql.Tx) interface{} {
		return NewOrganizationRepository(repo.db, tx)
	}
	return repo
}

func (repo *OrganizationRepository) Insert(org *Organization) (*Organization, error) {
	org.CreatedAt = NowUnix()
	org.UpdatedAt = NowUnix()

	_, err := repo.db.Exec(
		"insert into organization (id, name, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?)",
		org.Id,
		org.Name,
		org.CreatedAt,
		org.UpdatedAt,
		org.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert organization")
	}

	return org, nil;
}

func (repo *OrganizationRepository) Update(org *Organization) (*Organization, error) {
	org.UpdatedAt = NowUnix()

	_, err := repo.db.Exec(
		"update organization set name = ?, updated_at = ?, deleted_at = ? where id = ?",
		org.Name,
		org.UpdatedAt,
		org.DeletedAt,
		org.Id,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to update organization")
	}

	return org, nil;
}
