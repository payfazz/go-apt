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
