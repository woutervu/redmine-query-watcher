package main

import (
	"fmt"

	"github.com/mattn/go-redmine"
)

type Issue struct {
	Id          int
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
	GetIssueById(id int) (*Issue, error)
	GetIssuesByQueryId(queryId int) ([]*Issue, error)
}

type RedmineService struct {
	Config        *Config
	RedmineClient *redmine.Client
}

func (r *RedmineService) GetIssueById(id int) (*Issue, error) {
	client, err := r.getClient()
	if err != nil {
		return nil, err
	}

	ri, err := client.Issue(id)
	if err != nil {
		return nil, err
	}

	issue, err := r.transformIssue(ri)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (r *RedmineService) GetIssuesByQueryId(queryId int) ([]*Issue, error) {
	client, err := r.getClient()
	if err != nil {
		return nil, err
	}

	ris, err := client.IssuesByQuery(r.Config.QueryId)
	if err != nil {
		return nil, err
	}

	c := len(ris)
	issues := make([]*Issue, c)

	for i, issue := range ris {
		ti, err := r.transformIssue(&issue)
		if err != nil {
			return nil, err
		}

		issues[i] = ti
	}

	return issues, nil
}

func (r *RedmineService) getClient() (*redmine.Client, error) {
	if r.RedmineClient != nil {
		return r.RedmineClient, nil
	}

	if r.Config == nil {
		return nil, fmt.Errorf("config was not properly set")
	}

	r.RedmineClient = redmine.NewClient(r.Config.RedmineURL, r.Config.RedmineToken)

	return r.RedmineClient, nil
}

func (r *RedmineService) transformIssue(ri *redmine.Issue) (*Issue, error) {
	if ri == nil {
		return nil, fmt.Errorf("nil pointer `*redmine.Issue`")
	}

	assignedToName := ""
	if ri.AssignedTo != nil {
		assignedToName = ri.AssignedTo.Name
	}

	issue := Issue{
		Id:          ri.Id,
		ProjectCode: ri.Project.Name[:3],
		Subject:     ri.Subject,
		Status:      r.getStatus(ri.Status.Id),
		Assignee:    assignedToName,
	}

	return &issue, nil
}

func (r *RedmineService) getStatus(statusId int) status {
	return status(statusId)
}

func getRedmineService() (*RedmineService, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	return &RedmineService{Config: config}, nil
}
