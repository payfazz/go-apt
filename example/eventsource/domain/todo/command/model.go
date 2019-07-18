package command

import (
	"encoding/json"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo/data"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
	"time"
)

type Todo struct {
	Id        string     `json:"id"`
	Text      string     `json:"text"`
	Completed bool       `json:"completed"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (t *Todo) ApplyAll(events ...*fazzeventsource.Event) error {
	for _, ev := range events {
		err := t.Apply(ev)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Todo) Apply(event *fazzeventsource.Event) error {
	switch event.Type {
	case data.EVENT_TODO_CREATED:
		return t.ApplyTodoCreated(event)
	case data.EVENT_TODO_UPDATED:
		return t.ApplyTodoUpdated(event)
	case data.EVENT_TODO_DELETED:
		return t.ApplyTodoDeleted(event)
	}
	return nil
}

func (t *Todo) ApplyTodoCreated(event *fazzeventsource.Event) error {
	payload := &data.TodoCreated{}
	err := json.Unmarshal(event.Data, payload)
	if err != nil {
		return err
	}

	t.Text = payload.Text
	t.CreatedAt = event.CreatedAt
	t.UpdatedAt = event.CreatedAt
	return nil
}

func (t *Todo) ApplyTodoUpdated(event *fazzeventsource.Event) error {
	payload := &data.TodoUpdated{}
	err := json.Unmarshal(event.Data, payload)
	if err != nil {
		return err
	}
	if payload.Text != nil {
		t.Text = *payload.Text
	}
	if payload.Completed != nil {
		t.Completed = *payload.Completed
	}
	t.UpdatedAt = event.CreatedAt
	return nil

}

func (t *Todo) ApplyTodoDeleted(event *fazzeventsource.Event) error {
	t.DeletedAt = event.CreatedAt
	return nil
}
