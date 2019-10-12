package document

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
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

func (repo *DocumentRepository) FindById(id string) *shared.Document {
	row := repo.QueryRow(
		"select id, folder_id, organization_id, created_at, updated_at, deleted_at from document where id = ? and deleted_at is null",
		id,
	)

	var document shared.Document
	err := row.Scan(&document.Id, &document.FolderId, &document.OrganizationId, &document.CreatedAt, &document.UpdatedAt, &document.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &document
}

func (repo *DocumentRepository) Find(userId string, organizationIds []string, folderIds []string, documentIds []string, folderId *string, pagination *shared.Pagination) ([]shared.Document, error) {
	query := "select distinct d.id, d.folder_id, d.organization_id, d.created_at, d.updated_at, d.deleted_at from document d WHERE "

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

	// add the specific check for a specific folder
	if folderId != nil {
		query = fmt.Sprintf("%s AND d.folder_id = ?", query)
		params = append(params, *folderId)
	} else {
		query = fmt.Sprintf("%s AND d.folder_id is null", query)
	}

	query = fmt.Sprintf("%s AND d.deleted_at is null ORDER BY d.created_at ASC", query)

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

	documents := make([]shared.Document, 0)
	for rows.Next() {
		var document shared.Document
		err := rows.Scan(&document.Id, &document.FolderId, &document.OrganizationId, &document.CreatedAt, &document.UpdatedAt, &document.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse document")
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func (repo *DocumentRepository) Insert(document *shared.Document) error {
	document.CreatedAt = util.NowUnix()
	document.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into document (id, folder_id, organization_id, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?)",
		document.Id,
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

func (repo *DocumentRepository) Update(document *shared.Document) error {
	document.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"update document set folder_id = ?, updated_at = ?, deleted_at = ? where id = ?",
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