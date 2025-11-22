package services

import "fmt"

type ErrUserExistsInAnotherTeam struct{}

func (e *ErrUserExistsInAnotherTeam) Error() string {
	return "Some of users already exists in another team"
}

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "resource not found"
}

type ErrTeamExists struct {
	TeamName string
}

func (e *ErrTeamExists) Error() string {
	return fmt.Sprintf("%s already exists", e.TeamName)
}
