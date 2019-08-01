package snapstore

import (
	"context"
	"encoding/json"
)

// SnapshotStore is interface for aggregate storage
type SnapshotStore interface {
	Save(ctx context.Context, id string, data json.RawMessage) error
	Find(ctx context.Context, id string) (json.RawMessage, error)
}

// EmptySnapshotStore is empty snapshot store implementation
type EmptySnapshotStore struct{}

// Save is empty snapshot store save implementation
func (EmptySnapshotStore) Save(ctx context.Context, id string, data json.RawMessage) error {
	return nil
}

// Find is empty snapshot store find implementation
func (EmptySnapshotStore) Find(ctx context.Context, id string) (json.RawMessage, error) {
	return nil, nil
}
