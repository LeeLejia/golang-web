package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateWorkGroup(t *testing.T) {
	InitTestDB()
	token := LoginTest()

	v := url.Values{}
	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)
	v.Add("name", "testing")
	v.Add("ownerID", "2")
	v.Add("ownerName", "123")

	req, err := http.NewRequest("POST", "/work_groups/create", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateWorkGroup)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}

	req, err = http.NewRequest("GET", "/work_groups/read?account=admin&osType=web&types=info&roleType=super&order=created_at DESC&oId=10&token="+token, strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ListEmployees)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)
	v.Add("name", "testing")
	v.Add("ownerID", "2")
	v.Add("ownerName", "123")

	req, err = http.NewRequest("POST", "/work_groups/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateWorkGroup)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
