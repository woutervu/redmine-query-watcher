package main

import "fmt"

func main() {
	rs, err := getRedmineService()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c, err := getConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	is, err := rs.GetIssuesByQueryId(c.QueryId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, issue := range is {
		fmt.Println(issue)
	}
}
