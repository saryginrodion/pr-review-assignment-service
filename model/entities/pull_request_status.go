package entities

import "database/sql/driver"

type PullRequestStatus string

const (
	PULL_REQUEST_OPEN   = PullRequestStatus("OPEN")
	PULL_REQUEST_MERGED = PullRequestStatus("MERGED")
)

// From https://gorm.io/docs/data_types.html

func (s *PullRequestStatus) Scan(value any) error {
	*s = PullRequestStatus(value.([]byte))
	return nil
}

func (self PullRequestStatus) Value() (driver.Value, error) {
	return string(self), nil
}
