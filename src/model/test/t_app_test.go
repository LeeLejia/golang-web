package test

import (
	"testing"
	"fmt"
	"../../common/conf"
	"github.com/bitly/go-simplejson"
	"time"
	m ".."
	"github.com/cjwddz/fast-model"
)

func TestT_app(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/task.toml")
	md, err := m.GetAppModel()
	if err != nil {
		t.Fail()
	}
	// 正常插入
	app := m.T_app{
		Icon:          "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADhUlEQVRoQ+1aTUwTQRj9ZkvZtkARhGi0RBHDT6IBE0JMSKhCOMihSBC5ICoxJurVxJP8yJEYL0QTD0qACxKCePBCFCXpBQ4UgyEQFVSKVNTYSmmXdnfMNE6zbFpYoGVbXG7tPmbe+973vkk6iyDO/1Ao/o39XDWAUBRb2hhbdx07JOW0TkBDv+cMI8BTQHA0tsj/Y4NhXmDgam+d/g3lFxQQII9hJCaJS0gJCM5SEUEBjX2euZitvLSqGOa76/XZ5OuAgIY+7jyDhMF4qD7lKGCmpreefR4Q0NjnaQUELfEkADC0ddfrW1UBirmmOqBY6YNngpoBZT1QM6Bs/UE9B5Q24P92oOFUwigCwD0TfrPUCUuBxmpKZXzS7987BP24nS9YXQMjeSYXF9bp7U6hY+lotqWCzSULP7CuTdoWhULxJp0WdjKFRYV/ODyJEML0WXIiFGEMznuvOcenXzhXLi7iAu6YtaM56cxBsrBjRVi+O+wrDSXgcr933d4GLTg7LewKAOCmAc5EBWyGi7iArgs61/Qyb7O7MFQeTyi7MeR10bYgm4UjRp7dOq19W5KlMd9+6bW3VLA/iFNSAVLcshsOhxSxnRYyZzNjTcWJJW2vuNkVDgwdVaxp+IN/tHfCX0Y3kSOAkJaLi6gD96vYsRQW9l0f5AIZIJ/T9HCItMRmAkh2msvZA2s8dpD/DycgMwnsHed0yRQXMQGGRHA9qtYZxRUnk6T2hLZUHGZKzMVhG908SQsZGgaZMAZX+wj37eNPnCcXFzEBZHSSnuf8eIbjwUMXNrKo6Mtv3krDTImNL/DBnz4IdmpJ0I8t8PmrPkgVZ2UzXMQEPKll7X4B3O+WhEXxovmZTBoJIw3zRr0tZ1rJPuG3EmI6+wemfNYX0/y6sZmzH800l7N5tLViUkB7pdaalao5eXPIi2kLiCv1uIadTWDAsNF8l1ZWrtAdtxAJ70OLDn918lPSQ4suTvNBwnytWAvh5rsiLURO0LxMZs7uxBnf3Tg4LsVkxJgUHbiNLPJMLG784zBpPTm4HTsgO1S7DdxKiHebm6z9VAGyyhRFkOpAFIsra2nVAVlliiJIdSCKxZW19J5yIO4v+Yhll5555hHAEVn2KQzCAJ97LuoDl/F756KbqCG39QhDV6w6QSqPEVwJ+aqBuCsCmYixlz0EYGzkYlvavSHfVlG4xbe0/V/rXqJPNUQmQgAAAABJRU5ErkJggg==",
		AppId:         "AGJLJWEFJ298R5pF",
		Name:          "test_app",
		Version:       "1.0",
		Describe:      "a app for test",
		Developer:     15,
		Valid:         true,
		File:          "XXXXXXXXXXXXX",
		Src:           "XXXXXX",
		Expend:        simplejson.New(),
		DownloadCount: 124,
		CreatedAt:     time.Now(),
	}
	err = md.Insert(app)
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
	err = md.Update(sq.And2("!=", "name", "old app").Set2("name","aaaaaaa"))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	err = md.Delete(cond.Reset().And2("!=", "name", "old app").And2(">","id",1000))
	fmt.Println(rs)
	fmt.Println(md.CountAll())
	fmt.Println(md.Count(cond))
}