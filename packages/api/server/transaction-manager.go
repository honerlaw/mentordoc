package server

import (
	"database/sql"
)

type Transactionable interface {
	InjectTransaction(tx *sql.Tx) interface{}
}

type TransactionManager struct {
	db *sql.DB
	tx *sql.Tx
}

func NewTransactionManager(db *sql.DB, tx *sql.Tx) *TransactionManager {
	return &TransactionManager{
		db: db,
		tx: tx,
	}
}

func (manager *TransactionManager) InjectTransaction(tx *sql.Tx) interface{} {
	return NewTransactionManager(manager.db, tx)
}

func (manager *TransactionManager) Transact(obj Transactionable, handle func(obj interface{}) (interface{}, error)) (interface{}, error) {
	tx, err := manager.getTransaction()
	if err != nil {
		return nil, err
	}

	injected := obj.InjectTransaction(tx)

	resp, err := handle(injected)
	
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}

func (manager *TransactionManager) getTransaction() (*sql.Tx, error) {
	if manager.tx != nil {
		return manager.tx, nil
	}
	tx, err := manager.db.Begin();
	if err != nil {
		return nil, err
	}
	return tx, nil
}
