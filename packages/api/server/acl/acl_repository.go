package acl

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server"
	"log"
	"reflect"
	"strings"
)

type AclRepository struct {
	server.Repository

	db *sql.DB
	tx *sql.Tx
}

func NewAclRepository(db *sql.DB, tx *sql.Tx) *AclRepository {
	repo := &AclRepository{
		db: db,
		tx: tx,
	}
	return repo
}

func (repo *AclRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewAclRepository(repo.db, tx)
}

func (repo *AclRepository) FindStatement(resourceName string, resourceId string, action string) (*Statement, error) {
	row := repo.QueryRow(
		"select id, resource_name, resource_id, action, created_at, updated_at, deleted_at from statement where resource_name = ? and resource_id = ? and action = ?",
		resourceName,
		resourceId,
		action,
	)

	statement := &Statement{}
	err := row.Scan(statement.Id, statement.ResourceName, statement.ResourceID, statement.Action, statement.CreatedAt, statement.UpdatedAt, statement.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return statement, nil
}

func (repo *AclRepository) FindStatements(resourceName string, resourceId string) ([]*Statement, error) {
	rows, err := repo.Query(
		"select id, resource_name, resource_id, action, created_at, updated_at, deleted_at from statement where resource_name = ? and resource_id = ?",
		resourceName,
		resourceId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statements := make([]*Statement, 0)

	for rows.Next() {
		statement := &Statement{}
		err := rows.Scan(statement.Id, statement.ResourceName, statement.ResourceID, statement.Action, statement.CreatedAt, statement.UpdatedAt, statement.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, nil
}

func (repo *AclRepository) InsertStatement(statement *Statement) (*Statement, error) {
	statement.CreatedAt = server.NowUnix()
	statement.UpdatedAt = server.NowUnix()

	_, err := repo.Exec(
		"insert into statement (id, resource_name, resource_id, action, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?)",
		statement.Id,
		strings.ToLower(statement.ResourceName),
		statement.ResourceID,
		strings.ToLower(statement.Action),
		statement.CreatedAt,
		statement.UpdatedAt,
		statement.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert statement")
	}

	return statement, nil;
}

func (repo *AclRepository) InsertPolicy(policy *Policy) (*Policy, error) {
	policy.CreatedAt = server.NowUnix()
	policy.UpdatedAt = server.NowUnix()

	_, err := repo.Exec(
		"insert into statement (id, created_at, updated_at, deleted_at) values (?, ?, ?, ?)",
		policy.Id,
		policy.CreatedAt,
		policy.UpdatedAt,
		policy.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to insert policy")
	}

	return policy, nil;
}

func (repo *AclRepository) LinkStatementToPolicy(policy *Policy, statement *Statement) error {
	_, err := repo.Exec(
		"insert into policy_statement (policy_id, statement_id) values (?, ?)",
		policy.Id,
		statement.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to link statement to policy")
	}

	return nil
}

func (repo *AclRepository) UnlinkStatementToPolicy(policy *Policy, statement *Statement) error {
	_, err := repo.Exec(
		"delete from policy_statement where policy_id = ? and statement_id = ?",
		policy.Id,
		statement.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to unlink statement from policy")
	}

	return nil
}

type FindPolicyForResourceResult struct {
	PolicyId string
	Id string
	ResourceName string
	ResourceId string
	Action string
}

func (repo *AclRepository) FindPolicyForResource(resourceName string, resourceId string, actions []string) (*Policy, error) {

	// another sscenario, 20 users in the orrganization, user A creates document A, we need to also crerate the statements and assign them to all users

	// instead we ccould denote a group of statements, these give access to all resourcces with the name (so no id)
	//

	// so we want to give acceess to all documents in a given folder to 100k users
	// we would needd to find eaech document, find the policy for each document, and then assign them the policy...

	// find all the statement ids for the given resource and actions

	// so how can we check if a policy exists with exactly those statements
	// we can find policies that contain them

	// assumptions
	// - every resource will basically have a statement for all possible combinations of actions on that resource... So we can assume when we create
	//   a resource, we can just immediately create all the statements for that resource, when we delete, we can remove them all
	// - policy - defines a set of statements, e.g. org:owner would have statements for the org, folder, document levels
	// - when we assign

	// thhis will select all policies and get all statements and everything across the board... next up we need to only
	// select all policies that contain the statements we need
	// then run a seccnod query to selecct all the statements in each of thhoose policies, then compaer the results to see if a policy matche

	// find all unique policy ids that match the given resources
	ids, err := repo.findPolicyIdsWithResources(resourceName, resourceId, actions)
	if err != nil {
		return nil, err
	}

	idsInterface := repo.convertToInterfaceArray(ids)

	// now lets find all of the policies and their statements
	rows, err := repo.Query(
		fmt.Sprintf("select p.id as policyId, s.id, s.resourceName, s.resourceId, s.action from policy p join policy_statement ps on p.id = ps.policy_id join statement s on s.id = ps.statement_id where p.id in (%s)", repo.buildPlaceholderString(ids)),
		idsInterface...
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close();

	// generate a map of policy id => statements in policy
	resultMap := make(map[string][]*FindPolicyForResourceResult)
	for rows.Next() {
		result := &FindPolicyForResourceResult{}
		err := rows.Scan(result.PolicyId, result.Id, result.ResourceName, result.ResourceId, result.Action)
		if err != nil {
			return nil, err
		}

		// does not exist in the map so add it
		if _, ok := resultMap[result.PolicyId]; !ok {
			resultMap[result.PolicyId] = make([]*FindPolicyForResourceResult, 0)
		}

		resultMap[result.PolicyId] = append(resultMap[result.PolicyId], result)
	}

	// next lets just check each to see if they contain the values we are looking for


}

// Find all of the policy ids that contain the given resources, we can then use this to find
func (repo *AclRepository) findPolicyIdsWithResources(resourceName string, resourceId string, actions []string) ([]int, error) {
	values := make([]interface{}, len(actions) + 2)
	values[0] = resourceName;
	values[1] = resourceId;
	for i, action := range actions {
		values[i + 2] = action
	}

	rows, err := repo.Query(
		fmt.Sprintf("select distinct p.id from policy p join policy_statement ps on p.id = ps.policy_id join statement s on s.id = ps.statement_id and s.resourceName = ? and s.resourceId = ? and action in (%s)", repo.buildPlaceholderString(actions)),
		values...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (repo *AclRepository) buildPlaceholderString(values interface{}) string {
	val := reflect.ValueOf(values)
	placeholders := make([]string, val.Len())
	for i := 0; i < val.Len(); i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func (repo *AclRepository) convertToInterfaceArray(values interface{}) []interface{} {
	val := reflect.ValueOf(values)
	target := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		target[i] = val.Index(i).Interface()
	}
	return target;
}