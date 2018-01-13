package model

import (
	"fmt"
	"../pdb"
	"strings"
	"time"
)

type DbModel struct {
	sc SqlController
}
type SqlController struct {
	// 表名
	TableName string
	// 插入操作的列
	InsertColumns []string
	insertColumns string
	// 插入操作的占位符
	insertPlaceHold string
	// 查找操作获取到的列名
	QueryColumns []string
	queryColumns string
	// 查找操作列占位符
	queryPlaceHold string
	// 返回需要插入的列的集合obj，对应InsertColumns
	InSertFields func(obj interface{}) []interface{}
	// 将queryColumns返回的对象赋值到具体model对象,对应QueryColumns
	QueryField2Obj func(fields []interface{}) interface{}
}
// 获取一个model对象
func GetModel(sqlController SqlController) (DbModel,error){
	if sqlController.TableName==""{
		return DbModel{},fmt.Errorf("请配置TableName！")
	}
	if sqlController.InsertColumns==nil{
		return DbModel{},fmt.Errorf("InsertColumns为空")
	}
	if sqlController.QueryColumns==nil{
		return DbModel{},fmt.Errorf("QueryColumns为空")
	}
	if sqlController.InSertFields ==nil{
		return DbModel{},fmt.Errorf("CFMap为空")
	}
	sqlController.insertColumns = strings.Join(sqlController.InsertColumns,",")
	sqlController.queryColumns = strings.Join(sqlController.QueryColumns,",")
	ics:=make([]string,len(sqlController.InsertColumns))
	for i:=range sqlController.InsertColumns{
		ics[i]=fmt.Sprintf("$%d",i+1)
	}
	sqlController.insertPlaceHold = strings.Join(ics,",")

	qcs:=make([]string,len(sqlController.QueryColumns))
	for i:=range sqlController.QueryColumns{
		qcs[i]=fmt.Sprintf("$%d",i+1)
	}
	sqlController.queryPlaceHold = strings.Join(qcs,",")
	return DbModel{sc:sqlController},nil
}
// 获取表名
func (m *DbModel) GetTableName() string {
	return m.sc.TableName
}
// 设置插入列
func (m *DbModel) SetInsertColumns(columns []string,insertFileds func(obj interface{}) []interface{}){
	ics:=make([]string,len(columns))
	for i:=range columns{
		ics[i]=fmt.Sprintf("$%d",i+1)
	}
	m.sc.InsertColumns = columns
	m.sc.insertPlaceHold = strings.Join(ics,",")
	m.sc.insertColumns = strings.Join(columns,",")
	m.sc.InSertFields = insertFileds
}
// 设置搜索列
func (m *DbModel) SetQueryColumns(columns []string){
	qcs :=make([]string,len(columns))
	for i:=range columns{
		qcs[i]=fmt.Sprintf("$%d",i+1)
	}
	m.sc.QueryColumns = columns
	m.sc.queryPlaceHold = strings.Join(qcs,",")
	m.sc.queryColumns = strings.Join(m.sc.QueryColumns,",")
}
// 插入操作
func (m *DbModel) Insert(obj interface{}) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(%s) "+
		"VALUES(%s)", m.GetTableName(), m.sc.insertColumns, m.sc.insertPlaceHold))
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.sc.InSertFields(obj)...)
	return
}
/**
获取数据
conditionAndLimit：where id > $1 order by  limit $2 offset $3
 */
func (m *DbModel) Query(condiAOrderALimit string, args ...interface{}) (result []interface{}, err error) {
	sql:=fmt.Sprintf("SELECT %s FROM %s %s", m.sc.queryColumns, m.GetTableName(), condiAOrderALimit)
	fmt.Println(sql)
	stmt, err := pdb.Session.Prepare(sql)
	if err != nil {
		return result, err
	}
	rows,err:=stmt.Query(args...)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		refs := make([]interface{},0, len(m.sc.QueryColumns))
		for  range m.sc.QueryColumns {
			var ref interface{}
			refs = append(refs, &ref)
		}
		err = rows.Scan(refs...)
		row:=make([]interface{},len(m.sc.QueryColumns))
		for i:=range row{
			row[i] = *refs[i].(*interface{})
		}
		if err == nil {
			result = append(result, m.sc.QueryField2Obj(row))
		}
	}
	return result, err
}
/**
获取记录数量
condition:where id > $1 or name == $2
args: 5 'cjwddz'
 */
func (m *DbModel) Count(condition string,args ...interface{}) (count int, err error) {
	count = 0
	stmt,err:=pdb.Session.Prepare(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.GetTableName(), condition))
	if err!=nil{
		return 0,err
	}
	err = stmt.QueryRow(args...).Scan(&count)
	return
}
/**
更新数据
setAndCondition: SET id = $1,name=$2 where id = $3
args: 3 'cjwddz' 5
 */
func (m *DbModel) Update(setAndCondition string,args ...interface{}) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s %s", m.GetTableName(), setAndCondition))
	if err != nil {
		return
	}
	_, err = stmt.Exec(args...)
	return
}
/**
删除数据
condition: where id = $3
args: 3 'cjwddz' 5
 */
func (m *DbModel) Delete(condition string,args ...interface{})error{
	stmt,err:=pdb.Session.Prepare(fmt.Sprintf("DELETE FROM %s %s", m.GetTableName(), condition))
	if err!=nil{
		return err
	}
	_,err = stmt.Exec(args...)
	return err
}


/**
安全断言
 */
func GetInt(field interface{},def int)int{
	if field == nil{
		return def
	}
	if rs,ok:=field.(int);ok{
		return rs
	}
	if rs,ok:=field.(int64);ok{
		return int(rs)
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return def
}
func GetInt64(field interface{},def int64)int64{
	if field == nil{
		return def
	}
	if rs,ok:=field.(int64);ok{
		return int64(rs)
	}
	if rs,ok:=field.(int);ok{
		return int64(rs)
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return def
}
func GetByteArr(field interface{})[]byte{
	if field == nil{
		return nil
	}
	if rs,ok:=field.([]byte);ok{
		return rs
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return nil
}
func GetString(field interface{})string{
	if field ==nil{
		return ""
	}
	if rs,ok:=field.(string);ok{
		return rs
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return ""
}
func GetBool(field interface{},def bool)bool{
	if field==nil{
		return def
	}
	if rs,ok:=field.(bool);ok{
		return rs
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return def
}
func GetFloat(field interface{},def float32)float32{
	if field == nil{
		return def
	}
	if rs,ok:=field.(float32);ok{
		return rs
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return def
}
func GetFloat64(field interface{},def float64)float64{
	if field == nil{
		return def
	}
	if rs,ok:=field.(float64);ok{
		return rs
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return def
}
func GetTime(field interface{},def time.Time)time.Time{
	if field == nil{
		return def
	}
	if rs,ok:=field.(time.Time);ok{
		return rs
	}
	fmt.Println(fmt.Println("注意，该处强转可能出现错误！请在编译前检查！值："))
	fmt.Println(fmt.Println(field))
	return def
}