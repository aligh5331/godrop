package repo

type Transaction interface {
	Commit() error
	Rollback() error
}

type TxKey struct{}
