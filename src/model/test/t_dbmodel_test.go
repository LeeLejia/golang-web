package test

import (
	"testing"
	"fmt"
	"../../common/conf"
	"../../pdb"
	"github.com/bitly/go-simplejson"
	"time"
	".."
)

func TestT(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	sc:=model.SqlController {
		TableName:      "t_app",
		InsertColumns:  []string{"icon","app_id","version","name","describe","developer","valid","file","src","expend","download_count","created_at"},
		QueryColumns:   []string{"id","icon","app_id","version","name","describe","developer","valid","file","src","expend","download_count","created_at"},
		InSertFields:   InsertFields,
		QueryField2Obj: QueryField2Obj,
	}
	m,err:=model.GetModel(sc)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	app:=&model.T_app{
		Icon:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADhUlEQVRoQ+1aTUwTQRj9ZkvZtkARhGi0RBHDT6IBE0JMSKhCOMihSBC5ICoxJurVxJP8yJEYL0QTD0qACxKCePBCFCXpBQ4UgyEQFVSKVNTYSmmXdnfMNE6zbFpYoGVbXG7tPmbe+973vkk6iyDO/1Ao/o39XDWAUBRb2hhbdx07JOW0TkBDv+cMI8BTQHA0tsj/Y4NhXmDgam+d/g3lFxQQII9hJCaJS0gJCM5SEUEBjX2euZitvLSqGOa76/XZ5OuAgIY+7jyDhMF4qD7lKGCmpreefR4Q0NjnaQUELfEkADC0ddfrW1UBirmmOqBY6YNngpoBZT1QM6Bs/UE9B5Q24P92oOFUwigCwD0TfrPUCUuBxmpKZXzS7987BP24nS9YXQMjeSYXF9bp7U6hY+lotqWCzSULP7CuTdoWhULxJp0WdjKFRYV/ODyJEML0WXIiFGEMznuvOcenXzhXLi7iAu6YtaM56cxBsrBjRVi+O+wrDSXgcr933d4GLTg7LewKAOCmAc5EBWyGi7iArgs61/Qyb7O7MFQeTyi7MeR10bYgm4UjRp7dOq19W5KlMd9+6bW3VLA/iFNSAVLcshsOhxSxnRYyZzNjTcWJJW2vuNkVDgwdVaxp+IN/tHfCX0Y3kSOAkJaLi6gD96vYsRQW9l0f5AIZIJ/T9HCItMRmAkh2msvZA2s8dpD/DycgMwnsHed0yRQXMQGGRHA9qtYZxRUnk6T2hLZUHGZKzMVhG908SQsZGgaZMAZX+wj37eNPnCcXFzEBZHSSnuf8eIbjwUMXNrKo6Mtv3krDTImNL/DBnz4IdmpJ0I8t8PmrPkgVZ2UzXMQEPKll7X4B3O+WhEXxovmZTBoJIw3zRr0tZ1rJPuG3EmI6+wemfNYX0/y6sZmzH800l7N5tLViUkB7pdaalao5eXPIi2kLiCv1uIadTWDAsNF8l1ZWrtAdtxAJ70OLDn918lPSQ4suTvNBwnytWAvh5rsiLURO0LxMZs7uxBnf3Tg4LsVkxJgUHbiNLPJMLG784zBpPTm4HTsgO1S7DdxKiHebm6z9VAGyyhRFkOpAFIsra2nVAVlliiJIdSCKxZW19J5yIO4v+Yhll5555hHAEVn2KQzCAJ97LuoDl/F756KbqCG39QhDV6w6QSqPEVwJ+aqBuCsCmYixlz0EYGzkYlvavSHfVlG4xbe0/V/rXqJPNUQmQgAAAABJRU5ErkJggg==",
		AppId:"AGJLJWEFJ298R5pF",
		Name:"test_appoo",
		Version:"1.0",
		Describe:"a app for test",
		Developer:15,
		Valid:true,
		File:"XXXXXXXXXXXXX",
		Src:"XXXXXX",
		Expend:simplejson.New(),
		DownloadCount:124,
		CreatedAt:time.Now(),
	}
	// 插入数据
	err=m.Insert(app)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	// 查询数量
	count,err:=m.Count("")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("count:%d",count))
	count,err=m.Count("where id>$1",0)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("count:%d",count))
	// 更新数据
	err=m.Update("SET name = $1 WHERE name = $2","upName","test_appoo")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}else{
		fmt.Println("update success!")
	}
	// 获取信息
	res,err:=m.Query("")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	for _,obj:=range res{
		fmt.Println(*obj.(*model.T_app))
	}
}
func InsertFields(obj interface{}) []interface{} {
	fmt.Println(obj)
	app:=obj.(*model.T_app)
	expend := []byte{}
	if app.Expend != nil {
		bs, err := app.Expend.MarshalJSON()
		if err==nil{
			expend = bs
		}
	}
	return []interface{}{
		app.Icon,app.AppId,app.Version,app.Name,app.Describe,app.Developer,app.Valid,app.File,app.Src,expend,app.DownloadCount,app.CreatedAt,
	}
}
func QueryField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(model.GetByteArr(fields[10]))
	app:=&model.T_app{
		ID:model.GetInt64(fields[0],0),
		Icon:model.GetString(fields[1]),
		AppId:model.GetString(fields[2]),
		Name:model.GetString(fields[4]),
		Version:model.GetString(fields[3]),
		Describe:model.GetString(fields[5]),
		Developer:model.GetInt(fields[6],-1),
		Valid:model.GetBool(fields[7],true),
		File:model.GetString(fields[8]),
		Src:model.GetString(fields[9]),
		Expend:expend,
		DownloadCount:model.GetInt(fields[11],0),
		CreatedAt:model.GetTime(fields[12],time.Now()),
	}
	return app
}
