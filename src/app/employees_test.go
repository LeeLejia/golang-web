package app

import (
	"testing"
	"net/url"
	"net/http"
	"strings"
	"net/http/httptest"
)

func TestCreateEmployee(t *testing.T) {
	InitTestDB()
	token := LoginTest()
	v := url.Values{}
	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)
	v.Add("oId","1")
	v.Add("name", "赵医生a")
	v.Add("phone", "15966666668")
	v.Add("workNo","200000000")
	v.Add("title", "主任医师")
	v.Add("sex", "男")

	req, err := http.NewRequest("POST", "/employees/create", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateEmployee)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"code":200,"msg":"SUCCESS","data":null}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}

	req, err = http.NewRequest("GET", "/employees/read?account=admin&osType=web&types=info&roleType=super&order=created_at DESC&oId=10&token="+token, strings.NewReader(v.Encode()))
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
	v.Add("accounts","15966666668")
	v.Add("state","2")

	req, err = http.NewRequest("POST", "/employees/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateEmployee)
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
