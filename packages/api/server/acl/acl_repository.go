package acl

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/transaction"
	"log"
	"strings"
)

type AclRepository struct {
	transaction.Transactionable
	server.Repository

	db *sql.DB
	tx *sql.Tx
}

func NewAclRepository(db *sql.DB, tx *sql.Tx) *AclRepository {
	repo := &AclRepository{
		db: db,
		tx: tx,
	}
	repo.CloneWithTransaction = func(tx *sql.Tx) interface{} {
		return NewAclRepository(repo.db, tx)
	}
	return repo
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
