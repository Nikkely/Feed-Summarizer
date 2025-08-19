package database

import "context"

type Database interface {
	PutMulti(ctx context.Context, table string, key string, entities []any) error
}
