package transaction

import (
	"database/sql"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (factory *TransactionManager) Transact(obj TransactionableObject, handle func(obj interface{}) (interface{}, error)) (interface{}, error) {
	tx, err := factory.db.Begin();
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
