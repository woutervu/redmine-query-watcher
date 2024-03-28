package main

import "fmt"

func main() {
	rs, err := getRedmineService()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	is, err := rs.GetIssuesByQueryId(1274)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(is)
}
