package folder

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
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

func (repo *FolderRepository) Insert(folder *shared.Folder) error {
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

func (repo *FolderRepository) Update(folder *shared.Folder) error {
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

// @todo list should probably do this in a single query
func (repo *FolderRepository) ChildCount(folders []shared.Folder) error {
	if len(folders) == 0 {
		return nil
	}

	inQuery := fmt.Sprintf("parent_folder_id in (%s)", util.BuildSqlPlaceholderArray(folders))
	params := make([]interface{}, 0)
	for _, folder := range folders {
		params = append(params, folder.Id)
	}

	inQueryDocument := fmt.Sprintf("folder_id in (%s)", util.BuildSqlPlaceholderArray(folders))
	for _, folder := range folders {
		params = append(params, folder.Id)
	}

	// add for the second query
	for _, folder := range folders {
		params = append(params, folder.Id)
	}
	for _, folder := range folders {
		params = append(params, folder.Id)
	}

	query := fmt.Sprintf("select f.parent_folder_id, f.folderCount, d.folder_id, d.documentCount from (select parent_folder_id, count(id) as folderCount from folder where %s and deleted_at is null group by parent_folder_id) f left join (select folder_id, count(id) as documentCount from document where %s and deleted_at is null group by folder_id) d on f.parent_folder_id = d.folder_id UNION select f.parent_folder_id, f.folderCount, d.folder_id, d.documentCount from (select parent_folder_id, count(id) as folderCount from folder where %s and deleted_at is null group by parent_folder_id) f right join (select folder_id, count(id) as documentCount from document where %s and deleted_at is null group by folder_id) d on f.parent_folder_id = d.folder_id", inQuery, inQueryDocument, inQuery, inQueryDocument)

	rows, err := repo.Query(query, params...)
	if err != nil {
		log.Print(err)
		return errors.New("failed to find folder counts")
	}
	defer rows.Close()

	type FolderCount struct {
		ParentFolderId *string
		FolderCount    *int
		FolderId       *string
		DocumentCount  *int
	}

	for rows.Next() {
		var folderCount FolderCount
		err := rows.Scan(&folderCount.ParentFolderId, &folderCount.FolderCount, &folderCount.FolderId, &folderCount.DocumentCount)
		if err != nil {
			log.Print(err)
			return err
		}

		for i := 0; i < len(folders); i++ {
			folder := &folders[i];
			if folderCount.ParentFolderId != nil && *folderCount.ParentFolderId == folder.Id {
				folder.ChildCount += *folderCount.FolderCount;
			}
			if folderCount.FolderId != nil && *folderCount.FolderId == folder.Id {
				folder.ChildCount += *folderCount.DocumentCount;
			}
		}
	}

	return nil
}

func (repo *FolderRepository) Find(organizationIds []string, folderIds []string, parentFolderId *string, pagination *shared.Pagination) ([]shared.Folder, error) {
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

	query = fmt.Sprintf("%s AND deleted_at is null ORDER BY name ASC", query)

	// add the pagination portion of the query
	if pagination != nil {
		query = fmt.Sprintf("%s LIMIT ?, ?", query)
		params = append(params, pagination.Page*pagination.Count, pagination.Count)
	}

	rows, err := repo.Query(query, params...)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find folders")
	}
	defer rows.Close()

	folders := make([]shared.Folder, 0)
	for rows.Next() {
		var folder shared.Folder
		err := rows.Scan(&folder.Id, &folder.Name, &folder.ParentFolderId, &folder.OrganizationId, &folder.CreatedAt, &folder.UpdatedAt, &folder.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse folder")
		}
		folders = append(folders, folder)
	}

	err = repo.ChildCount(folders)
	if err != nil {
		return nil, errors.New("failed to get child count for folders")
	}

	return folders, nil
}

func (repo *FolderRepository) FindById(id string) *shared.Folder {
	row := repo.QueryRow(
		"select id, name, parent_folder_id, organization_id, created_at, updated_at, deleted_at from folder where id = ? and deleted_at is null",
		id,
	)

	var folder shared.Folder
	err := row.Scan(&folder.Id, &folder.Name, &folder.ParentFolderId, &folder.OrganizationId, &folder.CreatedAt, &folder.UpdatedAt, &folder.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	err = repo.ChildCount([]shared.Folder{folder})
	if err != nil {
		log.Print(err)
		return nil
	}

	return &folder
}
