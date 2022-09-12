package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andygrunwald/go-jira/v2/cloud/models/apps/insights"
)

// GetObjectSchemaList resource to find object schemas
// Reference: https://developer.atlassian.com/cloud/insight/rest/api-group-objectschema/#api-objectschema-list-get
func (i *InsightsService) GetObjectSchemaList(ctx context.Context, workspaceID string) (*insights.GenericList[insights.ObjectSchema], error) {
	apiEndPoint := fmt.Sprintf(`%s/jsm/insight/workspace/%s/v1/objectschema/list`, insightsURL, workspaceID)

	req, err := i.client.NewRequest(ctx, http.MethodGet, apiEndPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := i.client.client.Do(req)
	if err != nil {
		return nil, err
	}

	list := new(insights.GenericList[insights.ObjectSchema])
	err = json.NewDecoder(res.Body).Decode(&list)

	return list, err
}

// GetObjectSchemaAttributes find all object type attributes for this object schema
// Reference: https://developer.atlassian.com/cloud/insight/rest/api-group-objectschema/#api-objectschema-id-attributes-get
func (i *InsightsService) GetObjectSchemaAttributes(ctx context.Context, workspaceID, id string) ([]insights.ObjectTypeAttribute, error) {
	apiEndPoint := fmt.Sprintf(`%s/jsm/insight/workspace/%s/v1/objectschema/%s/attributes`, insightsURL, workspaceID, id)

	req, err := i.client.NewRequest(ctx, http.MethodGet, apiEndPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := i.client.client.Do(req)
	if err != nil {
		return nil, err
	}

	var attributes []insights.ObjectTypeAttribute
	err = json.NewDecoder(res.Body).Decode(&attributes)

	return attributes, err
}
