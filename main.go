package main

import (
	"fmt"
	"os"
	"time"
)

type Msg string

var lastUpdate int64 = time.Now().Unix() - 30
var issueChannel chan []*Issue = make(chan []*Issue, 1)

func main() {
	ec, err := appRun()
	if err != nil {
		fmt.Printf("Exit code: %d\nMessage: %s", ec, err)
		os.Exit(ec)
	}

	os.Exit(ec)
}

func appRun() (int, error) {
	m, err := getModel()
	if err != nil {
		return 1, err
	}

	tp, err := getTeaProgram(m)
	if err != nil {
		return 1, err
	}

	rs, err := getRedmineService()
	if err != nil {
		return 1, err
	}

	go func() {
		for {
			now := time.Now().Unix()
			nextRun := lastUpdate + 30
			timeToSleep := 30
			if now < nextRun {
				timeToSleep = int(lastUpdate) - int(now)
				time.Sleep(time.Second * time.Duration(timeToSleep))
			}

			ri, _ := rs.GetIssuesByQueryId(rs.Config.QueryId)
			issueChannel <- ri
			tp.Send(Msg("U"))

			lastUpdate = time.Now().Unix()
		}
	}()

	_, err = tp.Run()

	if err != nil {
		return 1, err
	}

	return 0, nil
}
