package acl

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	"log"
	"strings"
)

type UserRoleRepository struct {
	server.Repository

	db *sql.DB
	tx *sql.Tx
}

func NewUserRoleRepository(db *sql.DB, tx *sql.Tx) *UserRoleRepository {
	repo := &UserRoleRepository{
		db: db,
		tx: tx,
	}
	return repo
}

func (repo *UserRoleRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserRoleRepository(repo.db, tx)
}

func (repo *UserRoleRepository) Link(user model.User, role *model.Role, resourceId string) error {
	// check if its already been linked
	rows, err := repo.Query(
		"select user_id, role_id, resource_id from user_role where user_id = ? and role_id = ? and resource_id = ?",
		user.Id,
		role.Id,
		resourceId,
	)
	if err != nil {
		log.Print(err)
		return errors.New("failed to link role to user")
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count += 1
	}

	// role is already linked, so do nothing
	if count > 0 {
		return nil
	}

	// otherwise attempt to link the role to the user
	_, err = repo.Exec(
		"insert into user_role (user_id, role_id, resource_id) values (?, ?, ?)",
		user.Id,
		role.Id,
		resourceId,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to link role to user")
	}

	return nil
}

func (repo *UserRoleRepository) Unlink(user model.User, role *model.Role, resourceId string) error {
	_, err := repo.Exec(
		"delete from user_role where user_id = ? and role_id = ? and resource_id = ?",
		user.Id,
		role.Id,
		resourceId,
	)
	if err != nil {
		log.Print(err)
		return errors.New("failed to unlink role from user")
	}

	return nil
}

type ResourceRequest struct {
	ResourcePath []string
	ResourceIds  []string
}

type ResourceResponse struct {
	PermissionId string
	UserId       string
	ResourcePath string
	ResourceId   string
	Action       string
}

/*

In order to find general information about the resources in relation to the user, we need to know the full path to the
resource (e.g. organization:folder:document), we also need to know the parent id for each of those resources. So we would need
to know the id for the document, the id for the folder the document is in, and the id for the organization that the
folder (containing the document) is in.

This means we need the full path, and the correct ids for the hierarchical relationship of the items.

From there we can generate a query that basically checks each level (going towards the target resource) and see
if the correct permission exists for that resource. The reason is that a user may have organization level role that
encompasses access to a document. So the only way to check that, is if we check the root level for that access.

For example, given the following input

resourcePath := []string{"organization", "folder", "document"}
resourceIds := []string{"5", "10", "15"}

The following where clause would be generated

((resource_path = "organization:folder:document" AND resource_id = 5") OR (resource_path = "folder:document" AND resource_id = "10") OR (resource_path = "document" AND resource_id = "10"))

@todo clean this method up / split it up

*/
func (repo *UserRoleRepository) GetDataForResources(userId string, requests []ResourceRequest, action *string) ([]*ResourceResponse, error) {
	params := make([]interface{}, 0)
	clauses := make([]string, 0)

	for _, request := range requests {

		// if these aren't the same length, we have a mismatch and could produce an invalid query
		if len(request.ResourcePath) != len(request.ResourceIds) {
			return nil, errors.New("resource paths and ids must be the same length")
		}

		// build all of the where clauses for the given request
		requestClauses := make([]string, len(request.ResourcePath)*len(request.ResourceIds))
		for idIndex, id := range request.ResourceIds {

			// generate the path which is a slice from the idIndex to the cap
			path := strings.Join(request.ResourcePath[idIndex:cap(request.ResourcePath)], ":")

			clause := "(p.resource_path = ? AND ur.resource_id = ?)"

			params = append(params, path, id)

			requestClauses = append(requestClauses, clause)
		}

		clauses = append(clauses, fmt.Sprintf("(%s)", strings.Join(requestClauses, " OR ")))
	}

	whereClause := fmt.Sprintf("(%s)", strings.Join(clauses, " OR "))

	// add the user id
	params = append(params, userId)
	query := fmt.Sprintf("select p.id, p.resource_path, ur.resource_id, p.action, ur.user_id from user_role ur join role_permission rp on rp.role_id = ur.role_id join permission p on p.id = rp.permission_id where %s AND ur.user_id = ?", whereClause);

	// an action was specified, so also filter on that
	if action != nil {
		params = append(params, *action)
		query += " AND action = ?"
	}

	rows, err := repo.Query(
		query,
		params...
	)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to fetch resource data")
	}
	defer rows.Close()

	results := make([]*ResourceResponse, 0)

	// basically we just need to check that something was returned
	for rows.Next() {

		res := &ResourceResponse{}
		err := rows.Scan(res.PermissionId, res.ResourcePath, res.ResourceId, res.Action, res.UserId)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to fetch resource data")
		}

		results = append(results, res)
	}

	return results, nil
}
