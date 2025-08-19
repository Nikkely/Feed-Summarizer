package database

import "context"

// Database defines the interface for data persistence operations.
// This interface abstracts the underlying storage mechanism,
// allowing for different implementations (e.g., Cloud Datastore, MySQL, etc.).
type Database interface {
	// PutMulti stores multiple entities in the specified table.
	// Parameters:
	//   - ctx: The context for the operation
	//   - table: The name of the table or collection to store the entities
	//   - key: The unique identifier for the group of entities
	//   - entities: The slice of entities to store
	//
	// Returns:
	//   - error: An error if the operation fails
	PutMulti(ctx context.Context, table string, key string, entities []any) error
}
