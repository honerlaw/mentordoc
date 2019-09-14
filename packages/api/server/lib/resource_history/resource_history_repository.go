package resource_history

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
)

type ResourceHistoryRepository struct {
	util.Repository
}

func NewResourceHistoryRepository(db *sql.DB, tx *sql.Tx) *ResourceHistoryRepository {
	repo := &ResourceHistoryRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *ResourceHistoryRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewResourceHistoryRepository(repo.Db, tx)
}

func (repo *ResourceHistoryRepository) Insert(history *shared.ResourceHistory) (*shared.ResourceHistory, error) {
	history.CreatedAt = util.NowUnix()
	history.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into resource_history (id, resource_id, resource_name, user_id, action, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?, ?)",
		history.Id,
		history.ResourceId,
		history.ResourceName,
		history.UserId,
		history.Action,
		history.CreatedAt,
		history.UpdatedAt,
		history.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert resource history")
	}

	return history, nil;
}

func (repo *ResourceHistoryRepository) FindOne(resourceId string, resourceName string, userId string, action string) *shared.ResourceHistory {
	row := repo.QueryRow(
		"select id, resource_id, resource_name, user_id, action, created_at, updated_at, deleted_at from resource_history where resource_id = ? and resource_name = ? and user_id = ? and action = ? and deleted_at is null",
		resourceId,
		resourceName,
		userId,
		action,
	)

	var history shared.ResourceHistory
	err := row.Scan(&history.Id, &history.ResourceId, &history.ResourceName, &history.UserId, &history.Action, &history.CreatedAt, &history.UpdatedAt, &history.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &history;
}