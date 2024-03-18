package main

type Issue struct {
	Id          string
	ProjectCode string
	Subject     string
	Status      status
	Assignee    string
}

type status int

const (
	New status = iota
	InProgress
	OnHold
	Solved
	Closed
)

type IssueServiceInterface interface {
	GetIssueById(id string) (*Issue, error)
	GetIssuesByQueryId(queryId int) ([]*Issue, error)
}

type RedmineService struct{}

func (r *RedmineService) GetIssueById(id string) (*Issue, error) {
	return nil, nil
}

func (r *RedmineService) GetIssuesByQueryId(queryId int) ([]*Issue, error) {
	return nil, nil
}
