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

type FolderRepository struct {
	util.Repository
}

func NewFolderRepository(db *sql.DB, tx *sql.Tx) *FolderRepository {
	repo := &FolderRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *FolderRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewFolderRepository(repo.Db, tx)
}

func (repo *FolderRepository) Insert(folder *model.Folder) error {
	folder.CreatedAt = util.NowUnix()
	folder.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into folder (id, name, parent_folder_id, organization_id, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?)",
		folder.Id,
		folder.Name,
		folder.ParentFolderId,
		folder.OrganizationId,
		folder.CreatedAt,
		folder.UpdatedAt,
		folder.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to insert folder")
	}

	return nil;
}

func (repo *FolderRepository) Update(folder *model.Folder) error {
	folder.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"update folder set name = ?, parent_folder_id = ?, updated_at = ?, deleted_at = ? where id = ?",
		folder.Name,
		folder.ParentFolderId,
		folder.UpdatedAt,
		folder.DeletedAt,
		folder.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to update folder")
	}

	return nil;
}

func (repo *FolderRepository) Find(organizationIds []string, folderIds []string, parentFolderId *string, pagination *model.Pagination) ([]model.Folder, error) {
	query := "select id, name, parent_folder_id, organization_id, created_at, updated_at, deleted_at from folder where"

	params := make([]interface{}, 0)

	// build the in queries
	inQueries := make([]string, 0)
	if len(organizationIds) > 0 {
		inQueries = append(inQueries, fmt.Sprintf("organization_id in (%s)", util.BuildSqlPlaceholderArray(organizationIds)))
		params = append(params, util.ConvertStringArrayToInterfaceArray(organizationIds)...)
	}
	if len(folderIds) > 0 {
		inQueries = append(inQueries, fmt.Sprintf("id in (%s)", util.BuildSqlPlaceholderArray(folderIds)))
		params = append(params, util.ConvertStringArrayToInterfaceArray(folderIds)...)
	}

	// tack on the in query
	query = fmt.Sprintf("%s (%s)", query, strings.Join(inQueries, " OR "))

	// add the parent folder portion of the where clause
	if parentFolderId != nil {
		query = fmt.Sprintf("%s AND parent_folder_id = ?", query)
		params = append(params, *parentFolderId)
	} else {
		query = fmt.Sprintf("%s AND parent_folder_id is null", query)
	}

	query = fmt.Sprintf("%s ORDER BY name ASC", query)

	// add the pagination portion of the query
	if pagination != nil {
		query = fmt.Sprintf("%s LIMIT ?, ?", query)
		params = append(params, pagination.Page * pagination.Count, pagination.Count)
	}

	rows, err := repo.Query(query, params...)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find folders")
	}
	defer rows.Close()

	folders := make([]model.Folder, 0)
	for rows.Next() {
		var folder model.Folder
		err := rows.Scan(&folder.Id, &folder.Name, &folder.ParentFolderId, &folder.OrganizationId, &folder.CreatedAt, &folder.UpdatedAt, &folder.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse folder")
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (repo *FolderRepository) FindById(id string) *model.Folder {
	row := repo.QueryRow(
		"select id, name, parent_folder_id, organization_id, created_at, updated_at, deleted_at from folder where id = ?",
		id,
	)

	var folder model.Folder
	err := row.Scan(&folder.Id, &folder.Name, &folder.ParentFolderId, &folder.OrganizationId, &folder.CreatedAt, &folder.UpdatedAt, &folder.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &folder
}
