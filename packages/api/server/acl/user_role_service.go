package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/pkg/errors"
)
/**
TODO
- We need a way to fetch all of the actions for a given resource
-- this will probably be very similar to checking for access, just excluding the action scope and expanding to allow
-- an array of an array of paths and ids (e.g. [][]string for both paths and ids), might be better to convert this to a struct
-- of some sort
 */

type UserRoleService struct {
	roleRepository     *RoleRepository
	userRoleRepository *UserRoleRepository
}

func NewUserRoleService(roleRepository *RoleRepository, userRoleRepository *UserRoleRepository) *UserRoleService {
	return &UserRoleService{
		roleRepository:     roleRepository,
		userRoleRepository: userRoleRepository,
	}
}

func (service *UserRoleService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserRoleService(service.roleRepository.InjectTransaction(tx).(*RoleRepository),
		service.userRoleRepository.InjectTransaction(tx).(*UserRoleRepository))
}

func (service *UserRoleService) LinkUserToRole(user model.User, roleName string, resourceId string) error {
	role := service.roleRepository.Find(roleName)
	if role != nil {
		return errors.New("failed to find role")
	}

	return service.userRoleRepository.Link(user, role, resourceId)
}

func (service *UserRoleService) UserCanAccessResource(user model.User, path []string, ids []string, action string) (bool, error) {
	return service.userRoleRepository.CanAccessResource(user.Id, path, ids, action)
}

func (service *UserRoleService) UserActionsForResource() (bool, error) {

}
