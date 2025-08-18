package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Datastore struct {
	client *datastore.Client
}

func NewDatastoreClient(ctx context.Context, projectID string) (*Datastore, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create datastore client: %w", err)
	}
	return &Datastore{client: client}, nil
}

func (d *Datastore) PutMulti(ctx context.Context, kind string, keys []string, entities []interface{}) error {
	if len(entities) == 0 {
		return fmt.Errorf("no entities to put")
	}

	dsKeys := make([]*datastore.Key, len(entities))
	for i, key := range keys {
		dsKeys[i] = datastore.NameKey(kind, key, nil)
	}

	if _, err := d.client.PutMulti(ctx, dsKeys, entities); err != nil {
		return fmt.Errorf("failed to put entities: %w", err)
	}

	return nil
}
