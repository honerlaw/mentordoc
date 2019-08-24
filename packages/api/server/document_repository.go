package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"log"
	"strings"
)

type DocumentRepository struct {
	util.Repository
}

func NewDocumentRepository(db *sql.DB, tx *sql.Tx) *DocumentRepository {
	repo := &DocumentRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *DocumentRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewDocumentRepository(repo.Db, tx)
}

func (repo *DocumentRepository) FindById(id string) *model.Document {
	row := repo.QueryRow(
		"select id, name, folder_id, organization_id, created_at, updated_at, deleted_at from document where id = ?",
		id,
	)

	var document model.Document
	err := row.Scan(&document.Id, &document.Name, &document.FolderId, &document.OrganizationId, &document.CreatedAt, &document.UpdatedAt, &document.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &document
}

func (repo *DocumentRepository) Find(organizationIds []string, folderIds []string, documentIds []string, pagination *model.Pagination) ([]model.Document, error) {
	query := "select id, name, folder_id, organization_id, created_at, updated_at, deleted_at from document where"

	params := make([]interface{}, 0)

	// build the in queries
	inQueries := make([]string, 0)
	if len(organizationIds) > 0 {
		inQueries = append(inQueries, fmt.Sprintf("organization_id in (%s)", util.BuildSqlPlaceholderArray(organizationIds)))
		params = append(params, util.ConvertStringArrayToInterfaceArray(organizationIds)...)
	}
	if len(folderIds) > 0 {
		inQueries = append(inQueries, fmt.Sprintf("folder_id in (%s)", util.BuildSqlPlaceholderArray(folderIds)))
		params = append(params, util.ConvertStringArrayToInterfaceArray(folderIds)...)
	}
	if len(documentIds) > 0 {
		inQueries = append(inQueries, fmt.Sprintf("id in (%s)", util.BuildSqlPlaceholderArray(documentIds)))
		params = append(params, util.ConvertStringArrayToInterfaceArray(documentIds)...)
	}

	// tack on the in query
	query = fmt.Sprintf("%s (%s)", query, strings.Join(inQueries, " OR "))

	// add the pagination portion of the query
	if pagination != nil {
		query = fmt.Sprintf("%s LIMIT ?, ?", query)
		params = append(params, pagination.Page * pagination.Count, pagination.Count)
	}

	rows, err := repo.Query(query, params...)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find documents")
	}
	defer rows.Close()

	documents := make([]model.Document, 0)
	for rows.Next() {
		var document model.Document
		err := rows.Scan(&document.Id, &document.Name, &document.FolderId, &document.OrganizationId, &document.CreatedAt, &document.UpdatedAt, &document.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse document")
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func (repo *DocumentRepository) Insert(document *model.Document) error {
	document.CreatedAt = util.NowUnix()
	document.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into document (id, name, folder_id, organization_id, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?)",
		document.Id,
		document.Name,
		document.FolderId,
		document.OrganizationId,
		document.CreatedAt,
		document.UpdatedAt,
		document.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to insert document")
	}

	return nil;
}

func (repo *DocumentRepository) Update(document *model.Document) error {
	document.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"update document set name = ?, folder_id = ?, updated_at = ?, deleted_at = ? where id = ?",
		document.Name,
		document.FolderId,
		document.UpdatedAt,
		document.DeletedAt,
		document.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to update document")
	}

	return nil;
}