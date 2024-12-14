package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/trivago/tgo/tcontainer"
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

	fmt.Print("Custom field name (i.e. customfield_10220): ")
	customFieldName, _ := r.ReadString('\n')

	fmt.Print("Custom field value: ")
	customFieldValue, _ := r.ReadString('\n')

	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		APIToken: strings.TrimSpace(password),
	}

	client, err := jira.NewClient(strings.TrimSpace(jiraURL), tp.Client())
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		os.Exit(1)
	}

	unknowns := tcontainer.NewMarshalMap()
	unknowns[customFieldName] = customFieldValue

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: "myuser",
			},
			Reporter: &jira.User{
				Name: "youruser",
			},
			Description: "Test Issue",
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: "PROJ1",
			},
			Summary:  "Just a demo issue",
			Unknowns: unknowns,
		},
	}

	issue, _, err := client.Issue.Create(context.Background(), &i, &jira.CreateIssueOptions{
		UpdateHistory: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %v\n", issue.Key, issue.Self)
}
