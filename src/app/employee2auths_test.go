package app

import (
	"testing"
	"net/url"
	"net/http"
	"strings"
	"net/http/httptest"
)

func TestAddEmployee2Auth(t *testing.T) {
	InitTestDB()
	token := LoginTest()
	v := url.Values{}
	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)

	req, err := http.NewRequest("GET", "/employees/auth/read?account=admin&osType=web&types=info&roleType=super&order=created_at DESC&auth=1&token="+token, strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListEmployee2Auth)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"code":200,"msg":"SUCCESS","data":{}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}
	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)

	req, err = http.NewRequest("GET", "/employees/auth/no_auth?account=admin&osType=web&types=info&roleType=super&order=created_at DESC&auth=1&oId=10&token="+token, strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ListEmployeeWithoutAuth)
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
	v.Add("auth","1")
	v.Add("employeeIds", "1,3")
	v.Add("method","add")

	req, err = http.NewRequest("POST", "/employees/auth/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(AddEmployee2Auth)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"code":200,"msg":"SUCCESS","data":null}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
