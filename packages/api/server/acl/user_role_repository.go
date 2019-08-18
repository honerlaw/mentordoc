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

/*

In order to find if a person has access to a resource, we need to know the full path to the resource
e.g. organization:folder:document, we also need to know the parent id for each of those resources. So we would need
to know the id for the document, the id for the folder the document is in, and the id for the organization that the
folder (containing the document) is in.

This means we need the full path, and the correct ids for the hierarchical relationship of the items.

From there we can generate a query that basically checks each level (going towards the target resource) and see
if the correct permission exists for that resource. The reason is that a user may have organization level role that
encompasses access to a document. So the only way to check that, is if we check the root level for that access.

For example, given the following input

resourcePath := []string{"organization", "folder"}
resourceIds := []string{"5", "10"}

The following where clause would be generated

((resource_path = "organization" AND resource_id = 5") OR (resource_path = "organization:folder" AND resource_id = 5"))
OR ((resource_path = "organization:folder" AND resource_id = 10"))

*/
func (repo *UserRoleRepository) CanAccessResource(userId string, resourcePath []string, resourceIds []string, action string) (bool, error) {

	// if these aren't the same length, we have a mismatch and could produce an invalid query
	if len(resourcePath) != len(resourceIds) {
		return false, errors.New("resource paths and ids must be the same length")
	}

	// start building the wonderful where clause
	params := make([]interface{}, 2 * len(resourcePath) * len(resourceIds) + 1)
	clauses := make([]string, len(resourcePath) * len(resourceIds))
	for idIndex, id := range resourceIds {

		subClauses := make([]string, len(resourcePath))

		// so we only want the sub path from the id index to the end of the path (so we are only searching down
		// the tree instead of potentially up it)
		for pathIndex := idIndex; pathIndex < len(resourcePath); pathIndex++ {
			path := strings.Join(resourcePath[pathIndex:cap(resourcePath)], ":")

			clause := "(p.resource_path = ? AND ur.resource_id = ?)"

			subClauses = append(subClauses, clause)
			params = append(params, path, id)
		}

		clauses = append(clauses, fmt.Sprintf("(%s)", strings.Join(subClauses, " OR ")))
	}
	whereClause := strings.Join(clauses, " OR ")

	// add the action that we are also checking for
	params = append(params, action, userId)

	rows, err := repo.Query(
		fmt.Sprintf("select p.id, p.resource_path, ur.resource_id, p.action, ur.user_id from user_role ur join role_permission rp on rp.role_id = ur.role_id join permission p on p.id = rp.permission_id where %s AND p.action = ? AND ur.user_id = ?", whereClause),
		params...
	)
	if err != nil {
		log.Print(err)
		return false, errors.New("failed to check if resource can be accessed")
	}
	defer rows.Close()

	// basically we just need to check that something was returned
	count := 0
	for rows.Next() {
		count += 1
	}

	return count > 0, nil
}