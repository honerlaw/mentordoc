package server

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"log"
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

func (repo *FolderRepository) FindRoots(organizationId string) ([]*model.Folder, error) {
	rows, err := repo.Query(
		"select id, name, parent_folder_id, organization_id, created_at, updated_at, deleted_at from folder where organization_id = ?",
		organizationId,
	)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find root folders")
	}
	defer rows.Close()

	folders := make([]*model.Folder, 0)

	for rows.Next() {
		var folder model.Folder
		err := rows.Scan(&folder.Id, &folder.Name, &folder.ParentFolderId, &folder.OrganizationId, &folder.CreatedAt, &folder.UpdatedAt, &folder.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse roots folder")
		}
		folders = append(folders, &folder)
	}

	return folders, nil
}

func (repo *FolderRepository) FindChildren(folderParentId string, organizationId string) ([]*model.Folder, error) {
	rows, err := repo.Query(
		"select id, name, parent_folder_id, organization_id, created_at, updated_at, deleted_at from folder where parent_folder_id = ? and organization_id = ?",
		folderParentId,
		organizationId,
	)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find child folders")
	}
	defer rows.Close()

	folders := make([]*model.Folder, 0)

	for rows.Next() {
		var folder model.Folder
		err := rows.Scan(&folder.Id, &folder.Name, &folder.ParentFolderId, &folder.CreatedAt, &folder.UpdatedAt, &folder.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse child folder")
		}
		folders = append(folders, &folder)
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
