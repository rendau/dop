package mock

import (
	"sync"

	"github.com/rendau/dop/adapters/logger"
)

type St struct {
	lg      logger.Lite
	testing bool

	q  []Req
	mu sync.Mutex
}

type Req struct {
	Topic string
	Key   string
	Value any
}

func New(lg logger.Lite, testing bool) *St {
	return &St{
		lg:      lg,
		testing: testing,
		q:       make([]Req, 0),
	}
}

func (m *St) SendJson(topic, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.testing {
		m.lg.Infow("KRP: SendJson", "topic", topic, "key", key, "value", value)
		return nil
	}

	req := Req{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	if len(m.q) > 100 {
		m.q = make([]Req, 0)
	}

	m.q = append(m.q, req)

	return nil
}

func (m *St) SendManyJson(topic, key string, value []any) error {
	var err error

	for _, v := range value {
		err = m.SendJson(topic, key, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *St) PullAll() []Req {
	m.mu.Lock()
	defer m.mu.Unlock()

	q := m.q

	m.q = make([]Req, 0)

	return q
}

func (m *St) Clean() {
	_ = m.PullAll()
}
