package main

import (
	"context"
	"fmt"

	jira "github.com/andygrunwald/go-jira/v2/onpremise"
)

func main() {
	jiraURL := "<your-jira-instance>"

	// See "Using Personal Access Tokens"
	// https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html
	tp := jira.BearerAuthTransport{
		Token: "<persona-access-token>",
	}
	client, err := jira.NewClient(tp.Client(), jiraURL)
	if err != nil {
		panic(err)
	}

	u, _, err := client.User.GetSelf(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Email: %v\n", u.EmailAddress)
	fmt.Println("Success!")
}
