package esrepo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

// SnapshotSaver return function that will save event to snapshot
func SnapshotSaver(store snapstore.SnapshotStore, newAggregate esfazz.AggregateFactory) EventListener {
	return func(ctx context.Context, events ...*esfazz.Event) error {
		for _, ev := range events {
			agg := newAggregate(ev.Aggregate.Id)

			jsonSnapshot, err := store.Find(ctx, agg.GetId())
			if err != nil {
				return err
			}

			if jsonSnapshot != nil {
				err = json.Unmarshal(jsonSnapshot, agg)
				if err != nil {
					return err
				}
			}

			if agg.GetVersion() != ev.Aggregate.Version {
				return errors.New("cannot apply event: aggregate version is different from snapshot")
			}

			err = agg.Apply(ev)
			if err != nil {
				return err
			}

			jsonNewSnap, err := json.Marshal(agg)
			if err != nil {
				return err
			}

			err = store.Save(ctx, agg.GetId(), jsonNewSnap)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
