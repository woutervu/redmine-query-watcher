package main

import (
	"fmt"
	"math/rand"

	"github.com/mattn/go-redmine"
)

var anon bool = true

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
	Unknown
)

var redmineServiceInstance *RedmineService

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
		Status:      getStatus(ri.Status),
		Assignee:    assignedToName,
	}

	if anon == true {
		issue.anonymizeIssue()
	}

	return &issue, nil
}

func (i *Issue) anonymizeIssue() {
	id := rand.Intn(199999) + 100000
	pcs := []string{"ABC", "DEF", "GHI", "JKL"}
	sub := []string{"Can't login!", "Issue with customer", "HTTP 500 error on checkout", "URGENT: DEADLOCK ERRORS"}
	st := []status{New, InProgress, OnHold, Solved}
	as := []string{"John Doe", "Jane Doe", "Walter White", "Jesse Pinkman", "Leslie Pollos"}

	i.Id = id
	i.ProjectCode = getRandomString(pcs)
	i.Subject = getRandomString(sub)
	i.Status = getRandomStatus(st)
	i.Assignee = getRandomString(as)
}

func getStatus(redmineStatus *redmine.IdName) status {
	switch redmineStatus.Name {
	case "New":
		return New
	case "Support: solved":
		return Solved
	case "Support: on hold":
		return OnHold
	case "Support: in progress":
		return InProgress
	case "Closed":
		return Closed
	}
	return Unknown
}

func getRedmineService() (*RedmineService, error) {
	if redmineServiceInstance != nil {
		return redmineServiceInstance, nil
	}

	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	rs := RedmineService{Config: config}
	redmineServiceInstance = &rs

	return redmineServiceInstance, nil
}

func getRandomString(s []string) string {
	length := len(s)
	if length == 0 {
		return ""
	}

	index := rand.Intn(length - 1)

	return s[index]
}

func getRandomStatus(st []status) status {
	length := len(st)
	if length == 0 {
		return Unknown
	}

	index := rand.Intn(length - 1)

	return st[index]
}
