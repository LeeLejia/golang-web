package app

import (
	"testing"
	"net/url"
	"net/http"
	"strings"
	"net/http/httptest"
)

func TestListSystemLogs(t *testing.T) {
	InitTestDB()
	token := LoginTest()
	v := url.Values{}
	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)

	req, err := http.NewRequest("GET", "/system_logs/read?account=admin&osType=web&roleType=super&token="+token, strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListSystemLogs)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"code":200,"msg":"SUCCESS","data":{"systemLogs":[],"total":0}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}

	req, err = http.NewRequest("GET", "/system_logs/read?account=admin&osType=web&types=info&roleType=super&order=created_at DESC&start=2017-08-17 21:04:01&end=2017-08-19 21:04:01&token="+token, strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ListSystemLogs)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"code":200,"msg":"SUCCESS","data":{"systemLogs":[],"total":0}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
