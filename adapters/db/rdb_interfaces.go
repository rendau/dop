package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type RDBFull interface {
	RDBConnectionWithHelpers
	RDBContextTransaction
}

type RDBMin interface {
	RDBConnection
	RDBContextTransaction
}

type RDBConnection interface {
	DbExec(ctx context.Context, sql string, args ...any) error
	DbQuery(ctx context.Context, sql string, args ...any) (RDBRows, error)
	DbQueryRow(ctx context.Context, sql string, args ...any) RDBRow
	DbExecM(ctx context.Context, sql string, argMap map[string]any) error
	DbQueryM(ctx context.Context, sql string, argMap map[string]any) (RDBRows, error)
	DbQueryRowM(ctx context.Context, sql string, argMap map[string]any) RDBRow
	HErr(err error) error
}

type RDBConnectionWithHelpers interface {
	RDBConnection

	HfList(ctx context.Context, ops RDBListOptions) (int64, error)
	HfGenerateSort(rNames []string, allowed map[string]string) []string
	HfGet(ctx context.Context, ops RDBGetOptions) error
	HfCreate(ctx context.Context, ops RDBCreateOptions) error
	HfUpdate(ctx context.Context, ops RDBUpdateOptions) error
	HfGetCUFields(obj any) (map[string]any, map[string]bool)
	HfOptionalWhere(conds []string) string
	HfDelete(ctx context.Context, ops RDBDeleteOptions) error
}

type RDBContextTransaction interface {
	RenewContextTransaction(ctx context.Context) error
	TransactionFn(ctx context.Context, f func(context.Context) error) error
	TransactionAddAsyncCallback(ctx context.Context, f func())
}

type RDBRows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type RDBRow interface {
	Scan(dest ...any) error
}

type RDBConSt interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
