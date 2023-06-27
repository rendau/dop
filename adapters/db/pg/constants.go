package pg

import (
	"regexp"
	"time"
)

type transactionCtxKeyT int8

const (
	ErrPrefix         = "pg-error"
	transactionCtxKey = transactionCtxKeyT(1)
)

var defaultOptions = OptionsSt{
	Timezone:          "Asia/Almaty",
	MaxConns:          100,
	MinConns:          5,
	MaxConnLifetime:   30 * time.Minute,
	MaxConnIdleTime:   15 * time.Minute,
	HealthCheckPeriod: 20 * time.Second,
	FieldTag:          "db",
}

var (
	queryParamRegexp = regexp.MustCompile(`(?si)\$\{[^}]+\}`)
)
