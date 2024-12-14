package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFilterService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/2/filter"
	raw, err := ioutil.ReadFile("./mocks/all_filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEndpoint)
		fmt.Fprint(writer, string(raw))
	})

	filters, _, err := testClient.Filter.GetList()
	if filters == nil {
		t.Error("Expected Filters list. Filters list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestFilterService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/2/filter/10000"
	raw, err := ioutil.ReadFile("./mocks/filter.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEndpoint)
		fmt.Fprint(writer, string(raw))
	})

	filter, _, err := testClient.Filter.Get(10000)
	if filter == nil {
		t.Errorf("Expected Filter, got nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}

}

func TestFilterService_GetFavouriteList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/2/filter/favourite"
	raw, err := ioutil.ReadFile("./mocks/favourite_filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEndpoint)
		fmt.Fprint(writer, string(raw))
	})

	filters, _, err := testClient.Filter.GetFavouriteList()
	if filters == nil {
		t.Error("Expected Filters list. Filters list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestFilterService_GetMyFilters(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/filter/my"
	raw, err := ioutil.ReadFile("./mocks/my_filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEndpoint)
		fmt.Fprint(writer, string(raw))
	})

	opts := GetMyFiltersQueryOptions{}
	filters, _, err := testClient.Filter.GetMyFilters(&opts)
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if filters == nil {
		t.Errorf("Expected Filters, got nil")
	}
}

func TestFilterService_Search(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/filter/search"
	raw, err := ioutil.ReadFile("./mocks/search_filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEndpoint)
		fmt.Fprint(writer, string(raw))
	})

	opt := FilterSearchOptions{}
	filters, _, err := testClient.Filter.Search(&opt)
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if filters == nil {
		t.Errorf("Expected Filters, got nil")
	}
}
