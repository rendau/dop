package mock

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/rendau/dop/adapters/logger"
)

type St struct {
	lg      logger.Lite
	testing bool

	q  []Req
	mu sync.Mutex

	smsCodeRegexp *regexp.Regexp
}

type Req struct {
	Phone string
	Msg   string
}

func New(lg logger.Lite, testing bool) *St {
	return &St{
		lg:            lg,
		testing:       testing,
		q:             make([]Req, 0),
		smsCodeRegexp: regexp.MustCompile(`([0-9]{4})`),
	}
}

func (m *St) Send(phone string, msg string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.testing {
		m.lg.Infow("Sms sent", "phone", phone, "msg", msg)
		return true
	}

	req := Req{
		Phone: phone,
		Msg:   msg,
	}

	if len(m.q) > 100 {
		m.q = make([]Req, 0)
	}

	m.q = append(m.q, req)

	return true
}

func (m *St) SendAsync(phone string, msg string) bool {
	return m.Send(phone, msg)
}

func (m *St) PullAll() []Req {
	m.mu.Lock()
	defer m.mu.Unlock()

	q := m.q

	m.q = make([]Req, 0)

	return q
}

func (m *St) PullCode() int {
	smsReqs := m.PullAll()
	if len(smsReqs) < 1 {
		return 0
	}

	matches := m.smsCodeRegexp.FindStringSubmatch(smsReqs[0].Msg)
	if len(matches) == 2 {
		code, _ := strconv.ParseInt(matches[1], 10, 64)
		return int(code)
	}

	return 0
}

func (m *St) Clean() {
	_ = m.PullAll()
}
