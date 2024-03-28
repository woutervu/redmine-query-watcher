package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

type IssueCol []*Issue

type QueryCol map[int][]*Issue

type MockIssueService struct {
	CurrentId       int
	IssueCollection IssueCol
	QueryCollection QueryCol
}

func TestGetIssueByIdExistingIssue(t *testing.T) {
	is := getMockService()
	is.generateIssues(5)
	issue := is.IssueCollection[0]

	fetchedIssue, err := is.GetIssueById(issue.Id)
	if err != nil {
		t.Errorf("Expected *Issue, got error: %s", err)
	}

	if fetchedIssue == nil {
		t.Errorf("Expected Issue with ID `%d`, got nil.", issue.Id)
	}

	if fetchedIssue != nil && fetchedIssue.Id != issue.Id {
		t.Errorf("Expected Issue with ID `%d`, got %d.", issue.Id, fetchedIssue.Id)
	}
}

func TestGetIssueByIdNonExistentIssue(t *testing.T) {
	is := getMockService()
	issueId := 145678
	fetchedIssue, err := is.GetIssueById(issueId)
	if fetchedIssue != nil && err == nil {
		t.Errorf("Expected nil *Issue, got issue with ID `%d`", fetchedIssue.Id)
	}
}

func TestGetIssuesByQueryIdExistent(t *testing.T) {
	is := getMockService()
	queryId := 1359
	amount := 5
	is.generateQueryIdWithIssues(queryId, amount)

	issues, err := is.GetIssuesByQueryId(queryId)
	if err != nil {
		t.Errorf("Expected []*Issue, got error: %s", err)
	}

	if len(issues) != amount {
		t.Errorf("Expected %d issues, got %d", amount, len(issues))
	}
}

func TestGetIssuesByQueryIdNonExistent(t *testing.T) {
	is := getMockService()
	queryId := 1359
	amount := 1
	is.generateQueryIdWithIssues(queryId, amount)

	issues, err := is.GetIssuesByQueryId(queryId + 1)
	if issues != nil {
		t.Errorf("Expected error, got []*Issue with `%d` items.", len(issues))
	}

	if err == nil {
		t.Error("Expected error, got nil.")
	}
}

func (m *MockIssueService) GetIssueById(id int) (*Issue, error) {
	for _, issue := range m.IssueCollection {
		if issue.Id == id {
			return issue, nil
		}
	}

	return nil, fmt.Errorf("Issue with ID `%d` not found.", id)
}

func (m *MockIssueService) GetIssuesByQueryId(queryId int) ([]*Issue, error) {
	queryCol, exists := m.QueryCollection[queryId]

	if !exists {
		return nil, fmt.Errorf("Query with ID `%d` doesn't exist.", queryId)
	}

	return queryCol, nil
}

func getMockService() MockIssueService {
	return MockIssueService{
		CurrentId: 100000,
	}
}

func (m *MockIssueService) generateIssues(amount int) {
	projectCodes := make([]string, 0, 5)
	projectCodes = append(projectCodes, "ABC", "DEF", "GHI", "JKL", "MNO")
	assignees := make([]string, 0, 6)
	assignees = append(assignees, "", "John Doe", "Jane Doe", "Walter White", "Jesse Pinkman", "Leslie Pollos")

	for i := 0; i < amount; i++ {
		m.CurrentId = m.CurrentId + 1
		issue := Issue{
			Id:          m.CurrentId,
			ProjectCode: projectCodes[rand.Intn(4)],
			Subject:     "Test issue #" + strconv.Itoa(m.CurrentId),
			Status:      status(rand.Intn(4)),
			Assignee:    assignees[rand.Intn(5)],
		}

		m.IssueCollection = append(m.IssueCollection, &issue)
	}
}

func (m *MockIssueService) generateQueryIdWithIssues(queryId int, amount int) {
	if len(m.IssueCollection) < amount {
		m.generateIssues(amount - len(m.IssueCollection))
	}
	m.QueryCollection = make(QueryCol)
	m.QueryCollection[queryId] = append(m.QueryCollection[queryId], m.IssueCollection[:amount]...)
}
