package acl

import (
	"errors"
	"github.com/honerlaw/mentordoc/server/model"
	"reflect"
)

type modelAclData struct {
	Path         []string
	StructFields []string
}

type AclWrapperService struct {
	aclService *AclService
	data       map[string]*modelAclData
}

type AclWrappedModel struct {
	Model   interface{} `json:"model"`
	Actions []string    `json:"actions"`
}

func NewAclWrapperService(aclService *AclService) *AclWrapperService {
	return &AclWrapperService{
		aclService: aclService,

		// @todo move this data to be on the model itself, then we can just use reflection to get it all
		data: map[string]*modelAclData{
			"*model.Organization": {
				Path:         []string{"organization"},
				StructFields: []string{"Id"},
			},
			"*model.Folder": {
				Path:         []string{"organization", "folder"},
				StructFields: []string{"OrganizationId", "Id"},
			},
			"*model.Document": {
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
func (service *AclWrapperService) Wrap(user *model.User, modelSlice interface{}) ([]*AclWrappedModel, error) {
	s := reflect.ValueOf(modelSlice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("a slice of models is required")
	}

	// convert the passed models to an interface array so we can work with it...
	models := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		models[i] = s.Index(0).Interface()
	}

	paths := make([][]string, len(models))
	ids := make([][]string, len(models))

	for index, m := range models {
		modelName := reflect.TypeOf(m).String()
		modelData := service.data[modelName]

		// find the ids needed for the model and paths
		modelIds, err := service.getIdsForModel(m, modelData)
		if err != nil {
			return nil, err
		}
		paths[index] = modelData.Path
		ids[index] = modelIds
	}

	// find all of the actions that can be done on the given model
	resp, err := service.aclService.UserActionsForResources(user, paths, ids)
	if err != nil {
		return nil, err
	}

	// prepopulate the model map with all of the passed models
	modelMap := make(map[string]*AclWrappedModel)
	for _, m := range models {
		id := reflect.ValueOf(m).Elem().FieldByName("Id").Interface().(string)
		modelMap[id] = &AclWrappedModel{
			Model:   m,
			Actions: make([]string, 0),
		}
	}

	// merge the actions together with the model
	for _, res := range resp {
		wrapper := modelMap[res.ResourceId]
		wrapper.Actions = append(wrapper.Actions, res.Action)
	}

	// convert the map to an array
	modelArray := make([]*AclWrappedModel, 0)
	for _, wrapper := range modelMap {
		modelArray = append(modelArray, wrapper)
	}

	return modelArray, nil

}

/**
Grab the values of the id fields off the model to be used with the acl service
*/
func (service *AclWrapperService) getIdsForModel(model interface{}, data *modelAclData) ([]string, error) {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return nil, errors.New("model must be pointer to struct")
	}

	derefValue := value.Elem()

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
