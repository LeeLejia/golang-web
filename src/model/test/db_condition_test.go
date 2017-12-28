package test

import (
	"testing"
	"net/http"
	"net/url"
	".."
	"fmt"
)

func TestCondition(t *testing.T){
	r:=&http.Request{}
	r.PostForm = url.Values{}
	r.PostForm.Add("name","stupi")
	r.PostForm.Add("id","123")
	r.PostForm.Add("valid","TrUe")
	r.PostForm.Add("pos","10")
	r.PostForm.Add("len","50")

	cond:=model.DbCondition{}
	fmt.Println(cond.GetWhere())
	fmt.Println(cond.GetParams())

	cond.And(r,"=","s_name")
	fmt.Println(cond.GetWhere())
	fmt.Println(cond.GetParams())

	cond.And(r,"=","s_name")
	fmt.Println(cond.GetWhere())
	fmt.Println(cond.GetParams())

	cond.And(r,">","i_id").Or(r,"!=","b_valid")
	fmt.Println(cond.GetWhere())
	fmt.Println(cond.GetParams())

	cond.Limit(r,"pos","len").Order("Order by id desc")
	fmt.Println(cond.GetWhere())
	fmt.Println(cond.GetParams())

}
