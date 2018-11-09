package store

import "github.com/zachlefevre/project_hopper_backend/com"

type EventStore struct {
	history persistence
}

func (store EventStore) PersistEvent(event *pb.Event) error {
	return store.history.add(event)
}

func (store EventStore) GetEvents() []*pb.Event {
	return store.history.getAll()
}

func NewEventStore() EventStore {
	return EventStore{history: &local_persistence{}}
}
