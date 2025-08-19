// Package database provides interfaces and implementations for data persistence.
package database

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
)

// Datastore provides access to Google Cloud Datastore operations.
// It wraps the datastore.Client to provide simplified access to common operations.
type Datastore struct {
	client *datastore.Client
}

// NewDatastoreClient creates a new Datastore instance with the specified project ID.
// Parameters:
//   - ctx: The context for the client connection
//   - projectID: The Google Cloud project ID to connect to
//
// Returns:
//   - *Datastore: A new Datastore instance
//   - error: An error if client creation fails
func NewDatastoreClient(ctx context.Context, projectID string) (*Datastore, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create datastore client: %w", err)
	}
	return &Datastore{client: client}, nil
}

// PutMulti stores multiple entities in Datastore.
// Each entity is stored with a named key and must be of type map[string]any.
// Parameters:
//   - ctx: The context for the operation
//   - kind: The kind (entity type) to store
//   - keys: The string keys for each entity
//   - entities: The entities to store, must be []map[string]any
//
// Returns:
//   - error: An error if the operation fails or if entities are of wrong type
func (d *Datastore) PutMulti(ctx context.Context, kind string, keys []string, entities []interface{}) error {
	if len(entities) == 0 {
		return fmt.Errorf("no entities to put")
	}

	dsKeys := make([]*datastore.Key, len(entities))
	dsEntities := make([]dynamicEntity, len(entities))

	for i, key := range keys {
		dsKeys[i] = datastore.NameKey(kind, key, nil)

		switch v := entities[i].(type) {
		case map[string]any:
			dsEntities[i] = dynamicEntity(v)
		default:
			return fmt.Errorf("unsupported entity type at index %d: %T", i, v)
		}
	}

	if _, err := d.client.PutMulti(ctx, dsKeys, dsEntities); err != nil {
		return fmt.Errorf("failed to put entities: %w", err)
	}

	return nil
}

// GenerateUUID generates a new UUID string.
// This function uses github.com/google/uuid to generate a Version 1 UUID,
// which is based on timestamp and MAC address.
//
// Returns:
//   - string: The generated UUID as a string
//   - error: An error if UUID generation fails
func GenerateUUID() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}
	return uuid.String(), nil
}

// GenerateUUIDs generates multiple UUID strings.
// This is a convenience function that calls GenerateUUID multiple times.
// If any UUID generation fails, all errors are collected and returned.
//
// Parameters:
//   - len: The number of UUIDs to generate
//
// Returns:
//   - []string: A slice of generated UUIDs
//   - error: An error if any UUID generation fails
func GenerateUUIDs(len int) ([]string, error) {
	var err error
	uuids := make([]string, len)
	for i := 0; i < len; i++ {
		uuid, uuidErr := GenerateUUID()
		if uuidErr != nil {
			err = errors.Join(err, uuidErr)
		}
		uuids[i] = uuid
	}
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUIDs: %w", err)
	}
	return uuids, nil
}

// dynamicEntity implements datastore.PropertyLoadSaver interface.
// It allows dynamic handling of arbitrary entity properties without
// requiring a predefined struct type.
type dynamicEntity map[string]any

// Save converts the dynamicEntity map to a PropertyList for Datastore storage.
// This method satisfies the datastore.PropertySaver interface.
//
// Returns:
//   - []datastore.Property: A list of properties ready for storage
//   - error: Always returns nil as map conversion cannot fail
func (d dynamicEntity) Save() ([]datastore.Property, error) {
	props := make([]datastore.Property, 0, len(d))
	for k, v := range d {
		props = append(props, datastore.Property{
			Name:  k,
			Value: v,
		})
	}
	return props, nil
}

// Load restores a dynamicEntity from a PropertyList retrieved from Datastore.
// This method satisfies the datastore.PropertyLoader interface.
//
// Parameters:
//   - props: The list of Datastore properties to convert
//
// Returns:
//   - error: Always returns nil as map loading cannot fail
func (d *dynamicEntity) Load(props []datastore.Property) error {
	m := make(map[string]any, len(props))
	for _, p := range props {
		m[p.Name] = p.Value
	}
	*d = m
	return nil
}
