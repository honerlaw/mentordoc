package organization

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
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

func (repo *OrganizationRepository) FindById(id string) *shared.Organization {
	row := repo.QueryRow(
		"select id, name, created_at, updated_at, deleted_at from organization where id = ?",
		id,
	)

	var organization shared.Organization
	err := row.Scan(&organization.Id, &organization.Name, &organization.CreatedAt, &organization.UpdatedAt, &organization.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &organization;
}

func (repo *OrganizationRepository) Insert(org *shared.Organization) (*shared.Organization, error) {
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

func (repo *OrganizationRepository) Update(org *shared.Organization) (*shared.Organization, error) {
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

func (repo *OrganizationRepository) Find(organizationIds []string) ([]shared.Organization, error) {
	if len(organizationIds) == 0 {
		return make([]shared.Organization, 0), nil
	}

	// build the in query
	inQuery := fmt.Sprintf("id in (%s) ORDER BY name ASC", util.BuildSqlPlaceholderArray(organizationIds))
	params := util.ConvertStringArrayToInterfaceArray(organizationIds)
	query := fmt.Sprintf("select distinct id, name, created_at, updated_at, deleted_at from organization where %s", inQuery)

	rows, err := repo.Query(query, params...)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find organizations")
	}
	defer rows.Close()

	orgs := make([]shared.Organization, 0)
	for rows.Next() {
		var org shared.Organization
		err := rows.Scan(&org.Id, &org.Name, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse organization")
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}