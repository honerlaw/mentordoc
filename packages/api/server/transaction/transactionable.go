package transaction

import "database/sql"

type TransactionableObject interface {
	InjectTransaction(tx *sql.Tx) interface{}
}

type Transactionable struct {
	CloneWithTransaction func(tx *sql.Tx) interface{}
}

func (obj *Transactionable) InjectTransaction(tx *sql.Tx) interface{} {
	return obj.CloneWithTransaction(tx)
}
