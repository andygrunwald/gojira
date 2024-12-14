package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	jira "github.com/andygrunwald/go-jira/v2/onpremise"
	"golang.org/x/term"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("Jira URL: ")
	jiraURL, _ := r.ReadString('\n')

	fmt.Print("Jira Username: ")
	username, _ := r.ReadString('\n')

	fmt.Print("Jira Password: ")
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client, err := jira.NewClient(strings.TrimSpace(jiraURL), tp.Client())
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				AccountID: "my-user-account-id",
			},
			Reporter: &jira.User{
				AccountID: "your-user-account-id",
			},
			Description: "Test Issue",
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: "PROJ1",
			},
			Summary: "Just a demo issue",
		},
	}

	issue, _, err := client.Issue.Create(context.Background(), &i, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Self)
}
