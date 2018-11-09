package natsutil

import (
	"fmt"
	"sync"

	"github.com/nats-io/nuid"

	"github.com/nats-io/go-nats-streaming"
)

type StreamingComponent struct {
	cmu sync.Mutex

	id string

	nc stan.Conn

	kind string
}

func NewStreamingComponent(kind string) *StreamingComponent {
	id := nuid.Next()
	return &StreamingComponent{
		id:   id,
		kind: kind,
	}
}

func (s *StreamingComponent) ConnectToNatsStreamingService(clusterId string, options ...stan.Option) error {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	nc, err := stan.Connect(clusterId, s.id, options...)
	if err != nil {
		return err
	}
	s.nc = nc
	return nil
}

func (s *StreamingComponent) NATS() stan.Conn {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	return s.nc
}

func (s *StreamingComponent) Name() string {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	return fmt.Sprintf("%s:%s", s.kind, s.id)
}

func (s *StreamingComponent) Close() error {
	s.NATS().Close()
	return nil
}
