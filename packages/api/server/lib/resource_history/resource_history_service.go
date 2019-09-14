package resource_history

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	uuid "github.com/satori/go.uuid"
)

type ResourceHistoryService struct {
	resourceHistoryRepository *ResourceHistoryRepository
}

func NewResourceHistoryService(
	resourceHistoryRepository *ResourceHistoryRepository,
) *ResourceHistoryService {
	service := &ResourceHistoryService{
		resourceHistoryRepository: resourceHistoryRepository,
	};
	return service
}

func (service *ResourceHistoryService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewResourceHistoryService(
		service.resourceHistoryRepository.InjectTransaction(tx).(*ResourceHistoryRepository),
	)
}

func (service *ResourceHistoryService) Create(resourceId string, resourceName string, userId string, action string) (*shared.ResourceHistory, error) {
	history := &shared.ResourceHistory{
		ResourceId:   resourceId,
		ResourceName: resourceName,
		UserId:       userId,
		Action:       action,
	}
	history.Id = uuid.NewV4().String()

	history, err := service.resourceHistoryRepository.Insert(history)
	if err != nil {
		return nil, shared.NewNotFoundError("could not create resource history")
	}

	return history, nil;
}
