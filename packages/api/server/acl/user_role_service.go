package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/pkg/errors"
)

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

func (service *UserRoleService) LinkUserToRole(user *model.User, roleName string, resourceId string) error {
	role := service.roleRepository.Find(roleName)
	if role == nil {
		return errors.New("failed to find role")
	}

	return service.userRoleRepository.Link(user, role, resourceId)
}

/*
Check if a user can access a specific resource for the given action
 */
func (service *UserRoleService) UserCanAccessResource(user *model.User, path []string, ids []string, action string) (bool, error) {
	request := []ResourceRequest{
		{
			ResourcePath: path,
			ResourceIds:  ids,
			Action: &action,
		},
	}

	data, err := service.userRoleRepository.GetDataForResources(user, request)

	if err != nil {
		return false, err
	}

	return len(data) > 0, nil
}

/**
Fetch the actions that the user an do on each resource, this data needs to be merged with the actual models elsewhere
 */
func (service *UserRoleService) UserActionsForResources(user *model.User, paths [][]string, ids [][]string) ([]*ResourceResponse, error) {
	if len(paths) != len(ids) {
		return nil, errors.New("path and ids must be the same length")
	}

	// build all of the requests
	requests := make([]ResourceRequest, len(paths))
	for index, path := range paths {
		requests[index] = ResourceRequest{
			ResourcePath: path,
			ResourceIds: ids[index],
		}
	}

	data, err := service.userRoleRepository.GetDataForResources(user, requests)

	if err != nil {
		return nil, err
	}

	return data, nil
}

/**
Find the resource data for the specified resource path and action. E.g. to find all viewable documents
 */
func (service *UserRoleService) UserActionableResourcesByPath(user *model.User, path []string, action string) ([]*ResourceResponse, error) {
	request := []ResourceRequest{
		{
			ResourcePath: path,
			Action: &action,
		},
	}

	data, err := service.userRoleRepository.GetDataForResources(user, request)

	if err != nil {
		return nil, err
	}

	return data, nil
}
