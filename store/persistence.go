package store

import "github.com/zachlefevre/project_hopper_backend/com"

type persistence interface {
	add(event *pb.Event) error
	getAll() []*pb.Event
}

type local_persistence struct {
	cache []*pb.Event
}

func (lp *local_persistence) add(event *pb.Event) error {
	lp.cache = append(lp.cache, event)
	return nil
}

func (lp local_persistence) getAll() []*pb.Event {
	return lp.cache
}
