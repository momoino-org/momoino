package testutils

import (
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

type AnyUUIDArg struct{}

var _ sqlmock.Argument = (*AnyUUIDArg)(nil)

func (AnyUUIDArg) Match(v driver.Value) bool {
	_, ok := v.(uuid.UUID)

	if ok {
		return true
	}

	if err := uuid.Validate(v.(string)); err == nil {
		return true
	}

	return false
}
