package main

import (
	"context"
	"fmt"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
)

func main() {
	jiraURL := "https://go-jira-opensource.atlassian.net/"

	// Jira docs: https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/
	// Create a new API token: https://id.atlassian.com/manage-profile/security/api-tokens
	tp := jira.BasicAuthTransport{
		Username: "<username>",
		APIToken: "<api-token>",
	}
	client, err := jira.NewClient(jiraURL, tp.Client())
	if err != nil {
		panic(err)
	}

	component, _, err := client.Component.Get(context.Background(), "10000")
	if err != nil {
		panic(err)
	}
	if component == nil {
		fmt.Println("Component not found")
	}

	// Update component
	component.Name = "New Name"

	updatedComponent, _, err := client.Component.Update(context.Background(), component)
	if err != nil {
		panic(err)
	}

	// updatedComponent SHOULD -in-theory- be the same as the component sent to the update method
	if updatedComponent.Name != component.Name {
		fmt.Println("This should not happen, received component different from the sent one!")
	}

	fmt.Printf("Updated component: %+v\n", updatedComponent)
	fmt.Println("Success!")
}
