package entities

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type PullRequestStatus string

const (
	PULL_REQUEST_OPEN   = PullRequestStatus("OPEN")
	PULL_REQUEST_MERGED = PullRequestStatus("MERGED")
)

// From https://gorm.io/docs/data_types.html

func (s *PullRequestStatus) Scan(value any) error {
	if value == nil {
		return errors.New("Failed to parse PullRequestStatus")
	}

	switch value := value.(type) {
	case string:
		*s = PullRequestStatus(value)
	case []byte:
		*s = PullRequestStatus(value)
	default:
		return fmt.Errorf("Cannot convert %T to PullRequestStatus", value)
	}
	return nil
}

func (self PullRequestStatus) Value() (driver.Value, error) {
	return string(self), nil
}
