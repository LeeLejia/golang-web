package app

import (
	"fmt"
	"../conf"
	"../model"
	"../pdb"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"../common/logger"
)

func InitTestDB() {
	conf.Init("../../../app.toml")
	pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBNameTest)
	logger.InitLogger()
	stmt, err := pdb.Session.Prepare("DELETE FROM " + model.UserTableName())
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()
	stmt, err = pdb.Session.Prepare("DELETE FROM " + model.SystemLogTableName())
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()
	stmt, err = pdb.Session.Prepare("DELETE FROM " + model.OperateLogTableName())
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()
	stmt, err = pdb.Session.Prepare("DELETE FROM " + model.Employee2authTableName())
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()
	stmt, err = pdb.Session.Prepare("DELETE FROM " + model.EmployeeTableName() + " WHERE phone='15966666668'")
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()
	stmt, err = pdb.Session.Prepare("DELETE FROM " + model.WorkGroupTableName())
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Exec()
}

func TestRegister(t *testing.T) {
	InitTestDB()
	v := url.Values{}
	v.Add("account", "15999999999")
	v.Add("pwd", "123456")

	req, err := http.NewRequest("POST", "/register", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)
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
	v.Add("account", "")
	v.Add("pwd", "123456")

	req, err = http.NewRequest("POST", "/register", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Register)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"code":500,"msg":"手机号不能为空"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}
	v.Add("account", "15999999999")
	v.Add("pwd", "")

	req, err = http.NewRequest("POST", "/register", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Register)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"code":500,"msg":"请输入密码"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}
	v.Add("account", "1")
	v.Add("pwd", "123456")

	req, err = http.NewRequest("POST", "/register", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Register)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"code":500,"msg":"手机号格式错误"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}
	v.Add("account", "15999999999")
	v.Add("pwd", "123456")

	req, err = http.NewRequest("POST", "/register", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Register)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"code":500,"msg":"账号已存在，请直接登录"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	v = url.Values{}
	v.Add("account", "15999999999")
	v.Add("pwd", "123456")
	req, err = http.NewRequest("POST", "/login", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "")
	v.Add("pwd", "123456")
	req, err = http.NewRequest("POST", "/login", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "15999999999")
	v.Add("pwd", "")
	req, err = http.NewRequest("POST", "/login", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "15999999989")
	v.Add("pwd", "123456")
	req, err = http.NewRequest("POST", "/login", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "15999999999")
	v.Add("pwd", "1234567")
	v.Add("username", "test")
	v.Add("state", "2")
	v.Add("type", "3")
	req, err = http.NewRequest("POST", "/users/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateUser)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "")
	v.Add("pwd", "1234567")
	v.Add("username", "test")
	v.Add("state", "2")
	v.Add("type", "3")
	req, err = http.NewRequest("POST", "/users/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateUser)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("account", "15999999998")
	v.Add("pwd", "1234567")
	v.Add("username", "test")
	v.Add("state", "2")
	v.Add("type", "3")
	req, err = http.NewRequest("POST", "/users/update", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateUser)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	req, err = http.NewRequest("POST", "/users/read", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ListUsers)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	v = url.Values{}
	v.Add("q", "15999999999")
	v.Add("state", "2")
	req, err = http.NewRequest("POST", "/users/read", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ListUsers)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
