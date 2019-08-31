package acl

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
	"strings"
)

type UserRoleRepository struct {
	util.Repository
}

func NewUserRoleRepository(db *sql.DB, tx *sql.Tx) *UserRoleRepository {
	repo := &UserRoleRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *UserRoleRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserRoleRepository(repo.Db, tx)
}

func (repo *UserRoleRepository) Link(user *shared.User, role *Role, resourceId string) error {
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

func (repo *UserRoleRepository) Unlink(user *shared.User, role *Role, resourceId string) error {
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

/*
There are two types of querys that need to be supported

1. We know the resource (path add ids) and user, but do we have access to it
2. We know the resource path and user and action, but want to know the resource ids

This method builds the where clause depending on what we are trying to fetch.

1. We can specify a path
2. We can specify a path and ids
3. We can specify a path and an action
4. We can specify a path, ids, and an action

The where clause that is built is fairly straight forward.

We know that every resource is denoted by a full path. E.g. for access to a document, its either denoted as
organization:folder:document, folder:document, or document. So this part of the where clause will always be built out.

If a resource ids are also specified, they will be added to the query, if the action is specified, it is added to the
query.

Examples based on what is specified:

1. a document path ([]string{"organization", "folder", "document"})

((resource_path = "organization:folder:document") OR (resource_path = "folder:document") OR (resource_path = "folder:document")) and userId = "12345"

2. a document path ([]string{"organization", "folder", "document"}) and ids ([]string{"5", "10", "15"})

((resource_path = "organization:folder:document" AND resource_id = "5") OR (resource_path = "folder:document" AND resource_id = "10") OR (resource_path = "folder:document" AND resource_id = "15")) and userId = "12345"

3. a document path ([]string{"organization", "folder", "document"}) and action (string{"view"})

((resource_path = "organization:folder:document" AND action = "view") OR (resource_path = "folder:document" AND action = "view") OR (resource_path = "folder:document" AND action = "view")) and userId = "12345"

4. a document path ([]string{"organization", "folder", "document"}) and ids ([]string{"5", "10", "15"}) and action (string{"view"})

((resource_path = "organization:folder:document" AND resource_id = "5" AND action = "view") OR (resource_path = "folder:document" AND resource_id = "10" AND action = "view") OR (resource_path = "folder:document" AND resource_id = "15" AND action = "view")) and userId = "12345"

Given these various different ways to filter based on a single user, we can then use the resulting data to query for
the actual data. This allows ACL to be decoupled from the application logic, but at the cost of needing to run a
minimum of 2 queries (one for the acl data, one for the application data), we could in the future have a single query,
if there is a need to do so
*/
func (repo *UserRoleRepository) buildWhereClause(userId string, requests []ResourceRequest) (*string, []interface{}, error) {
	params := make([]interface{}, 0)
	clauses := make([]string, 0)
	for _, request := range requests {

		// if these aren't the same length, we have a mismatch and could produce an invalid query
		if request.ResourceIds != nil && len(request.ResourcePath) != len(request.ResourceIds) {
			return nil, nil, errors.New("resource paths and ids must be the same length")
		}

		// build all of the where clauses for the given request
		requestClauses := make([]string, 0)
		for pathIndex := 0; pathIndex < len(request.ResourcePath); pathIndex++ {

			// generate the path which is a slice from the idIndex to the cap
			path := strings.Join(request.ResourcePath[pathIndex:cap(request.ResourcePath)], ":")

			clause := "p.resource_path = ?"
			params = append(params, path)

			if request.ResourceIds != nil {
				id := request.ResourceIds[pathIndex];
				clause += " AND ur.resource_id = ?"
				params = append(params, id)
			}

			if request.Action != nil {
				clause += " AND p.action = ?"
				params = append(params, *request.Action)
			}

			requestClauses = append(requestClauses, fmt.Sprintf("(%s)", clause))
		}

		clauses = append(clauses, fmt.Sprintf("(%s)", strings.Join(requestClauses, " OR ")))
	}

	clause := fmt.Sprintf("(%s)", strings.Join(clauses, " OR "))

	return &clause, params, nil
}

func (repo *UserRoleRepository) GetDataForResources(user *shared.User, requests []ResourceRequest) ([]ResourceResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("must supply at least one resource request")
	}

	whereClause, params, err := repo.buildWhereClause(user.Id, requests)
	if err != nil {
		return nil, err
	}

	// add the user id
	params = append(params, user.Id)
	query := fmt.Sprintf("select distinct p.id, p.resource_path, ur.resource_id, p.action, ur.user_id from user_role ur join role_permission rp on rp.role_id = ur.role_id join permission p on p.id = rp.permission_id where %s AND ur.user_id = ?", *whereClause);

	rows, err := repo.Query(
		query,
		params...
	)
	if err != nil {
		return nil, errors.New("failed to fetch resource data")
	}
	defer rows.Close()

	results := make([]ResourceResponse, 0)

	// basically we just need to check that something was returned
	for rows.Next() {

		var res ResourceResponse
		err := rows.Scan(&res.PermissionId, &res.ResourcePath, &res.ResourceId, &res.Action, &res.UserId)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to fetch resource data")
		}

		results = append(results, res)
	}

	return results, nil
}
