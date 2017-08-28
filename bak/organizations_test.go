package app

import (
	"fmt"
	"../src/common"
	"../src/conf"
	"../src/model"
	"../src/pdb"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func LoginTest() (token string) {
	return common.SaveSession(1, "admin", "web", "super")
}

func TestCreateOrganization(t *testing.T) {
	conf.Init("../../../app.toml")
	pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBNameTest)
	stmt, err := pdb.Session.Prepare("DELETE FROM " + model.OrganizationTableName())
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()

	token := LoginTest()

	v := url.Values{}
	v.Add("name", "test")
	v.Add("code", "test")
	v.Add("contacts", "test")
	v.Add("contactsPhone", "test")
	v.Add("email", "test")
	v.Add("businessLicense", "test")
	v.Add("codePic", "test")

	req, err := http.NewRequest("POST", "/organizations/create", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrganization)
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
	v.Add("q", "15999999999")
	v.Add("state", "2")
	req, err = http.NewRequest("GET", "/organizations/read", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ListOrganizations)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("id", "1")
	v.Add("name", "test")
	v.Add("code", "test")
	v.Add("contacts", "test")
	v.Add("contactsPhone", "test")
	v.Add("email", "test")
	v.Add("businessLicense", "test")
	v.Add("codePic", "test")

	v.Add("account", "admin")
	v.Add("osType", "web")
	v.Add("roleType", "super")
	v.Add("token", token)
	req, err = http.NewRequest("POST", "/organizations/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateOrganization)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `{"code":200,"msg":"SUCCESS","data":null}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
