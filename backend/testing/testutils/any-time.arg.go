package testutils

import (
	"database/sql/driver"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type AnyTimeArg struct{}

var _ sqlmock.Argument = (*AnyTimeArg)(nil)

func (AnyTimeArg) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
