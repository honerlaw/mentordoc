package server

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"log"
)

type OrganizationRepository struct {
	util.Repository
}

func NewOrganizationRepository(db *sql.DB, tx *sql.Tx) *OrganizationRepository {
	repo := &OrganizationRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *OrganizationRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewOrganizationRepository(repo.Db, tx)
}

func (repo *OrganizationRepository) FindById(id string) *model.Organization {
	row := repo.QueryRow(
		"select id, name, created_at, updated_at, deleted_at from organization where id = ?",
		id,
	)

	var organization model.Organization
	err := row.Scan(&organization.Id, &organization.Name, &organization.CreatedAt, &organization.UpdatedAt, &organization.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &organization;
}

func (repo *OrganizationRepository) Insert(org *model.Organization) (*model.Organization, error) {
	org.CreatedAt = util.NowUnix()
	org.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
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

func (repo *OrganizationRepository) Update(org *model.Organization) (*model.Organization, error) {
	org.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
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
