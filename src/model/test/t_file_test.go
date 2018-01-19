package test

import (
	"testing"
	"fmt"
	"../../common/conf"
	"time"
	m ".."
	"github.com/cjwddz/fast-model"
)

func TestT_file(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	md, err := m.GetFileModel()
	if err != nil {
		t.Fail()
	}
	// 正常插入
	file := m.T_File{
		Key:"test fileKey",
		Type:"test fileType",
		Name:"test fileName",
		Owner:-1,
		CreatedAt:time.Now(),
	}
	err = md.Insert(file)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	cond := model.DbCondition{}
	rs, err := md.Query(cond.And2(">", "id", 0).Limit2(4, -1))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	sq:=model.DbSetCondition{}
	err = md.Update(sq.And2("!=", "name", "test fileName").Set2("name","test fileNameUpdated"))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(rs)
	fmt.Println(md.CountAll())
	fmt.Println(md.Count(cond))
}