package acl

import (
	"errors"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"reflect"
	"strings"
)

type modelAclData struct {
	Path         []string
	StructFields []string
}

type AclWrapperService struct {
	aclService *AclService
	data       map[string]*modelAclData
}

func NewAclWrapperService(aclService *AclService) *AclWrapperService {
	return &AclWrapperService{
		aclService: aclService,

		// @todo move this data to be on the model itself, then we can just use reflection to get it all
		data: map[string]*modelAclData{
			"Organization": {
				Path:         []string{"organization"},
				StructFields: []string{"Id"},
			},
			"Folder": {
				Path:         []string{"organization", "folder"},
				StructFields: []string{"OrganizationId", "Id"},
			},
			"Document": {
				Path:         []string{"organization", "folder", "document"},
				StructFields: []string{"OrganizationId", "FolderId", "Id"},
			},
		},
	}
}

/*
wraps the given set of organization / folder / document models with the actions that the given user is capable of doing
to each resource. This way the client will know what info to show, we only need to fetch this data as it is going
out back to the client
 */
func (service *AclWrapperService) Wrap(user *shared.User, modelSlice interface{}) ([]AclWrappedModel, error) {
	s := reflect.ValueOf(modelSlice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("a slice of models is required")
	}

	// convert the passed models to an interface array so we can work with it...
	models := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		models[i] = s.Index(i).Interface()
	}

	paths := make([][]string, len(models))
	ids := make([][]string, len(models))

	for index, m := range models {
		resourceData, err := service.GetResourceDataForModel(m);
		if err != nil {
			return nil, err
		}
		paths[index] = resourceData.ResourcePath
		ids[index] = resourceData.ResourceIds
	}


	// find all of the actions that can be done on the given model
	resp, err := service.aclService.UserActionsForResources(user, paths, ids)
	if err != nil {
		return nil, err
	}

	// wrap the model with the acl data
	wrappedSlice := make([]AclWrappedModel, 0)
	for index, m := range models {

		wrapper := &AclWrappedModel{
			Model:   m,
			Actions: make([]string, 0),
		}

		// go over the responses and see if any belong to this modedl
		for _, res := range resp {
			for _, id := range ids[index] {
				if id == res.ResourceId {
					wrapper.Actions = append(wrapper.Actions, res.Action)
				}
			}
		}

		wrappedSlice = append(wrappedSlice, *wrapper)
	}

	return wrappedSlice, nil
}

func (service *AclWrapperService) GetResourceDataForModel(m interface{}) (*ResourceData, error) {
	modelName := reflect.TypeOf(m).String()
	modelData, ok := service.data[strings.Split(strings.ReplaceAll(modelName, "*", ""), ".")[1]]

	if !ok {
		return nil, errors.New("could not find data for model")
	}

	modelIds, err := service.getIdsForModel(m, modelData)
	if err != nil {
		return nil, err
	}
	return &ResourceData{
		ResourceIds: modelIds,
		ResourcePath: modelData.Path,
	}, nil
}

/**
Grab the values of the id fields off the model to be used with the acl service
*/
func (service *AclWrapperService) getIdsForModel(model interface{}, data *modelAclData) ([]string, error) {
	value := reflect.ValueOf(model)

	derefValue := value
	if value.Kind() == reflect.Ptr {
		derefValue = value.Elem()
	}

	if derefValue.Kind() != reflect.Struct {
		return nil, errors.New("model must be a struct")
	}

	ids := make([]string, len(data.StructFields))

	for index, field := range data.StructFields {
		fieldValue := derefValue.FieldByName(field)

		var value string
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				value = ""
			} else {
				value = fieldValue.Elem().Interface().(string)
			}
		} else {
			value = fieldValue.Interface().(string)
		}

		ids[index] = value
	}

	return ids, nil
}
