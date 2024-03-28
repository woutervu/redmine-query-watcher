package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type IssueCol []*Issue

type QueryCol map[int][]*Issue

type MockIssueService struct {
	CurrentId       int
	IssueCollection IssueCol
	QueryCollection QueryCol
}

func TestGetIssueByIdExistingIssue(t *testing.T) {
	assert := assert.New(t)

	is := getMockService()
	is.generateIssues(5)
	issue := is.IssueCollection[0]

	fetchedIssue, err := is.GetIssueById(issue.Id)

	assert.NotNil(fetchedIssue)
	assert.Nil(err)
	assert.Equal(issue, fetchedIssue)
}

func TestGetIssueByIdNonExistentIssue(t *testing.T) {
	assert := assert.New(t)

	is := getMockService()
	issueId := 145678
	fetchedIssue, err := is.GetIssueById(issueId)
	assert.Nil(fetchedIssue)
	assert.Error(err)
}

func TestGetIssuesByQueryIdExistent(t *testing.T) {
	assert := assert.New(t)

	is := getMockService()
	queryId := 1359
	is.generateQueryIdWithIssues(queryId, 5)
	expected := is.QueryCollection[queryId]

	issues, err := is.GetIssuesByQueryId(queryId)
	assert.Equal(expected, issues)
	assert.Nil(err)
}

func TestGetIssuesByQueryIdNonExistent(t *testing.T) {
	assert := assert.New(t)

	is := getMockService()
	is.generateQueryIdWithIssues(1359, 1)

	issues, err := is.GetIssuesByQueryId(1234)
	assert.Nil(issues)
	assert.Error(err)
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
		CurrentId:       100000,
		QueryCollection: make(QueryCol),
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
	if amount <= 0 {
		return
	}

	lastIndex := len(m.IssueCollection)
	if lastIndex > 0 {
		lastIndex = lastIndex - 1
	}

	m.generateIssues(amount)

	from := lastIndex
	to := lastIndex + amount
	if lastIndex > 0 {
		from = from + 1
		to = to + 1
	}

	m.QueryCollection[queryId] = append(m.QueryCollection[queryId], m.IssueCollection[from:to]...)
}
